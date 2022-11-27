package redis

// https://redis.io/commands/type/
// TYPE key
func TypeCommand(c *Client, args [][]byte) {
	if len(args) < 2 {
		c.Conn().WriteError(WrongNumOfArgsErr)
		return
	}

	key := string(args[1])
	db := c.Db()

	maybeItem, _ := db.GetOrExpire(key, true)

	if maybeItem == nil {
		c.Conn().WriteString("none")
	} else {
		c.Conn().WriteString(maybeItem.TypeFancy())
	}
}
