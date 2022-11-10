package redis

import (
	"fmt"
	"strconv"
)

func LRangeCommand(c *Client, args [][]byte) {
	if len(args) < 3 {
		c.Conn().WriteError(fmt.Sprintf(WrongNumOfArgsErr, "lrange"))
		return
	}

	key := string(args[1])

	start, err := strconv.Atoi(string(args[2]))
	if err != nil {
		c.Conn().WriteError(fmt.Sprintf("%s: %s", InvalidIntErr, err.Error()))
		return
	}

	end, err := strconv.Atoi(string(args[3]))
	if err != nil {
		c.Conn().WriteError(fmt.Sprintf("%s: %s", InvalidIntErr, err.Error()))
		return
	}

	db := c.Db()
	item, _ := db.GetOrExpire(key, true)

	if item == nil {
		c.Conn().WriteArray(0)
		return
	} else if item.Type() != ValueTypeList {
		c.Conn().WriteError(WrongTypeErr)
		return
	}

	l := item.(*List)
	values := l.LRange(start, end)

	c.Conn().WriteArray(len(values))
	for _, v := range values {
		c.Conn().WriteBulkString(v)
	}
}
