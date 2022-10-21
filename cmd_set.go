package redis

import (
	"fmt"
	"go-redis/ref"
	"strconv"
	"strings"
	"time"

	"github.com/tidwall/redcon"
)

const (
	SetNone uint = iota
	SetEx
	SetPx
	SetExat
	SetPxat
)

// SET key value [NX | XX] [GET] [EX seconds | PX milliseconds |
// EXAT unix-time-seconds | PXAT unix-time-milliseconds | KEEPTTL]
func SetCommand(c *Client, cmd redcon.Command) {
	if len(cmd.Args) == 0 {
		c.Conn().WriteError("no argument passed to handler. This should not be possible")
		return
	} else if len(cmd.Args) == 1 {
		c.Conn().WriteError(fmt.Sprintf("wrong number of arguments for '%s' command", cmd.Args[0]))
		return
	}

	key := string(cmd.Args[1])
	var value string
	if len(cmd.Args) > 1 {
		value = string(cmd.Args[2])
	}

	var expire *time.Time = nil
	var expireMode uint = SetNone
	var NX bool = false
	var XX bool = false

	if len(cmd.Args) > 2 {
		for i := 3; i < len(cmd.Args); i++ {
			arg := strings.ToLower(string(cmd.Args[i]))
			switch arg {
			default:
				c.Conn().WriteError(SyntaxErr)
				return
			case "ex":
				if expireMode != SetNone {
					c.Conn().WriteError(SyntaxErr)
					return
				}

				// We require 1 more argument for EX
				if len(cmd.Args) == i+1 {
					c.Conn().WriteError(SyntaxErr)
					return
				}

				// read next arg
				i++
				seconds, err := strconv.ParseUint(string(cmd.Args[i]), 10, 64)
				if err != nil {
					c.Conn().WriteError(fmt.Sprintf("%s: %s", InvalidIntErr, err.Error()))
					return
				}
				if seconds == 0 {
					c.Conn().WriteError("invalid expire time in 'set' command")
					return
				}
				expire = ref.Time(time.Now().Add(time.Duration(seconds * uint64(time.Second))))
				expireMode = SetEx
				continue
			case "px":
				if expireMode != SetNone {
					c.Conn().WriteError(SyntaxErr)
					return
				}

				// We require 1 more argument for PX
				if len(cmd.Args) == i {
					c.Conn().WriteError(SyntaxErr)
					return
				}

				// read next arg
				i++
				milliseconds, err := strconv.ParseUint(string(cmd.Args[i]), 10, 64)
				if err != nil {
					c.Conn().WriteError(fmt.Sprintf("%s: %s", InvalidIntErr, err.Error()))
					return
				}
				if milliseconds == 0 {
					c.Conn().WriteError("invalid expire time in 'set' command")
					return
				}
				expire = ref.Time(time.Now().Add(time.Duration(milliseconds * uint64(time.Millisecond))))
				expireMode = SetPx
				continue
			case "nx":
				if XX {
					c.Conn().WriteError(SyntaxErr)
					return
				}
				NX = true
				continue
			case "xx":
				if NX {
					c.Conn().WriteError(SyntaxErr)
					return
				}
				XX = true
				continue
			}
		}
	}

	// clients selected db
	db := c.Db()

	// TODO: Should we lock the database here?
	// The `Exists` and `Set` calls below should be atomic
	exists := db.Exists(&key)
	if NX && exists || XX && !exists {
		c.Conn().WriteNull()
		return
	}

	db.Set(&key, NewString(&value), expire)
	c.Conn().WriteString("OK")
}
