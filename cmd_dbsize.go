package redis

import "strings"

// https://redis.io/commands/dbsize/
func DbSizeCommand(c *Client, args [][]byte) {
	// TODO: For now it only returns stub values
	var str strings.Builder
	str.WriteString("redis_version:255.255.255\n")
	str.WriteString("redis_git_sha1:f36eb5a1")
	c.Conn().WriteBulkString(str.String())
}
