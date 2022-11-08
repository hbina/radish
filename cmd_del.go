package redis

func DelCommand(c *Client, args [][]byte) {
	db := c.Db()
	keys := make([]*string, 0, len(args)-1)
	for i := 1; i < len(args); i++ {
		k := string(args[i])
		keys = append(keys, &k)
	}
	dels := db.Delete(keys...)
	c.Conn().WriteInt(dels)
}
