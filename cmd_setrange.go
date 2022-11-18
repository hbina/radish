package redis

import (
	"fmt"
	"strconv"
	"time"
)

// https://redis.io/commands/setrange/
// SETRANGE key offset value
func SetrangeCommand(c *Client, args [][]byte) {
	if len(args) != 4 {
		c.Conn().WriteError(fmt.Sprintf(WrongNumOfArgsErr, args[0]))
		return
	}

	key := string(args[1])
	offsetStr := string(args[2])
	value := string(args[3])
	db := c.Db()

	// Parse byteOffset
	byteOffset64, err := strconv.ParseInt(offsetStr, 10, 32)

	if err != nil || byteOffset64 < 0 {
		c.Conn().WriteError("ERR bit offset is not an integer or out of range")
		return
	}

	byteOffset := int(byteOffset64)

	// Redis strings can only go up to 512MB
	if byteOffset+len(value) > 536870911 {
		c.Conn().WriteError("ERR string exceeds maximum allowed size (proto-max-bulk-len)")
		return
	}

	maybeItem, _ := db.GetOrExpire(key, true)

	if maybeItem != nil && maybeItem.Type() != ValueTypeString {
		c.Conn().WriteError(WrongTypeErr)
	} else {
		if maybeItem == nil {
			maybeItem = NewString(string(make([]byte, byteOffset)))
		}

		item := maybeItem.(*String)

		// Need to append item to have enough spaces for byteOffset
		// TODO: Optimize to use StringBuilder
		if item.Len() <= int(byteOffset) {
			newStr := item.inner + string(make([]byte, byteOffset-item.Len())) + value
			item = NewString(newStr)
		} else {
			newStr := item.inner[:byteOffset] + value + item.inner[byteOffset+len(value):]
			item = NewString(newStr)
		}

		db.Set(key, item, time.Time{})
		c.Conn().WriteInt(item.Len())
	}
}
