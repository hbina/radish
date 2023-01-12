package cmd

import (
	"fmt"
	"strconv"

	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/types"
	"github.com/hbina/radish/internal/util"
)

// https://redis.io/commands/getbit/
// GETBIT key offset
func GetbitCommand(c *pkg.Client, args [][]byte) {
	if len(args) != 3 {
		c.WriteError(fmt.Sprintf(util.WrongNumOfArgsErr, args[0]))
		return
	}

	key := string(args[1])
	offsetStr := string(args[2])
	db := c.Db()

	byteOffset64, err := strconv.ParseInt(offsetStr, 10, 32)

	if err != nil {
		c.WriteError(util.SyntaxErr)
		return
	}

	byteOffset := int(byteOffset64)

	// Calculate the bitoffset
	bitOffset := byteOffset % 8

	// Recalibrate byteOffset because we can only get 1 byte at a time
	byteOffset /= 8

	maybeItem, _ := db.Get(key)

	if maybeItem == nil {
		c.WriteInt(0)
	} else if maybeItem.Type() != types.ValueTypeString {
		c.WriteError(util.WrongTypeErr)
	} else {
		// Some tricky bit operations.
		// Please verify!
		mask := byte(0x80 >> byte(bitOffset))
		item := maybeItem.(*types.String)

		if item.Len() > byteOffset {
			oldBit := 0

			if mask&item.AsBytes()[byteOffset] > 0 {
				oldBit++
			}

			c.WriteInt(int(oldBit))
		} else {
			c.WriteInt(0)
		}
	}
}
