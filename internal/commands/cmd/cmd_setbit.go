package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/types"
	"github.com/hbina/radish/internal/util"
)

// https://redis.io/commands/setbit/
// SETBIT key offset value
func SetbitCommand(c *pkg.Client, args [][]byte) {
	if len(args) != 4 {
		c.WriteError(fmt.Sprintf(util.WrongNumOfArgsErr, args[0]))
		return
	}

	key := string(args[1])
	offsetStr := string(args[2])
	bitStr := string(args[3])
	db := c.Db()

	// Parse byteOffset
	byteOffset64, err := strconv.ParseInt(offsetStr, 10, 32)

	if err != nil || byteOffset64 < 0 {
		c.WriteError("ERR bit offset is not an integer or out of range")
		return
	}

	byteOffset := int(byteOffset64)

	// Calculate the bitOffset
	bitOffset := byteOffset % 8

	// Recalibrate byteOffset because we can only get 1 byte at a time
	byteOffset /= 8

	// Parse bitOffset
	if bitStr != "0" && bitStr != "1" {
		c.WriteError(util.InvalidIntErr)
		return
	}

	bit, err := strconv.ParseBool(bitStr)

	// Should not happen but you never know
	if err != nil {
		c.WriteError(util.SyntaxErr)
		return
	}

	maybeItem, _ := db.Get(key)

	if maybeItem != nil && maybeItem.Type() != types.ValueTypeString {
		c.WriteError(util.WrongTypeErr)
	} else {
		// Some tricky bit operations.
		// Please verify!

		if maybeItem == nil {
			maybeItem = types.NewString(string(make([]byte, byteOffset+1)))
		}

		mask := byte(0x80 >> byte(bitOffset))
		item := maybeItem.(*types.String)

		// Need to append item to have enough spaces for byteOffset
		if item.Len() <= int(byteOffset) {
			newStr := item.AsString() + string(make([]byte, byteOffset+1-item.Len()))
			item = types.NewString(newStr)
		}

		bytes := item.AsBytes()
		oldBit := 0

		if mask&bytes[byteOffset] > 0 {
			oldBit++
		}

		if bit {
			bytes[byteOffset] = bytes[byteOffset] | mask
		} else {
			bytes[byteOffset] = bytes[byteOffset] & (0xFF ^ mask)
		}

		db.Set(key, types.NewString(string(bytes)), time.Time{})
		c.WriteInt(int(oldBit))
	}
}
