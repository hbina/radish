package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/types"
	"github.com/hbina/radish/internal/util"
)

// https://redis.io/commands/getex/
// GETEX key [EX seconds | PX milliseconds | EXAT unix-time-seconds | PXAT unix-time-milliseconds | PERSIST]
func GetexCommand(c *pkg.Client, args [][]byte) {
	if len(args) < 2 {
		c.Conn().WriteError(fmt.Sprintf(util.WrongNumOfArgsErr, args[0]))
		return
	}

	key := string(args[1])
	db := c.Db()
	item, _ := db.Get(key)

	var newTtl time.Time
	expireMode := SetExpireMode

	// Parse the optional arguments
	for i := 2; i < len(args); i++ {
		arg := strings.ToLower(string(args[i]))
		switch arg {
		default:
			c.Conn().WriteError(util.SyntaxErr)
			return
		case "ex":
			if expireMode != SetExpireMode {
				c.Conn().WriteError(util.SyntaxErr)
				return
			}

			// We require 1 more argument for EX
			if len(args) == i+1 {
				c.Conn().WriteError(util.SyntaxErr)
				return
			}
			i++

			ttl, err := util.ParseTtlFromUnitTime(string(args[i]), int64(time.Second))

			if ttl.IsZero() || err != nil {
				c.Conn().WriteError(util.InvalidIntErr)
				return
			}

			newTtl = ttl
			expireMode = SetExpireEx
		case "px":
			if expireMode != SetExpireMode {
				c.Conn().WriteError(util.SyntaxErr)
				return
			}

			// We require 1 more argument for PX
			if len(args) == i {
				c.Conn().WriteError(util.SyntaxErr)
				return
			}
			i++

			ttl, err := util.ParseTtlFromUnitTime(string(args[i]), int64(time.Millisecond))

			if ttl.IsZero() || err != nil {
				c.Conn().WriteError(util.InvalidIntErr)
				return
			}

			newTtl = ttl
			expireMode = SetExpirePx
		case "exat":
			if expireMode != SetExpireMode {
				c.Conn().WriteError(util.SyntaxErr)
				return
			}

			// We require 1 more argument for EXAT
			if len(args) == i {
				c.Conn().WriteError(util.SyntaxErr)
				return
			}
			i++

			ttl, err := util.ParseTtlFromTimestamp(string(args[i]), time.Second)

			if err != nil || ttl.IsZero() {
				c.Conn().WriteError(util.InvalidIntErr)
				return
			}

			newTtl = ttl
			expireMode = SetExpireExat
		case "pxat":
			if expireMode != SetExpireMode {
				c.Conn().WriteError(util.SyntaxErr)
				return
			}

			// We require 1 more argument for PX
			if len(args) == i {
				c.Conn().WriteError(util.SyntaxErr)
				return
			}
			i++

			ttl, err := util.ParseTtlFromTimestamp(string(args[i]), time.Millisecond)

			if err != nil || ttl.IsZero() {
				c.Conn().WriteError(util.InvalidIntErr)
				return
			}

			newTtl = ttl
			expireMode = SetExpirePxat
		case "persist":
			if expireMode != SetExpireMode {
				c.Conn().WriteError(util.SyntaxErr)
				return
			}

			newTtl = time.Time{}
			expireMode = SetExpirePersist
		}
	}

	if item == nil {
		if c.R3 {
			c.Conn().WriteNull()
		} else {
			c.Conn().WriteNullBulk()
		}
		return
	}

	if item.Type() == types.ValueTypeString {
		v := item.Value().(string)
		c.Conn().WriteBulkString(v)
		// Only write the expiry ttl if the GET operation is successful
		db.SetExpiry(key, newTtl)
		return
	} else {
		c.Conn().WriteError(fmt.Sprintf("%s: key is a %s not a %s", util.WrongTypeErr, item.TypeFancy(), types.ValueTypeFancyString))
		return
	}
}
