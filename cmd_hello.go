package redis

// https://redis.io/commands/hello/
// HELLO [protover [AUTH username password] [SETNAME clientname]]
// Stub implementation
func HelloCommand(c *Client, args [][]byte) {
	c.Conn().WriteString("OK")
}
