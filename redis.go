package redis

import (
	"fmt"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/tidwall/redcon"
)

const (
	SyntaxErr             = "ERR syntax error"
	InvalidIntErr         = "ERR value is not an integer or out of range"
	InvalidFloatErr       = "ERR value is not a valid float"
	WrongTypeErr          = "WRONGTYPE Operation against a key holding the wrong kind of value"
	WrongNumOfArgsErr     = "ERR wrong number of arguments for '%s' command"
	ZeroArgumentErr       = "ERR zero argument passed to the handler. This is an implementation bug"
	DeserializationErr    = "ERR unable to deserialize '%s' into a valid object"
	OptionNotSupportedErr = "ERR option '%s' is not currently supported"
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
			if len(cmd.Args) == 0 {
				c.Conn().WriteError(ZeroArgumentErr)
				return
			}

			// TODO: Check that args is not empty
			// TODO: Remove the first argument from argument to command handlers
			mu.Lock()
			defer mu.Unlock()
			cmdl := strings.ToLower(string(cmd.Args[0]))
			commandHandler := c.Redis().CommandHandlerFn(cmdl)
			if commandHandler != nil {
				(*commandHandler)(c, cmd.Args)
			} else {
				c.Redis().UnknownCommandFn()(c, cmd)
			}
		},
		unknownCommand: func(c *Client, cmd redcon.Command) {
			c.Conn().WriteError(fmt.Sprintf("ERR unknown command '%s'", cmd.Args[0]))
		},
		commands: make(Commands, 0),
	}
	r.redisDbs = make(map[DatabaseId]*RedisDb, 0)
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
		NewCommand("sadd", SaddCommand),
		NewCommand("smembers", SmembersCommand),
		NewCommand("smismember", SmismemberCommand),
		NewCommand("zadd", ZaddCommand),
		NewCommand("dump", DumpCommand),
		NewCommand("exists", ExistsCommand),
		NewCommand("restore", RestoreCommand),
		NewCommand("pttl", PttlCommand),
		NewCommand("debug", DebugCommand),
		NewCommand("srem", SremCommand),
		NewCommand("sintercard", SintercardCommand),
		NewCommand("sinter", SinterCommand),
		NewCommand("sinterstore", SinterstoreCommand),
		NewCommand("scard", ScardCommand),
		NewCommand("sismember", SismemberCommand),
		NewCommand("sunion", SunionCommand),
		NewCommand("sunionstore", SunionstoreCommand),
		NewCommand("sdiff", SdiffCommand),
		NewCommand("sdiffstore", SdiffstoreCommand),
		NewCommand("spop", SpopCommand),
		NewCommand("srandmember", SrandmemberCommand),
		NewCommand("smove", SmoveCommand),
		NewCommand("watch", WatchCommand),
		NewCommand("multi", MultiCommand),
		NewCommand("exec", ExecCommand),
		NewCommand("flushdb", FlushDbCommand),
		NewCommand("dbsize", DbSizeCommand),
		NewCommand("setx", SetXCommand),
		NewCommand("setnx", SetNxCommand),
		NewCommand("expire", ExpireCommand),
		NewCommand("setex", SetexCommand),
		NewCommand("getex", GetexCommand),
		NewCommand("getdel", GetdelCommand),
		NewCommand("mget", MgetCommand),
		NewCommand("getset", GetsetCommand),
		NewCommand("mset", MsetCommand),
		NewCommand("msetnx", MsetnxCommand),
		NewCommand("strlen", StrlenCommand),
		NewCommand("setbit", SetbitCommand),
		NewCommand("getbit", GetbitCommand),
		NewCommand("setrange", SetrangeCommand),
		NewCommand("getrange", GetrangeCommand),
		NewCommand("lcs", LcsCommand),
		NewCommand("zrange", ZrangeCommand),
		NewCommand("type", TypeCommand),
		NewCommand("zcard", ZcardCommand),
		NewCommand("zscore", ZscoreCommand),
		NewCommand("zincrby", ZincrbyCommand),
		NewCommand("zrem", ZremCommand),
		NewCommand("zrevrange", ZrevrangeCommand),
		NewCommand("zrank", ZrankCommand),
		NewCommand("zrevrank", ZrevrankCommand),
		NewCommand("zrangebyscore", ZrangebyscoreCommand),
		NewCommand("zrevrangebyscore", ZrevrangebyscoreCommand),
		NewCommand("zcount", ZcountCommand),
		NewCommand("zrangebylex", ZrangebylexCommand),
		NewCommand("zrevrangebylex", ZrevrangebylexCommand),
	})

	// NOTE: Taken by dumping from `CONFIG GET *`.
	// Is meaningless for the moment.
	// TODO: Implement parser for redis.conf and remove this.
	r.configDb = map[string]string{
		"rdbchecksum":                     "yes",
		"daemonize":                       "no",
		"io-threads-do-reads":             "no",
		"lua-replicate-commands":          "yes",
		"always-show-logo":                "yes",
		"protected-mode":                  "yes",
		"rdbcompression":                  "yes",
		"rdb-del-sync-files":              "no",
		"activerehashing":                 "yes",
		"stop-writes-on-bgsave-error":     "yes",
		"dynamic-hz":                      "yes",
		"lazyfree-lazy-eviction":          "no",
		"lazyfree-lazy-expire":            "no",
		"lazyfree-lazy-server-del":        "no",
		"lazyfree-lazy-user-del":          "no",
		"repl-disable-tcp-nodelay":        "no",
		"repl-diskless-sync":              "no",
		"gopher-enabled":                  "no",
		"aof-rewrite-incremental-fsync":   "yes",
		"no-appendfsync-on-rewrite":       "no",
		"cluster-require-full-coverage":   "yes",
		"rdb-save-incremental-fsync":      "yes",
		"aof-load-truncated":              "yes",
		"aof-use-rdb-preamble":            "yes",
		"cluster-replica-no-failover":     "no",
		"cluster-slave-no-failover":       "no",
		"replica-lazy-flush":              "no",
		"slave-lazy-flush":                "no",
		"replica-serve-stale-data":        "yes",
		"slave-serve-stale-data":          "yes",
		"replica-read-only":               "yes",
		"slave-read-only":                 "yes",
		"replica-ignore-maxmemory":        "yes",
		"slave-ignore-maxmemory":          "yes",
		"jemalloc-bg-thread":              "yes",
		"activedefrag":                    "no",
		"syslog-enabled":                  "no",
		"cluster-enabled":                 "no",
		"appendonly":                      "no",
		"cluster-allow-reads-when-down":   "no",
		"aclfile":                         "",
		"unixsocket":                      "",
		"pidfile":                         "/var/run/redis/redis-server.pid",
		"replica-announce-ip":             "",
		"slave-announce-ip":               "",
		"masteruser":                      "",
		"masterauth":                      "",
		"cluster-announce-ip":             "",
		"syslog-ident":                    "redis",
		"dbfilename":                      "dump.rdb",
		"appendfilename":                  "appendonly.aof",
		"server_cpulist":                  "",
		"bio_cpulist":                     "",
		"aof_rewrite_cpulist":             "",
		"bgsave_cpulist":                  "",
		"ignore-warnings":                 "ARM64-COW-BUG",
		"supervised":                      "systemd",
		"syslog-facility":                 "local0",
		"repl-diskless-load":              "disabled",
		"loglevel":                        "notice",
		"maxmemory-policy":                "noeviction",
		"appendfsync":                     "everysec",
		"oom-score-adj":                   "no",
		"databases":                       "16",
		"port":                            "6379",
		"io-threads":                      "1",
		"auto-aof-rewrite-percentage":     "100",
		"cluster-replica-validity-factor": "10",
		"cluster-slave-validity-factor":   "10",
		"list-max-ziplist-size":           "-2",
		"tcp-keepalive":                   "300",
		"cluster-migration-barrier":       "1",
		"active-defrag-cycle-min":         "1",
		"active-defrag-cycle-max":         "25",
		"active-defrag-threshold-lower":   "10",
		"active-defrag-threshold-upper":   "100",
		"lfu-log-factor":                  "10",
		"lfu-decay-time":                  "1",
		"replica-priority":                "100",
		"slave-priority":                  "100",
		"repl-diskless-sync-delay":        "5",
		"maxmemory-samples":               "5",
		"timeout":                         "0",
		"replica-announce-port":           "0",
		"slave-announce-port":             "0",
		"tcp-backlog":                     "511",
		"cluster-announce-bus-port":       "0",
		"cluster-announce-port":           "0",
		"repl-timeout":                    "60",
		"repl-ping-replica-period":        "10",
		"repl-ping-slave-period":          "10",
		"list-compress-depth":             "0",
		"rdb-key-save-delay":              "0",
		"key-load-delay":                  "0",
		"active-expire-effort":            "1",
		"hz":                              "10",
		"min-replicas-to-write":           "0",
		"min-slaves-to-write":             "0",
		"min-replicas-max-lag":            "10",
		"min-slaves-max-lag":              "10",
		"maxclients":                      "10000",
		"active-defrag-max-scan-fields":   "1000",
		"slowlog-max-len":                 "128",
		"acllog-max-len":                  "128",
		"lua-time-limit":                  "5000",
		"cluster-node-timeout":            "15000",
		"slowlog-log-slower-than":         "10000",
		"latency-monitor-threshold":       "0",
		"proto-max-bulk-len":              "536870912",
		"stream-node-max-entries":         "100",
		"repl-backlog-size":               "1048576",
		"maxmemory":                       "0",
		"hash-max-ziplist-entries":        "512",
		"set-max-intset-entries":          "512",
		"zset-max-ziplist-entries":        "128",
		"active-defrag-ignore-bytes":      "104857600",
		"hash-max-ziplist-value":          "64",
		"stream-node-max-bytes":           "4096",
		"zset-max-ziplist-value":          "64",
		"hll-sparse-max-bytes":            "3000",
		"tracking-table-max-keys":         "1000000",
		"repl-backlog-ttl":                "3600",
		"auto-aof-rewrite-min-size":       "67108864",
		"tls-port":                        "0",
		"tls-session-cache-size":          "20480",
		"tls-session-cache-timeout":       "300",
		"tls-cluster":                     "no",
		"tls-replication":                 "no",
		"tls-auth-clients":                "yes",
		"tls-prefer-server-ciphers":       "no",
		"tls-session-caching":             "yes",
		"tls-cert-file":                   "",
		"tls-key-file":                    "",
		"tls-dh-params-file":              "",
		"tls-ca-cert-file":                "",
		"tls-ca-cert-dir":                 "",
		"tls-protocols":                   "",
		"tls-ciphers":                     "",
		"tls-ciphersuites":                "",
		"logfile":                         "",
		"client-query-buffer-limit":       "1073741824",
		"watchdog-period":                 "0",
		"dir":                             "",
		"save":                            "900 1 300 10 60 10000",
		"client-output-buffer-limit":      "normal 0 0 0 slave 268435456 67108864 60 pubsub 33554432 8388608 60",
		"unixsocketperm":                  "0",
		"slaveof":                         "",
		"notify-keyspace-events":          "",
		"bind":                            "127.0.0.1 ::1",
		"requirepass":                     "",
		"oom-score-adj-values":            "0 200 800",
	}
	return r
}

// Flush all keys synchronously
func (r *Redis) SyncFlushAll() {
	for _, v := range r.redisDbs {
		v.Clear()
	}
}

// Flush the selected db
func (r *Redis) SyncFlushDb(dbId DatabaseId) {
	d, exists := r.redisDbs[dbId]

	if exists {
		d.Clear()
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

func (r *Redis) GetConfigValue(key string) *string {
	v, e := r.configDb[key]
	if e {
		return &v
	}
	return nil
}

func (r *Redis) SetConfigValue(key string, value string) {
	r.configDb[key] = value
}

// NewClient creates new client and adds it to the redis.
func (r *Redis) NewClient(conn redcon.Conn) *Client {
	c := &Client{
		conn:     conn,
		redis:    r,
		clientId: r.NextClientId(),
	}
	return c
}

// NextClientId atomically gets and increments a counter to return the next client id.
func (r *Redis) NextClientId() ClientId {
	id := atomic.AddUint64(&r.nextClientId, 1)
	return ClientId(id)
}

// Clients gets the current connected clients.
func (r *Redis) Clients() Clients {
	return r.clients
}

func (r *Redis) getClients() Clients {
	return r.clients
}
