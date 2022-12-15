package pkg

import (
	"fmt"
	"net"
	"strings"
	"sync"

	"github.com/hbina/radish/internal/util"
)

type Redis struct {
	mu               *sync.RWMutex
	commands         map[string]*Command
	configs          map[string]string
	dbs              map[uint64]*Db
	blockingCommands map[string]*BlockingCommand
	retryList        []BlockedCommand
}

func Default(
	commands map[string]*Command,
	blockingCommands map[string]*BlockingCommand,
	configs map[string]string) *Redis {
	r := &Redis{
		mu:               new(sync.RWMutex),
		commands:         commands,
		blockingCommands: blockingCommands,
		configs:          configs,
		dbs:              make(map[uint64]*Db, 0),
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
	r.dbs[dbId] = NewRedisDb(dbId, r)
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
		conn:  &Conn{conn: conn},
		redis: r,
	}
	return c
}

func (r *Redis) HandleRequest(c *Client, args [][]byte) {
	util.Logger.Println(util.CollectArgs(args))

	if len(args) == 0 {
		c.Conn().WriteError(util.ZeroArgumentErr)
		return
	}

	// TODO: Check that args is not empty
	// TODO: Remove the first argument from argument to command handlers
	cmdName := strings.ToLower(string(args[0]))
	cmd := r.commands[cmdName]
	bcmd := r.blockingCommands[cmdName]

	cmdWrite := (cmd != nil && cmd.Flag&CMD_WRITE != 0) ||
		(bcmd != nil && bcmd.Flag&CMD_WRITE != 0)

	if cmdWrite {
		r.mu.Lock()
	} else {
		r.mu.RLock()
	}

	if cmd != nil {
		(cmd.Handler)(c, args)

		// Retry all the blocking commands
		r.HandleBlockedRequests()
	} else if bcmd != nil {
		err := (bcmd.Handler)(c, args)

		if err == BCMD_RETRY {
			r.AddBlockedRequest(c, args)
		}
	} else {
		c.Conn().WriteError(fmt.Sprintf("ERR unknown command '%s' with args '%s'", string(args[0]), args[1:]))
	}

	if cmdWrite {
		r.mu.Unlock()
	} else {
		r.mu.RUnlock()
	}
}

func (r *Redis) HandleBlockedRequests() {
	unfinished := make([]BlockedCommand, 0)

	for _, blockedCommand := range r.retryList {
		c := blockedCommand.c
		args := blockedCommand.args
		cmdName := strings.ToLower(string(args[0]))
		cmd := r.blockingCommands[cmdName]
		err := (cmd.Handler)(c, args)

		if err == BCMD_RETRY {
			unfinished = append(unfinished, blockedCommand)
		}
	}

	r.retryList = unfinished
}

func (r *Redis) AddBlockedRequest(c *Client, args [][]byte) {
	r.retryList = append(r.retryList, BlockedCommand{
		c:    c,
		args: args,
	})
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

		if resp != nil {
			util.Logger.Println(util.EscapeString(string(buffer)))
			buffer = leftover
			r.HandleRequest(client, util.ConvertRespToArgs(resp))
		}

		count, err = client.Read(tmp)

		if err != nil || count == 0 {
			return
		}
	}
}

func (r *Redis) RegisterCommands(cmds []*Command) {
	for _, cmd := range cmds {
		r.commands[cmd.Name] = cmd
	}
}

func (r *Redis) RegisterBlockingCommands(cmds []*BlockingCommand) {
	for _, cmd := range cmds {
		r.blockingCommands[cmd.Name] = cmd
	}
}
