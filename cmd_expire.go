package redis

import (
	"strings"
	"time"
)

const (
	ExpireMode = iota
	ExpireNx
	ExpireXx
	ExpireGt
	ExpireLt
)

// https://redis.io/commands/expire/
// EXPIRE key seconds [NX | XX | GT | LT]
//
// NX -- Set expiry only when the key has no expiry.
// XX -- Set expiry only when the key has an existing expiry.
// GT -- Set expiry only when the new expiry is greater than current one.
// LT -- Set expiry only when the new expiry is less than current one.
func ExpireCommand(c *Client, args [][]byte) {
	// Handle special case with RENAME
	// Renaming keyA to keyB will cause keyA to inherit all
	// the timeout characteristics of keyB.

	// Calling EXPIRE with negative time will cause it to delete the key

	if len(args) < 3 || len(args) > 4 {
		c.Conn().WriteError(WrongNumOfArgsErr)
		return
	}

	key := string(args[1])
	seconds := string(args[2])

	mode := ExpireMode

	// Parse option
	if len(args) == 4 {
		switch strings.ToLower(string(args[3])) {
		case "nx":
			{
				mode = ExpireNx
			}
		case "xx":
			{
				mode = ExpireNx
			}
		case "gt":
			{
				mode = ExpireNx
			}
		case "lt":
			{
				mode = ExpireNx
			}
		}
	}

	newTtl, err := ParseTtlFromUnitTime(seconds, int64(time.Second))

	if err != nil {
		c.Conn().WriteError(InvalidIntErr)
		return
	}

	item, oldTtl := c.Db().GetOrExpire(key, true)

	if item == nil {
		c.Conn().WriteInt(0)
		return
	}

	if mode == ExpireNx && time.Time.IsZero(oldTtl) {
		c.Db().SetExpiry(key, newTtl)
	} else if mode == ExpireXx && !time.Time.IsZero(oldTtl) {
		c.Db().SetExpiry(key, newTtl)
	} else if mode == ExpireGt && newTtl.After(oldTtl) {
		c.Db().SetExpiry(key, newTtl)
	} else if mode == ExpireLt && newTtl.Before(oldTtl) {
		c.Db().SetExpiry(key, newTtl)
	} else {
		c.Db().SetExpiry(key, newTtl)
	}
	c.Conn().WriteInt(1)
}
