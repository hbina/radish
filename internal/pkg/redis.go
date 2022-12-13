package pkg

import (
	"fmt"
	"net"
	"strings"
	"sync"

	"github.com/hbina/radish/internal/util"
)

const (
	SyntaxErr             = "ERR syntax error"
	InvalidIntErr         = "ERR value is not an integer or out of range"
	InvalidFloatErr       = "ERR value is not a valid float"
	InvalidLexErr         = "ERR min or max not valid string range item"
	WrongTypeErr          = "WRONGTYPE Operation against a key holding the wrong kind of value"
	WrongNumOfArgsErr     = "ERR wrong number of arguments for '%s' command"
	ZeroArgumentErr       = "ERR zero argument passed to the handler. This is an implementation bug"
	DeserializationErr    = "ERR unable to deserialize '%s' into a valid object"
	OptionNotSupportedErr = "ERR option '%s' is not currently supported"
	NegativeIntErr        = "ERR %s must be greater than 0"
	MustBePositiveErr     = "ERR %s must be positive"
)

type Redis struct {
	mu               *sync.RWMutex
	commands         map[string]*Command
	configs          map[string]string
	dbs              map[uint64]*Db
	blockingCommands map[string]*BlockingCommand
	retryList        []BlockedCommand
}

func Default(commands map[string]*Command, blockingCommands map[string]*BlockingCommand) *Redis {
	r := &Redis{
		mu:               new(sync.RWMutex),
		commands:         commands,
		blockingCommands: blockingCommands,
		configs:          createConfigs(),
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
		c.Conn().WriteError(ZeroArgumentErr)
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

// NOTE: Taken by dumping from `CONFIG GET *`.
// Is meaningless for the moment.
// TODO: Implement parser for redis.conf and remove this.
func createConfigs() map[string]string {
	return map[string]string{
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
}
