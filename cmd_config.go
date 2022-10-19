package redis

import (
	"fmt"
	"strings"

	"github.com/redis-go/redcon"
)

func ConfigCommand(c *Client, cmd redcon.Command) {
	if len(cmd.Args) < 2 {
		c.Conn().WriteError("wrong number of arguments for 'config' command")
		return
	}

	subcommand := string(cmd.Args[1])

	if strings.ToLower(subcommand) == "get" {
		if len(cmd.Args) < 3 {
			c.Conn().WriteError(fmt.Sprintf("Unknown subcommand or wrong number of arguments for '%s'. Try CONFIG HELP.", string(cmd.Args[1])))
			return
		}
		// TODO: Update these to match our actual implementation
		redisConfig := map[string]string{
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

		// Support glob style pattern matching

		result := make([]string, 0, (len(cmd.Args)-2)*2)
		for i := 2; i < len(cmd.Args); i++ {
			k := string(cmd.Args[i])
			v, e := redisConfig[k]
			if e {
				result = append(result, fmt.Sprintf("\"%s\"", k), fmt.Sprintf("\"%s\"", v))
			}
		}
		c.Conn().WriteArray(len(result))
		for _, v := range result {
			c.Conn().WriteString(v)
		}
	} else if strings.ToLower(subcommand) == "set" {
		// Currently no-op

		if len(cmd.Args) < 4 {
			c.Conn().WriteError(fmt.Sprintf("Unknown subcommand or wrong number of arguments for '%s'. Try CONFIG HELP.", string(cmd.Args[1])))
			return
		}

		c.Conn().Close()
	} else {
		c.Conn().WriteError(fmt.Sprintf("Unknown subcommand '%s'. Try CONFIG HELP.", subcommand))
	}
}
