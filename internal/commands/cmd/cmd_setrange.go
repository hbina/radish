package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/types"
	"github.com/hbina/radish/internal/util"
)

// https://redis.io/commands/setrange/
// SETRANGE key offset value
func SetrangeCommand(c *pkg.Client, args [][]byte) {
	if len(args) != 4 {
		c.Conn().WriteError(fmt.Sprintf(util.WrongNumOfArgsErr, args[0]))
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

	maybeItem, _ := db.Get(key)

	if maybeItem != nil && maybeItem.Type() != types.ValueTypeString {
		c.Conn().WriteError(util.WrongTypeErr)
	} else {
		if maybeItem == nil {

			if len(value) == 0 {
				db.Delete(key)
				c.Conn().WriteInt(0)
				return
			}

			maybeItem = types.NewString(string(make([]byte, byteOffset)))
		}

		item := maybeItem.(*types.String)

		// Need to append item to have enough spaces for byteOffset
		// TODO: Optimize to use StringBuilder
		if item.Len() <= int(byteOffset) {
			newStr := item.AsString() + string(make([]byte, byteOffset-item.Len())) + value
			item = types.NewString(newStr)
		} else {
			newStr := item.SubString(0, byteOffset) + value + item.SubString(byteOffset+len(value), item.Len())
			item = types.NewString(newStr)
		}

		db.Set(key, item, time.Time{})
		c.Conn().WriteInt(item.Len())
	}
}
