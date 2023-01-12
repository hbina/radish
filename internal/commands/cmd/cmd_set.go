package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/types"
	"github.com/hbina/radish/internal/util"
)

const (
	SetExpireMode = iota
	// Set key to expire after seconds
	SetExpireEx
	// Set key to expire after milliseconds
	SetExpirePx
	SetExpireExat
	SetExpirePxat
	SetExpirePersist
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
func SetCommand(c *pkg.Client, args [][]byte) {
	if len(args) < 3 {
		c.Conn().WriteError(fmt.Sprintf(util.WrongNumOfArgsErr, args[0]))
		return
	}

	key := string(args[1])
	value := string(args[2])

	var newTtl time.Time
	expireMode := SetExpireMode
	writeMode := SetWriteMode
	shouldGet := false

	// Parse the optional arguments
	for i := 3; i < len(args); i++ {
		arg := strings.ToLower(string(args[i]))
		switch arg {
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
		case "nx":
			if writeMode != SetWriteMode {
				c.Conn().WriteError(util.SyntaxErr)
				return
			}

			writeMode = SetWriteNx
		case "xx":
			if writeMode != SetWriteMode {
				c.Conn().WriteError(util.SyntaxErr)
				return
			}

			writeMode = SetWriteXx
		case "get":
			shouldGet = true
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

			// We require 1 more argument for EXAT
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
			expireMode = SetExpireExat
		default:
			c.Conn().WriteError(util.SyntaxErr)
			return
		}
	}

	var foundStr *types.String = nil

	if shouldGet {
		item, _ := c.Db().Get(key)
		if item != nil {
			if item.Type() == types.ValueTypeString {
				foundStr = item.(*types.String)
			} else {
				c.Conn().WriteError(util.WrongTypeErr)
				return
			}
		}
	}

	db := c.Db()
	exists := db.Exists(key)

	if writeMode == SetWriteNx && exists || writeMode == SetWriteXx && !exists {
		if shouldGet {
			if foundStr == nil {
				c.Conn().WriteNull()
			} else {
				c.Conn().WriteBulkString(foundStr.AsString())
			}
		} else {
			c.Conn().WriteNull()
		}
		return
	}

	db.Set(key, types.NewString(value), newTtl)

	if shouldGet {
		if foundStr == nil {
			c.Conn().WriteNull()
		} else {
			// We already checked that foundStr is a *types.String
			c.Conn().WriteBulkString(foundStr.AsString())
		}
	} else {
		c.Conn().WriteString("OK")
	}
}
