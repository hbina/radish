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
			c.WriteError("NOPROTO unsupported protocol version")
		}
	} else {
		c.WriteError("Unsupported operation")
	}
}

func writeStubResponse(c *pkg.Client) {
	c.WriteMap(7)
	c.WriteSimpleString("server")
	c.WriteSimpleString("redis")
	c.WriteSimpleString("version")
	c.WriteSimpleString("255.255.255")
	c.WriteSimpleString("proto")
	c.WriteInt(2)
	c.WriteSimpleString("id")
	c.WriteInt(12)
	c.WriteSimpleString("mode")
	c.WriteSimpleString("standalone")
	c.WriteSimpleString("role")
	c.WriteSimpleString("master")
	c.WriteSimpleString("modules")
	c.WriteArray(0)
}
