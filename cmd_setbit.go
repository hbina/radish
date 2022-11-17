package redis

import (
	"fmt"
	"strconv"
	"time"
)

// https://redis.io/commands/setbit/
// SETBIT key offset value
func SetbitCommand(c *Client, args [][]byte) {
	if len(args) != 4 {
		c.Conn().WriteError(fmt.Sprintf(WrongNumOfArgsErr, args[0]))
		return
	}

	key := string(args[1])
	offsetStr := string(args[2])
	bitStr := string(args[3])
	db := c.Db()

	// Parse byteOffset
	byteOffset64, err := strconv.ParseInt(offsetStr, 10, 32)

	if err != nil || byteOffset64 < 0 {
		c.Conn().WriteError("ERR bit offset is not an integer or out of range")
		return
	}

	byteOffset := int(byteOffset64)

	// Calculate the bitOffset
	bitOffset := byteOffset % 8

	// Recalibrate byteOffset because we can only get 1 byte at a time
	byteOffset /= 8

	// Parse bitOffset
	if bitStr != "0" && bitStr != "1" {
		c.Conn().WriteError(InvalidIntErr)
		return
	}

	bit, err := strconv.ParseBool(bitStr)

	// Should not happen but you never know
	if err != nil {
		c.Conn().WriteError(SyntaxErr)
		return
	}

	maybeItem, _ := db.GetOrExpire(key, true)

	if maybeItem != nil && maybeItem.Type() != ValueTypeString {
		c.Conn().WriteError(WrongTypeErr)
	} else {
		// Some tricky bit operations.
		// Please verify!

		if maybeItem == nil {
			numOfBytes := byteOffset
			if bitOffset > 0 {
				numOfBytes += 1
			}
			maybeItem = NewString(string(make([]byte, numOfBytes)))
		}

		mask := byte(0x80 >> byte(bitOffset))
		item := maybeItem.(*String)

		// Need to append item to have enough spaces for byteOffset
		if item.Len() <= int(byteOffset) {
			newStr := item.inner + string(make([]byte, byteOffset+1-item.Len()))
			item = NewString(newStr)
		}

		bytes := []byte(item.inner)
		oldBit := 0

		if mask&bytes[byteOffset] > 0 {
			oldBit++
		}

		if bit {
			bytes[byteOffset] = bytes[byteOffset] | mask
		} else {
			bytes[byteOffset] = bytes[byteOffset] & (0xFF ^ mask)
		}

		db.Set(key, NewString(string(bytes)), time.Time{})
		c.Conn().WriteInt(int(oldBit))
	}
}
