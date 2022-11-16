package redis

import (
	"fmt"
	"strings"
	"time"
)

const (
	SetExpireMode = iota
	// Set key to expire after seconds
	SetExpireEx
	// Set key to expire after milliseconds
	SetExpirePx
	SetExpireExat
	SetExpirePxat
)

const (
	SetWriteMode = iota
	// Only write if key doesnt already exists
	SetWriteNx
	// Only write if key already exists
	SetWriteXx
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

			ttl, err := ParseExpiryTime(string(args[i]), uint64(time.Second))

			if ttl.IsZero() || err != nil {
				c.Conn().WriteError(InvalidIntErr)
				return
			}

			expire = ttl
			expireMode = SetExpireEx
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

			ttl, err := ParseExpiryTime(string(args[i]), uint64(time.Millisecond))

			if ttl.IsZero() || err != nil {
				c.Conn().WriteError(InvalidIntErr)
				return
			}

			expire = ttl
			expireMode = SetExpirePx
			continue
		case "nx":
			if writeMode != SetWriteMode {
				c.Conn().WriteError(SyntaxErr)
				return
			}
			writeMode = SetWriteNx
			continue
		case "xx":
			if writeMode != SetWriteMode {
				c.Conn().WriteError(SyntaxErr)
				return
			}
			writeMode = SetWriteXx
			continue
		case "get":
			shouldGet = true
			continue
		}
	}

	var foundItem Item = nil

	if shouldGet {
		item, _ := c.Db().GetOrExpire(key, true)
		if item != nil {
			if item.Type() == ValueTypeString {
				foundItem = item
			} else {
				c.Conn().WriteError(WrongTypeErr)
				return
			}
		}
	}

	db := c.Db()
	exists := db.Exists(&key)

	if writeMode == SetWriteNx && exists || writeMode == SetWriteXx && !exists {
		c.Conn().WriteInt(0)
		return
	}

	db.Set(key, NewString(value), expire)

	if writeMode != SetWriteMode {
		// If we are not in default write mode, then we must have written something
		c.Conn().WriteInt(1)
	} else if shouldGet && foundItem == nil {
		// If we should be getting something but found nothing then return nil
		c.Conn().WriteNull()
	} else if shouldGet && foundItem != nil {
		// We found something and we already checked that its a valid string Item
		c.Conn().WriteBulkString(*foundItem.Value().(*string))
	} else {
		// Otherwise just print OK
		c.Conn().WriteString("OK")
	}
}
