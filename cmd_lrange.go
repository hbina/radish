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
	i := db.GetOrExpire(&key, true)
	if i == nil {
		c.Conn().WriteNull()
		return
	} else if i.Type() != ValueTypeList {
		c.Conn().WriteError(WrongTypeErr)
		return
	}

	l := i.(*List)
	values := l.LRange(start, end)

	c.Conn().WriteArray(len(values))
	for _, v := range values {
		c.Conn().WriteBulkString(v)
	}
}
