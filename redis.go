package redis

import (
	"log"
	"strings"
	"sync"

	"github.com/tidwall/redcon"
)

const (
	SyntaxErr         = "ERR syntax error"
	InvalidIntErr     = "ERR value is not an integer or out of range"
	InvalidFloatErr   = "ERR value is not a valid float"
	WrongTypeErr      = "WRONGTYPE Operation against a key holding the wrong kind of value"
	WrongNumOfArgsErr = "ERR wrong number of arguments for '%s' command"
)

// This is the redis server.
type Redis struct {
	mu *sync.RWMutex

	// databases/keyspaces
	redisDbs map[DatabaseId]*RedisDb
	configDb map[string]string

	commands       Commands
	unknownCommand UnknownCommand

	handler Handler

	accept  Accept
	onClose OnClose

	// TODO version
	// TODO log writer
	// TODO modules
	// TODO redis options type

	keyExpirer KeyExpirer

	clients      Clients
	nextClientId uint64
}

// A Handler is called when a request is received and after Accept
// (if Accept allowed the connection by returning true).
//
// For implementing an own handler see the default handler
// as a perfect example in the createDefault() function.
type Handler func(c *Client, cmd redcon.Command)

// Accept is called when a Client tries to connect and before everything else,
// the Client connection will be closed instantaneously if the function returns false.
type Accept func(c *Client) bool

// OnClose is called when a Client connection is closed.
type OnClose func(c *Client, err error)

// Client map
type Clients map[ClientId]*Client

// Client id
type ClientId uint64

// Gets the handler func.
func (r *Redis) HandlerFn() Handler {
	return r.handler
}

// Sets the handler func.
// Live updates (while redis is running) works.
func (r *Redis) SetHandlerFn(new Handler) {
	r.handler = new
}

// Gets the accept func.
func (r *Redis) AcceptFn() Accept {
	return r.accept
}

// Sets the accept func.
// Live updates (while redis is running) works.
func (r *Redis) SetAcceptFn(new Accept) {
	r.accept = new
}

// Gets the onclose func.
func (r *Redis) OnCloseFn() OnClose {
	return r.onClose
}

// Sets the onclose func.
// Live updates (while redis is running) works.
func (r *Redis) SetOnCloseFn(new OnClose) {
	r.onClose = new
}

func (r *Redis) KeyExpirer() KeyExpirer {
	return r.keyExpirer
}

func (r *Redis) SetKeyExpirer(ke KeyExpirer) {
	r.keyExpirer = ke
}

var defaultRedis *Redis

// Default redis server.
// Initializes the default redis if not already.
// You can change the fields or value behind the pointer
// of the returned redis pointer to extend/change the default.
func Default() *Redis {
	if defaultRedis != nil {
		return defaultRedis
	}
	defaultRedis = createDefault()
	return defaultRedis
}

func escapeNewLines(s string) string {
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\r", "\\r")
	return s
}

// createDefault creates a new default redis.
func createDefault() *Redis {
	// initialize default redis server
	mu := new(sync.RWMutex)
	r := &Redis{
		mu: mu,
		accept: func(c *Client) bool {
			return true
		},
		onClose: func(c *Client, err error) {
		},
		handler: func(c *Client, cmd redcon.Command) {
			mu.Lock()
			defer mu.Unlock()
			log.Println(escapeNewLines(string(cmd.Raw)))
			cmdl := strings.ToLower(string(cmd.Args[0]))
			commandHandler := c.Redis().CommandHandlerFn(cmdl)
			if commandHandler != nil {
				(*commandHandler)(c, cmd.Args)
			} else {
				c.Redis().UnknownCommandFn()(c, cmd)
			}
		},
		unknownCommand: func(c *Client, cmd redcon.Command) {
			c.Conn().WriteError("ERR unknown command '" + string(cmd.Args[0]) + "'")
		},
		commands: make(Commands, 0),
	}
	r.redisDbs = make(map[DatabaseId]*RedisDb, redisDbMapSizeDefault)
	r.RedisDb(0) // initializes default db 0
	r.keyExpirer = KeyExpirer(NewKeyExpirer(r))

	r.RegisterCommands([]*Command{
		NewCommand("ping", PingCommand, CMD_STALE, CMD_FAST),
		NewCommand("set", SetCommand, CMD_WRITE, CMD_DENYOOM),
		NewCommand("get", GetCommand, CMD_READONLY, CMD_FAST),
		NewCommand("del", DelCommand, CMD_WRITE),
		NewCommand("ttl", TtlCommand, CMD_READONLY, CMD_FAST),

		NewCommand("lpush", LPushCommand, CMD_WRITE, CMD_FAST, CMD_DENYOOM),
		NewCommand("rpush", RPushCommand, CMD_WRITE, CMD_FAST, CMD_DENYOOM),
		NewCommand("lpop", LPopCommand, CMD_WRITE, CMD_FAST),
		NewCommand("rpop", RPopCommand, CMD_WRITE, CMD_FAST),
		NewCommand("lrange", LRangeCommand, CMD_READONLY),
		NewCommand("config", ConfigCommand, CMD_WRITE),
		NewCommand("info", InfoCommand, CMD_READONLY),
		NewCommand("select", SelectCommand, CMD_WRITE),
		NewCommand("flushall", FlushAllCommand, CMD_WRITE),
		NewCommand("function", FunctionCommand),
		NewCommand("incr", IncrCommand, CMD_WRITE),
		NewCommand("incrby", IncrByCommand, CMD_WRITE),
		NewCommand("incrbyfloat", IncrByFloatCommand, CMD_WRITE),
		NewCommand("decr", DecrCommand),
		NewCommand("decrby", DecrByCommand),
		NewCommand("decrbyfloat", DecrByFloatCommand),
		NewCommand("object", ObjectCommand),
	})
	return r
}

// Flush all keys synchronously
func (db *Redis) SyncFlushAll() {
	for _, v := range db.redisDbs {
		v.SyncFlushAll()
	}
}

// RedisDb gets the redis database by its id or creates and returns it if not exists.
func (r *Redis) RedisDb(dbId DatabaseId) *RedisDb {
	getDb := func() *RedisDb { // returns nil if db not exists
		if db, ok := r.redisDbs[dbId]; ok {
			return db
		}
		return nil
	}

	db := getDb()
	if db != nil {
		return db
	}

	// NOTE: This differs from original Redis because the number of databases are configured
	// at compile time with redis.conf
	// However, it should be fine to always return a valid database unless some application
	// rely on it to fail to stop?

	// now really create db of that id
	r.redisDbs[dbId] = NewRedisDb(dbId, r)
	return r.redisDbs[dbId]
}
