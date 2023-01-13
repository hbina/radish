package cmd

import "github.com/hbina/radish/internal/pkg"

// https://redis.io/commands/hello/
// HELLO [protover [AUTH username password] [SETNAME clientname]]
// Stub implementation
func HelloCommand(c *pkg.Client, args [][]byte) {
	if len(args) == 1 {
		writeStubResponse(c)
	} else if len(args) == 2 {
		version := string(args[1])

		if version == "2" {
			c.UseResp2()
			writeStubResponse(c)
		} else if version == "3" {
			c.UseResp3()
			writeStubResponse(c)
		} else {
			c.Conn().WriteError("NOPROTO unsupported protocol version")
		}
	} else {
		c.Conn().WriteError("Unsupported operation")
	}
}

func writeStubResponse(c *pkg.Client) {
	c.Conn().WriteMap(7 * 2)
	c.Conn().WriteString("server")
	c.Conn().WriteString("redis")
	c.Conn().WriteString("version")
	c.Conn().WriteString("255.255.255")
	c.Conn().WriteString("proto")
	c.Conn().WriteInt(2)
	c.Conn().WriteString("id")
	c.Conn().WriteInt(12)
	c.Conn().WriteString("mode")
	c.Conn().WriteString("standalone")
	c.Conn().WriteString("role")
	c.Conn().WriteString("master")
	c.Conn().WriteString("modules")
	c.Conn().WriteArray(0)
}
