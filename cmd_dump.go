package redis

import (
	"encoding/json"
	"fmt"
)

// https://redis.io/commands/dump/
func DumpCommand(c *Client, args [][]byte) {
	if len(args) == 0 {
		c.Conn().WriteError("no argument passed to handler. This should not be possible")
		return
	} else if len(args) < 2 {
		c.Conn().WriteError(fmt.Sprintf(WrongNumOfArgsErr, args[0]))
		return
	}

	key := string(args[1])
	value, _ := c.Db().GetOrExpire(key, true)

	if value == nil {
		c.Conn().WriteNull()
		return
	}

	if value.Type() == ValueTypeString {
		str, err := json.Marshal(Kvp{
			Key:   key,
			Type:  value.TypeFancy(),
			Value: value.Value(),
		})

		if err != nil {
			c.Conn().WriteError(err.Error())
			return
		}

		c.Conn().WriteBulkString(string(str))

		return
	} else if value.Type() == ValueTypeList {
		arr := make([]string, 0)

		value.(*List).ForEachF(func(a string) {
			arr = append(arr, a)
		})

		str, err := json.Marshal(Kvp{
			Key:   key,
			Type:  value.TypeFancy(),
			Value: arr,
		})

		if err != nil {
			c.Conn().WriteError(err.Error())
			return
		}

		c.Conn().WriteBulkString(string(str))

		return
	} else if value.Type() == ValueTypeSet {
		str, err := json.Marshal(Kvp{
			Key:   key,
			Type:  value.TypeFancy(),
			Value: value.(*Set).inner,
		})

		if err != nil {
			c.Conn().WriteError(err.Error())
			return
		}

		c.Conn().WriteBulkString(string(str))

		return
	} else if value.Type() == ValueTypeZSet {
		str, err := json.Marshal(Kvp{
			Key:   key,
			Type:  value.TypeFancy(),
			Value: value.(*ZSet).inner,
		})

		if err != nil {
			c.Conn().WriteError(err.Error())
			return
		}

		c.Conn().WriteBulkString(string(str))

		return
	}

	c.Conn().WriteError(fmt.Sprintf("Dump for %s is not yet implemented", value.TypeFancy()))
}
