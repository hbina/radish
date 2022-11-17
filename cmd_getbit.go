package redis

import (
	"fmt"
	"strconv"
)

// https://redis.io/commands/getbit/
// GETBIT key offset
func GetbitCommand(c *Client, args [][]byte) {
	if len(args) != 3 {
		c.Conn().WriteError(fmt.Sprintf(WrongNumOfArgsErr, args[0]))
		return
	}

	key := string(args[1])
	offsetStr := string(args[2])
	db := c.Db()

	byteOffset64, err := strconv.ParseInt(offsetStr, 10, 32)

	if err != nil {
		c.Conn().WriteError(SyntaxErr)
		return
	}

	byteOffset := int(byteOffset64)

	// Calculate the bitoffset
	bitOffset := byteOffset % 8

	// Recalibrate byteOffset because we can only get 1 byte at a time
	byteOffset /= 8

	maybeItem, _ := db.GetOrExpire(key, true)

	if maybeItem == nil {
		c.Conn().WriteInt(0)
	} else if maybeItem.Type() != ValueTypeString {
		c.Conn().WriteError(WrongTypeErr)
	} else {
		// Some tricky bit operations.
		// Please verify!
		mask := byte(0x80 >> byte(bitOffset))
		item := maybeItem.(*String)

		if item.Len() > byteOffset {
			bytes := []byte(item.inner)
			oldBit := 0

			if mask&bytes[byteOffset] > 0 {
				oldBit++
			}

			c.Conn().WriteInt(int(oldBit))
		} else {
			c.Conn().WriteInt(0)
		}
	}
}
