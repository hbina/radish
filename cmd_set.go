package redis

import (
	"fmt"
	"strings"
	"time"
)

const (
	SetMode = iota
	SetEx
	SetPx
	SetExat
	SetPxat
)

const (
	SetExpireMode = iota
	SetExpireNx
	SetExpireXx
)

// https://redis.io/commands/set/
// SET key value [NX | XX] [GET] [EX seconds | PX milliseconds |
// EXAT unix-time-seconds | PXAT unix-time-milliseconds | KEEPTTL]
func SetCommand(c *Client, args [][]byte) {
	if len(args) == 0 {
		c.Conn().WriteError("no argument passed to handler. This should not be possible")
		return
	} else if len(args) < 3 {
		c.Conn().WriteError(fmt.Sprintf(WrongNumOfArgsErr, args[0]))
		return
	}

	key := string(args[1])
	value := string(args[2])

	var expire time.Time
	expireMode := SetMode
	writeMode := SetExpireMode
	shouldGet := false

	// Parse the optional arguments
	for i := 3; i < len(args); i++ {
		arg := strings.ToLower(string(args[i]))
		switch arg {
		default:
			c.Conn().WriteError(SyntaxErr)
			return
		case "ex":
			if expireMode != SetMode {
				c.Conn().WriteError(SyntaxErr)
				return
			}

			// We require 1 more argument for EX
			if len(args) == i+1 {
				c.Conn().WriteError(SyntaxErr)
				return
			}
			i++

			ttl, err := ParseExpiryTime(string(args[i]), uint64(time.Second))

			if ttl.IsZero() || err != nil {
				c.Conn().WriteError(InvalidIntErr)
				return
			}

			expire = ttl
			expireMode = SetEx
			continue
		case "px":
			if expireMode != SetMode {
				c.Conn().WriteError(SyntaxErr)
				return
			}

			// We require 1 more argument for PX
			if len(args) == i {
				c.Conn().WriteError(SyntaxErr)
				return
			}
			i++

			ttl, err := ParseExpiryTime(string(args[i]), uint64(time.Millisecond))

			if ttl.IsZero() || err != nil {
				c.Conn().WriteError(InvalidIntErr)
				return
			}

			expire = ttl
			expireMode = SetPx
			continue
		case "nx":
			if writeMode != SetExpireMode {
				c.Conn().WriteError(SyntaxErr)
				return
			}
			writeMode = SetExpireNx
			continue
		case "xx":
			if writeMode != SetExpireMode {
				c.Conn().WriteError(SyntaxErr)
				return
			}
			writeMode = SetExpireXx
			continue
		case "get":
			shouldGet = true
			continue
		}
	}

	found := false

	if shouldGet {
		item, _ := c.Db().GetOrExpire(key, true)
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
	if writeMode == SetExpireNx && exists || writeMode == SetExpireXx && !exists {
		c.Conn().WriteNull()
		return
	}

	db.Set(key, NewString(value), expire)

	if !found {
		c.Conn().WriteString("OK")
	}
}
