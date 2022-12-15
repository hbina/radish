package cmd

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/types"
	"github.com/hbina/radish/internal/util"
)

// https://redis.io/commands/restore/
// RESTORE key ttl serialized-value [REPLACE] [ABSTTL] [IDLETIME seconds] [FREQ frequency]
func RestoreCommand(c *pkg.Client, args [][]byte) {
	if len(args) < 4 {
		c.Conn().WriteError(fmt.Sprintf(util.WrongNumOfArgsErr, args[0]))
		return
	}

	key := string(args[1])
	ttl, err := util.ParseTtlFromUnitTime(string(args[2]), int64(time.Millisecond))

	// Do not fail on time.Time{}, RESTORE will simply ignore it
	if err != nil {
		c.Conn().WriteError(util.InvalidIntErr)
		return
	}

	isRestore := false

	// Parse the rest of options
	for i := 4; i < len(args); i++ {
		arg := strings.ToLower(string(args[i]))
		switch arg {
		case "replace":
			isRestore = true
		case "absttl":
			newTtl, err := util.ParseTtlFromTimestamp(string(args[2]), time.Millisecond)

			// Do not fail on time.Time{}, RESTORE will simply ignore it
			if err != nil {
				c.Conn().WriteError(util.InvalidIntErr)
				return
			}

			ttl = newTtl
		case "idletime":

			// We need 1 more argument for the time
			if len(args) == i+1 {
				c.Conn().WriteError(util.SyntaxErr)
			}

			i++

			// TODO: Use the given idle time.
		case "freq":
		default:
			c.Conn().WriteError(util.SyntaxErr)
			return
		}
	}

	db := c.Db()
	exists := db.Exists(key)

	if exists && !isRestore {
		c.Conn().WriteError("BUSYKEY Target key name already exists.")
		return
	}

	var kvp pkg.Kvp
	err = json.Unmarshal(args[3], &kvp)

	if err != nil {
		c.Conn().WriteError(fmt.Sprintf(util.DeserializationErr, string(args[3])))
		return
	}

	if kvp.Type == types.ValueTypeFancyString {
		set, ok := types.StringUnmarshal(kvp.Data)

		if !ok {
			c.Conn().WriteError(fmt.Sprintf(util.DeserializationErr, string(args[3])))
			return
		}

		db.Set(key, set, ttl)
	} else if kvp.Type == types.ValueTypeFancyList {
		set, ok := types.ListUnmarshal(kvp.Data)

		if !ok {
			c.Conn().WriteError(fmt.Sprintf(util.DeserializationErr, string(args[3])))
			return
		}

		db.Set(key, set, ttl)
	} else if kvp.Type == types.ValueTypeFancySet {
		set, ok := types.SetUnmarshal(kvp.Data)

		if !ok {
			c.Conn().WriteError(fmt.Sprintf(util.DeserializationErr, string(args[3])))
			return
		}

		db.Set(key, set, ttl)
	} else if kvp.Type == types.ValueTypeFancyZSet {
		set, ok := types.ZSetUnmarshal(kvp.Data)

		if !ok {
			c.Conn().WriteError(fmt.Sprintf(util.DeserializationErr, string(args[3])))
			return
		}

		db.Set(key, set, ttl)
	}
	c.Conn().WriteString("OK")
}
