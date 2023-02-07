package pkg

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/hbina/radish/internal/util"
)

type Redis struct {
	cmds    map[string]*Command         // List of supported commands
	configs map[string]string           // Configurations (Currently unused)
	dbs     map[uint64]*Db              // List of database currently maintained
	bcmds   map[string]*BlockingCommand // List of supported blocked commands
	rlist   map[*Client]*BlockedCommand // List of commands to be retried for which clients
	bcmdTtl chan *Client
}

func Default(
	commands map[string]*Command,
	blockingCommands map[string]*BlockingCommand,
	configs map[string]string) *Redis {
	r := &Redis{
		cmds:    commands,
		bcmds:   blockingCommands,
		configs: configs,
		dbs:     make(map[uint64]*Db, 0),
		rlist:   make(map[*Client]*BlockedCommand, 0),
		bcmdTtl: make(chan *Client, 1),
	}
	return r
}

// Flush all keys synchronously
func (r *Redis) SyncFlushAll() {
	for _, v := range r.dbs {
		v.Clear()
	}
}

// Flush the selected db
func (r *Redis) SyncFlushDb(dbId uint64) {
	d, exists := r.dbs[dbId]

	if exists {
		d.Clear()
	}
}

// GetDb gets the redis database by its id or creates and returns it if not exists.
func (r *Redis) GetDb(dbId uint64) *Db {
	db, ok := r.dbs[dbId]

	if ok {
		return db
	}

	// NOTE: This differs from original Redis because the number of databases are configured
	// at compile time with redis.conf
	// However, it should be fine to always return a valid database unless some application
	// rely on it to fail to stop?

	// now really create db of that id
	r.dbs[dbId] = NewRedisDb(dbId)
	return r.dbs[dbId]
}

func (r *Redis) GetConfigValue(key string) *string {
	v, e := r.configs[key]
	if e {
		return &v
	}
	return nil
}

func (r *Redis) SetConfigValue(key string, value string) {
	r.configs[key] = value
}

// NewClient creates new client and adds it to the redis.
func (r *Redis) NewClient(conn net.Conn) *Client {
	c := &Client{
		conn:  util.NewConn(conn),
		redis: r,
		dbId:  0,
		R3:    false,
	}
	return c
}

func (r *Redis) HandleRequest(c *Client, args [][]byte) {
	if len(args) == 0 {
		c.Conn().WriteError(util.ZeroArgumentErr)
		return
	}

	// TODO: Check that args is not empty
	// TODO: Remove the first argument from argument to command handlers
	cmdName := strings.ToLower(string(args[0]))
	cmd := r.cmds[cmdName]
	bcmd := r.bcmds[cmdName]

	c.Db().Lock()

	if cmd != nil {
		(cmd.Handler)(c, args)
		r.HandleBlockedRequests(true)
	} else if bcmd != nil {
		err := (bcmd.Handler)(c, args)
		if err != nil {
			r.rlist[err.c] = err
			go time.AfterFunc(err.duration, func() {
				r.bcmdTtl <- err.c
			})
		} else {
			r.HandleBlockedRequests(false)
		}
	} else {
		c.Conn().WriteError(fmt.Sprintf("ERR unknown command '%s' with args '%s'", string(args[0]), args[1:]))
	}

	c.Db().Unlock()
}

// SAFETY: Some of the checks here have been ommitted because
// we already checked for them when we first received the command
func (r *Redis) HandleBlockedRequests(new bool) {
	for _, bcmd := range r.rlist {
		if !bcmd.ttl.IsZero() && time.Now().After(bcmd.ttl) {
			delete(r.rlist, bcmd.c)
		} else {
			cmdName := strings.ToLower(string(bcmd.args[0]))
			cmd := r.bcmds[cmdName]
			err := (cmd.Handler)(bcmd.c, bcmd.args)

			if err != nil {
				if !new { // If not a new blocked command, use the old TTL
					err.ttl = bcmd.ttl
				}
				r.rlist[err.c] = err
			} else {
				delete(r.rlist, bcmd.c)
			}
		}
	}
}

func (r *Redis) HandleClient(client *Client) {
	buffer := make([]byte, 0, 1024)
	tmp := make([]byte, 1024)
	count, err := client.Read(tmp)

	if err != nil {
		util.Logger.Fatal(err)
	}

	for {
		buffer = append(buffer, tmp[:count]...)

		// Try to parse the current buffer as a RESP
		resp, leftover := util.ConvertBytesToRespType(buffer)

		for resp != nil {
			util.Logger.Println(util.EscapeString(string(buffer)))
			buffer = leftover
			r.HandleRequest(client, util.ConvertRespToArgs(resp))
			resp, leftover = util.ConvertBytesToRespType(buffer)
		}

		count, err = client.Read(tmp)

		if err != nil || count == 0 {
			return
		}
	}
}

func (r *Redis) RegisterCommands(cmds []*Command) {
	for _, cmd := range cmds {
		r.cmds[cmd.Name] = cmd
	}
}

func (r *Redis) RegisterBlockingCommands(cmds []*BlockingCommand) {
	for _, cmd := range cmds {
		r.bcmds[cmd.Name] = cmd
	}
}

func (r *Redis) StartKeyExpiryJob(tick time.Duration) {
	f := func() {
		ticker := time.NewTicker(tick)
		for range ticker.C {
			for _, db := range r.RedisDbs() {
				db.Lock()
				db.DeleteExpiredKeys()
				db.Unlock()
			}
		}
	}
	go f()
}

func (r *Redis) StartBcmdTimeoutJob() {
	f := func() {
		for c := range r.bcmdTtl {
			c.Db().Lock()
			if c.R3 {
				c.Conn().WriteNull()
			} else {
				c.Conn().WriteNullArray()
			}
			delete(r.rlist, c)
			c.Db().Unlock()
		}
	}
	go f()
}
