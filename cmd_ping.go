package redis

import (
	"fmt"
	"strings"
)

// https://redis.io/commands/ping/
func PingCommand(c *Client, args [][]byte) {
	if len(args) == 1 {
		c.Conn().WriteString("PONG")
	} else if len(args) == 2 {
		var buf strings.Builder
		buf.WriteString(string(args[1]))
		s := buf.String()
		c.Conn().WriteBulkString(s)
	} else {
		c.Conn().WriteError(fmt.Sprintf(WrongNumOfArgsErr, "ping"))
	}
}
