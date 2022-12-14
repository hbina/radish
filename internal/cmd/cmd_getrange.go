package cmd

import (
	"fmt"
	"strconv"

	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/types"
)

// https://redis.io/commands/getrange/
// GETRANGE key start end
func GetrangeCommand(c *pkg.Client, args [][]byte) {
	if len(args) != 4 {
		c.Conn().WriteError(fmt.Sprintf(pkg.WrongNumOfArgsErr, args[0]))
		return
	}

	key := string(args[1])
	startStr := string(args[2])
	endStr := string(args[3])
	db := c.Db()

	// Parse start index
	start64, err := strconv.ParseInt(startStr, 10, 64)

	if err != nil {
		c.Conn().WriteError("ERR start is not an integer or out of range")
		return
	}

	// TODO: this might be buggy in 32-bit computer
	start := int(start64)

	// Parse end index
	end64, err := strconv.ParseInt(endStr, 10, 64)

	if err != nil {
		c.Conn().WriteError("ERR end is not an integer or out of range")
		return
	}

	// We need to add 1 because its inclusive on both ends
	// TODO: this might be buggy in 32-bit computer
	end := int(end64) + 1

	maybeItem, _ := db.GetOrExpire(key, true)

	if maybeItem != nil && maybeItem.Type() != types.ValueTypeString {
		c.Conn().WriteError(pkg.WrongTypeErr)
	} else {
		if maybeItem == nil {
			c.Conn().WriteBulkString("")
			return
		}

		item := maybeItem.(*types.String)

		// If start/end is negative, recalibrate to the absolute index
		if start < 0 {
			start = item.Len() + start
			if start < 0 {
				start = 0
			}
		}
		if end <= 0 {
			end = item.Len() + end
			if start < 0 {
				start = 0
			}
		}

		if start > end {
			c.Conn().WriteBulkString("")
			return
		}

		if start >= item.Len() {
			start = item.Len()
		}

		if end >= item.Len() {
			end = item.Len()
		}

		str := item.Inner[start:end]
		c.Conn().WriteBulkString(str)
	}
}
