package commands

import (
	"github.com/hbina/radish/internal/commands/bcmd"
	"github.com/hbina/radish/internal/commands/cmd"
	"github.com/hbina/radish/internal/pkg"
)

func GenerateCommands() map[string]*pkg.Command {
	arr := []*pkg.Command{
		pkg.NewCommand("ping", cmd.PingCommand, pkg.CMD_READONLY),
		pkg.NewCommand("set", cmd.SetCommand, pkg.CMD_WRITE),
		pkg.NewCommand("get", cmd.GetCommand, pkg.CMD_READONLY),
		pkg.NewCommand("del", cmd.DelCommand, pkg.CMD_WRITE),
		pkg.NewCommand("ttl", cmd.TtlCommand, pkg.CMD_READONLY),
		pkg.NewCommand("lpush", cmd.LPushCommand, pkg.CMD_WRITE),
		pkg.NewCommand("rpush", cmd.RPushCommand, pkg.CMD_WRITE),
		pkg.NewCommand("lpop", cmd.LPopCommand, pkg.CMD_WRITE),
		pkg.NewCommand("rpop", cmd.RPopCommand, pkg.CMD_WRITE),
		pkg.NewCommand("lrange", cmd.LRangeCommand, pkg.CMD_READONLY),
		pkg.NewCommand("config", cmd.ConfigCommand, pkg.CMD_WRITE),
		pkg.NewCommand("info", cmd.InfoCommand, pkg.CMD_READONLY),
		pkg.NewCommand("select", cmd.SelectCommand, pkg.CMD_WRITE),
		pkg.NewCommand("flushall", cmd.FlushAllCommand, pkg.CMD_WRITE),
		pkg.NewCommand("function", cmd.FunctionCommand, pkg.CMD_WRITE),
		pkg.NewCommand("incr", cmd.IncrCommand, pkg.CMD_WRITE),
		pkg.NewCommand("incrby", cmd.IncrByCommand, pkg.CMD_WRITE),
		pkg.NewCommand("incrbyfloat", cmd.IncrByFloatCommand, pkg.CMD_WRITE),
		pkg.NewCommand("decr", cmd.DecrCommand, pkg.CMD_WRITE),
		pkg.NewCommand("decrby", cmd.DecrByCommand, pkg.CMD_WRITE),
		pkg.NewCommand("decrbyfloat", cmd.DecrByFloatCommand, pkg.CMD_WRITE),
		pkg.NewCommand("object", cmd.ObjectCommand, pkg.CMD_READONLY),
		pkg.NewCommand("sadd", cmd.SaddCommand, pkg.CMD_WRITE),
		pkg.NewCommand("smembers", cmd.SmembersCommand, pkg.CMD_WRITE),
		pkg.NewCommand("smismember", cmd.SmismemberCommand, pkg.CMD_WRITE),
		pkg.NewCommand("zadd", cmd.ZaddCommand, pkg.CMD_WRITE),
		pkg.NewCommand("dump", cmd.DumpCommand, pkg.CMD_READONLY),
		pkg.NewCommand("exists", cmd.ExistsCommand, pkg.CMD_READONLY),
		pkg.NewCommand("restore", cmd.RestoreCommand, pkg.CMD_WRITE),
		pkg.NewCommand("pttl", cmd.PttlCommand, pkg.CMD_READONLY),
		pkg.NewCommand("debug", cmd.DebugCommand, pkg.CMD_READONLY),
		pkg.NewCommand("srem", cmd.SremCommand, pkg.CMD_WRITE),
		pkg.NewCommand("sintercard", cmd.SintercardCommand, pkg.CMD_READONLY),
		pkg.NewCommand("sinter", cmd.SinterCommand, pkg.CMD_READONLY),
		pkg.NewCommand("sinterstore", cmd.SinterstoreCommand, pkg.CMD_WRITE),
		pkg.NewCommand("scard", cmd.ScardCommand, pkg.CMD_READONLY),
		pkg.NewCommand("sismember", cmd.SismemberCommand, pkg.CMD_READONLY),
		pkg.NewCommand("sunion", cmd.SunionCommand, pkg.CMD_READONLY),
		pkg.NewCommand("sunionstore", cmd.SunionstoreCommand, pkg.CMD_WRITE),
		pkg.NewCommand("sdiff", cmd.SdiffCommand, pkg.CMD_READONLY),
		pkg.NewCommand("sdiffstore", cmd.SdiffstoreCommand, pkg.CMD_WRITE),
		pkg.NewCommand("spop", cmd.SpopCommand, pkg.CMD_WRITE),
		pkg.NewCommand("srandmember", cmd.SrandmemberCommand, pkg.CMD_READONLY),
		pkg.NewCommand("smove", cmd.SmoveCommand, pkg.CMD_WRITE),
		pkg.NewCommand("watch", cmd.WatchCommand, pkg.CMD_READONLY),
		pkg.NewCommand("multi", cmd.MultiCommand, pkg.CMD_READONLY),
		pkg.NewCommand("exec", cmd.ExecCommand, pkg.CMD_READONLY),
		pkg.NewCommand("flushdb", cmd.FlushDbCommand, pkg.CMD_WRITE),
		pkg.NewCommand("dbsize", cmd.DbSizeCommand, pkg.CMD_READONLY),
		pkg.NewCommand("setx", cmd.SetXCommand, pkg.CMD_WRITE),
		pkg.NewCommand("setnx", cmd.SetNxCommand, pkg.CMD_WRITE),
		pkg.NewCommand("expire", cmd.ExpireCommand, pkg.CMD_WRITE),
		pkg.NewCommand("setex", cmd.SetexCommand, pkg.CMD_WRITE),
		pkg.NewCommand("getex", cmd.GetexCommand, pkg.CMD_READONLY),
		pkg.NewCommand("getdel", cmd.GetdelCommand, pkg.CMD_WRITE),
		pkg.NewCommand("mget", cmd.MgetCommand, pkg.CMD_WRITE),
		pkg.NewCommand("getset", cmd.GetsetCommand, pkg.CMD_WRITE),
		pkg.NewCommand("mset", cmd.MsetCommand, pkg.CMD_READONLY),
		pkg.NewCommand("msetnx", cmd.MsetnxCommand, pkg.CMD_WRITE),
		pkg.NewCommand("strlen", cmd.StrlenCommand, pkg.CMD_READONLY),
		pkg.NewCommand("setbit", cmd.SetbitCommand, pkg.CMD_WRITE),
		pkg.NewCommand("getbit", cmd.GetbitCommand, pkg.CMD_READONLY),
		pkg.NewCommand("setrange", cmd.SetrangeCommand, pkg.CMD_WRITE),
		pkg.NewCommand("getrange", cmd.GetrangeCommand, pkg.CMD_READONLY),
		pkg.NewCommand("lcs", cmd.LcsCommand, pkg.CMD_READONLY),
		pkg.NewCommand("zrange", cmd.ZrangeCommand, pkg.CMD_READONLY),
		pkg.NewCommand("type", cmd.TypeCommand, pkg.CMD_READONLY),
		pkg.NewCommand("zcard", cmd.ZcardCommand, pkg.CMD_READONLY),
		pkg.NewCommand("zscore", cmd.ZscoreCommand, pkg.CMD_READONLY),
		pkg.NewCommand("zincrby", cmd.ZincrbyCommand, pkg.CMD_WRITE),
		pkg.NewCommand("zrem", cmd.ZremCommand, pkg.CMD_WRITE),
		pkg.NewCommand("zrevrange", cmd.ZrevrangeCommand, pkg.CMD_READONLY),
		pkg.NewCommand("zrank", cmd.ZrankCommand, pkg.CMD_READONLY),
		pkg.NewCommand("zrevrank", cmd.ZrevrankCommand, pkg.CMD_READONLY),
		pkg.NewCommand("zrangebyscore", cmd.ZrangebyscoreCommand, pkg.CMD_READONLY),
		pkg.NewCommand("zrevrangebyscore", cmd.ZrevrangebyscoreCommand, pkg.CMD_READONLY),
		pkg.NewCommand("zcount", cmd.ZcountCommand, pkg.CMD_READONLY),
		pkg.NewCommand("zrangebylex", cmd.ZrangebylexCommand, pkg.CMD_READONLY),
		pkg.NewCommand("zrevrangebylex", cmd.ZrevrangebylexCommand, pkg.CMD_READONLY),
		pkg.NewCommand("zlexcount", cmd.ZlexcountCommand, pkg.CMD_READONLY),
		pkg.NewCommand("zremrangebyscore", cmd.ZremrangebyscoreCommand, pkg.CMD_WRITE),
		pkg.NewCommand("zremrangebylex", cmd.ZremrangebylexCommand, pkg.CMD_WRITE),
		pkg.NewCommand("zremrangebyrank", cmd.ZremrangebyrankCommand, pkg.CMD_WRITE),
		pkg.NewCommand("zinter", cmd.ZinterCommand, pkg.CMD_READONLY),
		pkg.NewCommand("zintercard", cmd.ZintercardCommand, pkg.CMD_READONLY),
		pkg.NewCommand("zinterstore", cmd.ZinterstoreCommand, pkg.CMD_WRITE),
		pkg.NewCommand("zunion", cmd.ZunionCommand, pkg.CMD_READONLY),
		pkg.NewCommand("zunioncard", cmd.ZunioncardCommand, pkg.CMD_READONLY),
		pkg.NewCommand("zunionstore", cmd.ZunionstoreCommand, pkg.CMD_WRITE),
		pkg.NewCommand("zdiff", cmd.ZdiffCommand, pkg.CMD_READONLY),
		pkg.NewCommand("zdiffcard", cmd.ZdiffcardCommand, pkg.CMD_READONLY),
		pkg.NewCommand("zdiffstore", cmd.ZdiffstoreCommand, pkg.CMD_WRITE),
		pkg.NewCommand("hello", cmd.HelloCommand, pkg.CMD_WRITE),
		pkg.NewCommand("zpopmin", cmd.ZpopminCommand, pkg.CMD_WRITE),
		pkg.NewCommand("zpopmax", cmd.ZpopmaxCommand, pkg.CMD_WRITE),
		pkg.NewCommand("zmpop", cmd.ZmpopCommand, pkg.CMD_WRITE),
		pkg.NewCommand("substr", cmd.SubstrCommand, pkg.CMD_READONLY),
	}

	res := make(map[string]*pkg.Command, len(arr))

	for _, r := range arr {
		res[r.Name] = r
	}

	return res
}

func GenerateBlockingCommands() map[string]*pkg.BlockingCommand {
	arr := []*pkg.BlockingCommand{
		pkg.NewBlockingCommand("bzmpop", bcmd.BzmpopCommand, pkg.CMD_WRITE),
		pkg.NewBlockingCommand("bzpopmin", bcmd.BzpopminCommand, pkg.CMD_WRITE),
		pkg.NewBlockingCommand("bzpopmax", bcmd.BzpopmaxCommand, pkg.CMD_WRITE),
	}

	res := make(map[string]*pkg.BlockingCommand, len(arr))

	for _, r := range arr {
		res[r.Name] = r
	}

	return res
}

// NOTE: Taken by dumping from `CONFIG GET *`.
// Is meaningless for the moment.
// TODO: Implement parser for redis.conf and remove this.
func GenerateConfigs() map[string]string {
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
