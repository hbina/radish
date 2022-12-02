package redis

// https://redis.io/commands/dbsize/
func DbSizeCommand(c *Client, args [][]byte) {
	db := c.Db()
	c.Conn().WriteInt(db.Len())
}
