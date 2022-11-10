package redis

import (
	"fmt"
	"strings"
)

const (
	PingTooManyArguments = "ERR wrong number of arguments for 'ping' command"
	ZeroArgument         = "ERR zero argument provided. This is a bug with the implementation"
)

func PingCommand(c *Client, args [][]byte) {
	if len(args) == 0 {
		c.Conn().WriteError(ZeroArgument)
	} else if len(args) == 1 {
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
