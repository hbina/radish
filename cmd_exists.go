package redis

import "fmt"

// https://redis.io/commands/exists/
func ExistsCommand(c *Client, args [][]byte) {
	if len(args) < 2 {
		c.Conn().WriteError(fmt.Sprintf(WrongNumOfArgsErr, args[0]))
		return
	}

	db := c.Db()
	count := 0

	for i := 1; i < len(args); i++ {
		value, _ := db.GetOrExpire(string(args[i]), true)
		if value != nil {
			count++
		}
	}

	c.Conn().WriteInt(count)
}
