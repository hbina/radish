package redis

import (
	"fmt"
	"strings"
	"sync"

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
	redisDbs map[uint64]*RedisDb
	configDb map[string]string

	commands       map[string]*Command
	unknownCommand UnknownCommand

	handler func(c *Client, cmd redcon.Command)

	keyExpirer *Expirer

	clients map[string]*Client
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
		handler: func(c *Client, cmd redcon.Command) {
			if len(cmd.Args) == 0 {
				c.Conn().WriteError(ZeroArgumentErr)
				return
			}

			// TODO: Check that args is not empty
			// TODO: Remove the first argument from argument to command handlers
			// fmt.Println(CollectArgs(cmd.Args))
			cmdl := strings.ToLower(string(cmd.Args[0]))
			command := c.Redis().Command(cmdl)

			fmt.Println(CollectArgs(cmd.Args))

			if command != nil {
				if command.flag&CMD_WRITE != 0 {
					mu.Lock()
				} else {
					mu.RLock()
				}

				(command.handler)(c, cmd.Args)

				if command.flag&CMD_WRITE != 0 {
					mu.Unlock()
				} else {
					mu.RUnlock()
				}
			} else {
				c.Redis().UnknownCommandFn()(c, cmd)
			}
		},
		unknownCommand: func(c *Client, cmd redcon.Command) {
			c.Conn().WriteError(fmt.Sprintf("ERR unknown command '%s'", cmd.Args[0]))
		},
		commands: make(map[string]*Command, 0),
		clients:  make(map[string]*Client),
		redisDbs: make(map[uint64]*RedisDb, 0),
	}
	r.keyExpirer = NewKeyExpirer(r)

	r.RegisterCommands([]*Command{
		NewCommand("ping", PingCommand, CMD_READONLY),
		NewCommand("set", SetCommand, CMD_WRITE),
		NewCommand("get", GetCommand, CMD_READONLY),
		NewCommand("del", DelCommand, CMD_WRITE),
		NewCommand("ttl", TtlCommand, CMD_READONLY),
		NewCommand("lpush", LPushCommand, CMD_WRITE),
		NewCommand("rpush", RPushCommand, CMD_WRITE),
		NewCommand("lpop", LPopCommand, CMD_WRITE),
		NewCommand("rpop", RPopCommand, CMD_WRITE),
		NewCommand("lrange", LRangeCommand, CMD_READONLY),
		NewCommand("config", ConfigCommand, CMD_WRITE),
		NewCommand("info", InfoCommand, CMD_READONLY),
		NewCommand("select", SelectCommand, CMD_WRITE),
		NewCommand("flushall", FlushAllCommand, CMD_WRITE),
		NewCommand("function", FunctionCommand, CMD_WRITE),
		NewCommand("incr", IncrCommand, CMD_WRITE),
		NewCommand("incrby", IncrByCommand, CMD_WRITE),
		NewCommand("incrbyfloat", IncrByFloatCommand, CMD_WRITE),
		NewCommand("decr", DecrCommand, CMD_WRITE),
		NewCommand("decrby", DecrByCommand, CMD_WRITE),
		NewCommand("decrbyfloat", DecrByFloatCommand, CMD_WRITE),
		NewCommand("object", ObjectCommand, CMD_READONLY),
		NewCommand("sadd", SaddCommand, CMD_WRITE),
		NewCommand("smembers", SmembersCommand, CMD_WRITE),
		NewCommand("smismember", SmismemberCommand, CMD_WRITE),
		NewCommand("zadd", ZaddCommand, CMD_WRITE),
		NewCommand("dump", DumpCommand, CMD_READONLY),
		NewCommand("exists", ExistsCommand, CMD_READONLY),
		NewCommand("restore", RestoreCommand, CMD_WRITE),
		NewCommand("pttl", PttlCommand, CMD_READONLY),
		NewCommand("debug", DebugCommand, CMD_READONLY),
		NewCommand("srem", SremCommand, CMD_WRITE),
		NewCommand("sintercard", SintercardCommand, CMD_READONLY),
		NewCommand("sinter", SinterCommand, CMD_READONLY),
		NewCommand("sinterstore", SinterstoreCommand, CMD_WRITE),
		NewCommand("scard", ScardCommand, CMD_READONLY),
		NewCommand("sismember", SismemberCommand, CMD_READONLY),
		NewCommand("sunion", SunionCommand, CMD_READONLY),
		NewCommand("sunionstore", SunionstoreCommand, CMD_WRITE),
		NewCommand("sdiff", SdiffCommand, CMD_READONLY),
		NewCommand("sdiffstore", SdiffstoreCommand, CMD_WRITE),
		NewCommand("spop", SpopCommand, CMD_WRITE),
		NewCommand("srandmember", SrandmemberCommand, CMD_READONLY),
		NewCommand("smove", SmoveCommand, CMD_WRITE),
		NewCommand("watch", WatchCommand, CMD_READONLY),
		NewCommand("multi", MultiCommand, CMD_READONLY),
		NewCommand("exec", ExecCommand, CMD_READONLY),
		NewCommand("flushdb", FlushDbCommand, CMD_WRITE),
		NewCommand("dbsize", DbSizeCommand, CMD_READONLY),
		NewCommand("setx", SetXCommand, CMD_WRITE),
		NewCommand("setnx", SetNxCommand, CMD_WRITE),
		NewCommand("expire", ExpireCommand, CMD_WRITE),
		NewCommand("setex", SetexCommand, CMD_WRITE),
		NewCommand("getex", GetexCommand, CMD_READONLY),
		NewCommand("getdel", GetdelCommand, CMD_WRITE),
		NewCommand("mget", MgetCommand, CMD_WRITE),
		NewCommand("getset", GetsetCommand, CMD_WRITE),
		NewCommand("mset", MsetCommand, CMD_READONLY),
		NewCommand("msetnx", MsetnxCommand, CMD_WRITE),
		NewCommand("strlen", StrlenCommand, CMD_READONLY),
		NewCommand("setbit", SetbitCommand, CMD_WRITE),
		NewCommand("getbit", GetbitCommand, CMD_READONLY),
		NewCommand("setrange", SetrangeCommand, CMD_WRITE),
		NewCommand("getrange", GetrangeCommand, CMD_READONLY),
		NewCommand("lcs", LcsCommand, CMD_READONLY),
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
func (r *Redis) SyncFlushDb(dbId uint64) {
	d, exists := r.redisDbs[dbId]

	if exists {
		d.Clear()
	}
}

// RedisDb gets the redis database by its id or creates and returns it if not exists.
func (r *Redis) RedisDb(dbId uint64) *RedisDb {
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
		conn:  conn,
		redis: r,
	}
	return c
}
