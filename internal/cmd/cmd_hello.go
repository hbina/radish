package cmd

import "github.com/hbina/radish/internal/pkg"

// https://redis.io/commands/hello/
// HELLO [protover [AUTH username password] [SETNAME clientname]]
// Stub implementation
func HelloCommand(c *pkg.Client, args [][]byte) {
	c.Conn().WriteString("OK")
}
