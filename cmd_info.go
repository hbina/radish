package redis

import "strings"

// https://redis.io/commands/info/
func InfoCommand(c *Client, args [][]byte) {
	// TODO: These are just stub values at the moment
	// Implementing this probably requires some modification to the build system
	var str strings.Builder
	str.WriteString("redis_version:255.255.255\r\n")
	str.WriteString("redis_git_sha1:f36eb5a1\r\n\r\n")
	str.WriteString("# Stats\r\n")
	str.WriteString("migrate_cached_sockets:0\r\n")
	c.Conn().WriteBulkString(str.String())
}
