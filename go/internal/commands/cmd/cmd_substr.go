package cmd

import (
	"fmt"
	"strconv"

	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/types"
	"github.com/hbina/radish/internal/util"
)

// https://redis.io/commands/substr/
// SUBSTR key start end
func SubstrCommand(c *pkg.Client, args [][]byte) {
	if len(args) != 4 {
		c.Conn().WriteError(fmt.Sprintf(util.WrongNumOfArgsErr, args[0]))
		return
	}

	db := c.Db()
	key := string(args[1])
	startStr := string(args[2])
	endStr := string(args[3])

	start64, err := strconv.ParseInt(startStr, 10, 64)

	if err != nil {
		c.Conn().WriteError(util.SyntaxErr)
		util.Logger.Println(err)
		return
	}

	end64, err := strconv.ParseInt(endStr, 10, 64)

	if err != nil {
		c.Conn().WriteError(util.SyntaxErr)
		util.Logger.Println(err)
		return
	}

	maybeItem, _ := db.Get(key)

	if maybeItem == nil {
		c.Conn().WriteBulkString("")
		return
	} else if maybeItem.Type() != types.ValueTypeString {
		c.Conn().WriteError(util.WrongTypeErr)
		return
	}

	item := maybeItem.(*types.String)
	itemLen64 := int64(item.Len())

	actualStart := start64
	if actualStart < 0 {
		actualStart = itemLen64 + start64
		if actualStart < 0 {
			actualStart = 0
		}
	} else if actualStart > itemLen64 {
		actualStart = itemLen64
	}

	// If start is beyond the length, just return empty string
	if actualStart > itemLen64 {
		util.Logger.Printf("%d", actualStart)
		c.Conn().WriteBulkString("")
		return
	}

	// The calculation for end is a bit weird because its inclusive
	actualEnd := end64
	if actualEnd < 0 {
		actualEnd = itemLen64 + end64 + 1
		if actualEnd < 0 {
			actualEnd = 1
		}
	} else if actualEnd > itemLen64 {
		actualEnd = itemLen64
	} else {
		actualEnd = actualEnd + 1
	}

	if actualStart > actualEnd {
		actualStart = actualEnd
	}

	util.Logger.Printf("%d %d", actualStart, actualEnd)
	substr := item.SubString(int(actualStart), int(actualEnd))
	c.Conn().WriteBulkString(substr)
}
