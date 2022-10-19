package redis

import (
	"strings"

	"github.com/tidwall/redcon"
)

const (
	PingTooManyArguments = "ERR wrong number of arguments for 'ping' command"
	ZeroArgument         = "ERR zero argument provided. This is a bug with the implementation"
)

// TODO: Since "PING" is already routed here, perhaps we
// should remove the first argument entirely so that
// it's not representable in the first place (zero arguments).
func PingCommand(c *Client, cmd redcon.Command) {
	if len(cmd.Args) == 0 {
		c.Conn().WriteError(ZeroArgument)
	} else if len(cmd.Args) == 1 {
		c.Conn().WriteString("PONG")
	} else if len(cmd.Args) == 2 {
		var buf strings.Builder
		buf.WriteString("\"")
		buf.Write(cmd.Args[1])
		buf.WriteString("\"")
		s := buf.String()
		c.Conn().WriteString(s)
	} else {
		c.Conn().WriteError(PingTooManyArguments)
	}
}
