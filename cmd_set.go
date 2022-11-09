package redis

import (
	"fmt"
	"go-redis/ref"
	"strconv"
	"strings"
	"time"
)

const (
	SetExpireMode uint = iota
	SetEx
	SetPx
	SetExat
	SetPxat
)

const (
	SetWriteMode uint = iota
	SetNx
	SetXx
)

func getExpiryTime(c *Client, arg string, multiplier uint64) *time.Time {
	unitTime, err := strconv.ParseUint(string(arg), 10, 64)
	if err != nil {
		c.Conn().WriteError(fmt.Sprintf("%s: %s", InvalidIntErr, err.Error()))
		return nil
	}
	if unitTime == 0 {
		c.Conn().WriteError("invalid expire time in 'set' command")
		return nil
	}

	return ref.Time(time.Now().Add(time.Duration(unitTime * multiplier)))
}

// SET key value [NX | XX] [GET] [EX seconds | PX milliseconds |
// EXAT unix-time-seconds | PXAT unix-time-milliseconds | KEEPTTL]
func SetCommand(c *Client, args [][]byte) {
	if len(args) == 0 {
		c.Conn().WriteError("no argument passed to handler. This should not be possible")
		return
	} else if len(args) < 3 {
		c.Conn().WriteError(fmt.Sprintf("wrong number of arguments for '%s' command", args[0]))
		return
	}

	key := string(args[1])
	value := string(args[2])

	var expire *time.Time = nil
	expireMode := SetExpireMode
	writeMode := SetWriteMode
	shouldGet := false

	// Parse the optional arguments
	for i := 3; i < len(args); i++ {
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

			expire = getExpiryTime(c, string(args[i]), uint64(time.Second))
			expireMode = SetEx
			continue
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

			expire = getExpiryTime(c, string(args[i]), uint64(time.Millisecond))
			expireMode = SetPx
			continue
		case "nx":
			if writeMode != SetWriteMode {
				c.Conn().WriteError(SyntaxErr)
				return
			}
			writeMode = SetNx
			continue
		case "xx":
			if writeMode != SetWriteMode {
				c.Conn().WriteError(SyntaxErr)
				return
			}
			writeMode = SetXx
			continue
		case "get":
			shouldGet = true
			continue
		}
	}

	found := false

	if shouldGet {
		item := c.Db().GetOrExpire(key, true)
		if item == nil {
			// c.Conn().WriteNull()
		} else {
			if item.Type() == ValueTypeString {
				v := *item.Value().(*string)
				c.Conn().WriteBulkString(v)
				found = true
			}
		}
	}

	db := c.Db()

	exists := db.Exists(&key)
	if writeMode == SetNx && exists || writeMode == SetXx && !exists {
		c.Conn().WriteNull()
		return
	}

	db.Set(key, NewString(value), expire)

	if !found {
		c.Conn().WriteString("OK")
	}
}
