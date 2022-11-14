package redis

// https://redis.io/commands/dbsize/
func DbSizeCommand(c *Client, args [][]byte) {
	c.Conn().WriteInt(c.Db().Len())
}
