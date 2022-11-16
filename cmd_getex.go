package redis

import (
	"fmt"
	"strings"
	"time"
)

// https://redis.io/commands/getex/
// GETEX key [EX seconds | PX milliseconds | EXAT unix-time-seconds | PXAT unix-time-milliseconds | PERSIST]
func GetexCommand(c *Client, args [][]byte) {
	if len(args) == 0 {
		c.Conn().WriteError("no argument passed to handler. This should not be possible")
		return
	} else if len(args) < 2 {
		c.Conn().WriteError(fmt.Sprintf(WrongNumOfArgsErr, args[0]))
		return
	}

	key := string(args[1])
	db := c.Db()
	item, _ := db.GetOrExpire(key, true)

	var newTtl time.Time
	expireMode := SetExpireMode

	// Parse the optional arguments
	for i := 2; i < len(args); i++ {
		arg := strings.ToLower(string(args[i]))
		switch arg {
		default:
			c.Conn().WriteError(SyntaxErr)
			return
		case "ex":
			if expireMode != SetExpireMode {
				c.Conn().WriteError(SyntaxErr)
				return
			}

			// We require 1 more argument for EX
			if len(args) == i+1 {
				c.Conn().WriteError(SyntaxErr)
				return
			}
			i++

			ttl, err := ParseTtlFromUnitTime(string(args[i]), int64(time.Second))

			if ttl.IsZero() || err != nil {
				c.Conn().WriteError(InvalidIntErr)
				return
			}

			newTtl = ttl
			expireMode = SetExpireEx
		case "px":
			if expireMode != SetExpireMode {
				c.Conn().WriteError(SyntaxErr)
				return
			}

			// We require 1 more argument for PX
			if len(args) == i {
				c.Conn().WriteError(SyntaxErr)
				return
			}
			i++

			ttl, err := ParseTtlFromUnitTime(string(args[i]), int64(time.Millisecond))

			if ttl.IsZero() || err != nil {
				c.Conn().WriteError(InvalidIntErr)
				return
			}

			newTtl = ttl
			expireMode = SetExpirePx
		case "exat":
			if expireMode != SetExpireMode {
				c.Conn().WriteError(SyntaxErr)
				return
			}

			// We require 1 more argument for PX
			if len(args) == i {
				c.Conn().WriteError(SyntaxErr)
				return
			}
			i++

			ttl, err := ParseTtlFromTimestamp(string(args[i]), time.Second)

			if err != nil || ttl.IsZero() {
				c.Conn().WriteError(InvalidIntErr)
				return
			}

			newTtl = ttl
			expireMode = SetExpireExat
		case "pxat":
			if expireMode != SetExpireMode {
				c.Conn().WriteError(SyntaxErr)
				return
			}

			// We require 1 more argument for PX
			if len(args) == i {
				c.Conn().WriteError(SyntaxErr)
				return
			}
			i++

			ttl, err := ParseTtlFromTimestamp(string(args[i]), time.Millisecond)

			if err != nil || ttl.IsZero() {
				c.Conn().WriteError(InvalidIntErr)
				return
			}

			newTtl = ttl
			expireMode = SetExpirePxat
		case "persist":
			if expireMode != SetExpireMode {
				c.Conn().WriteError(SyntaxErr)
				return
			}

			newTtl = time.Time{}
			expireMode = SetExpirePersist
		}
	}

	if item == nil {
		c.Conn().WriteNull()
		return
	}

	if item.Type() == ValueTypeString {
		v := item.Value().(string)
		c.Conn().WriteBulkString(v)
		// Only write the expiry ttl if the GET operation is successful
		db.SetExpiry(key, newTtl)
		return
	} else {
		c.Conn().WriteError(fmt.Sprintf("%s: key is a %s not a %s", WrongTypeErr, item.TypeFancy(), ValueTypeFancyString))
		return
	}
}
