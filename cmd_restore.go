package redis

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// https://redis.io/commands/restore/
// RESTORE key ttl serialized-value [REPLACE] [ABSTTL] [IDLETIME seconds] [FREQ frequency]
func RestoreCommand(c *Client, args [][]byte) {
	if len(args) < 4 {
		c.Conn().WriteError(fmt.Sprintf(WrongNumOfArgsErr, args[0]))
		return
	}

	key := string(args[1])
	ttl, err := ParseTtlFromUnitTime(string(args[2]), int64(time.Millisecond))

	// Do not fail on time.Time{}, RESTORE will simply ignore it
	if err != nil {
		c.Conn().WriteError(InvalidIntErr)
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
			newTtl, err := ParseTtlFromTimestamp(string(args[2]), time.Millisecond)

			// Do not fail on time.Time{}, RESTORE will simply ignore it
			if err != nil {
				c.Conn().WriteError(InvalidIntErr)
				return
			}

			ttl = newTtl
		case "idletime":

			// We need 1 more argument for the time
			if len(args) == i+1 {
				c.Conn().WriteError(SyntaxErr)
			}

			i++

			// TODO: Use the given idle time.
		case "freq":
		default:
			c.Conn().WriteError(SyntaxErr)
			return
		}
	}

	db := c.Db()
	exists := db.Exists(key)

	if exists && !isRestore {
		c.Conn().WriteError("BUSYKEY Target key name already exists.")
		return
	}

	var kvp Kvp
	err = json.Unmarshal(args[3], &kvp)

	if err != nil {
		c.Conn().WriteError(fmt.Sprintf(DeserializationErr, string(args[3])))
		return
	}

	if kvp.Type == ValueTypeFancyString {
		str, ok := kvp.Value.(string)

		if !ok {
			c.Conn().WriteError(fmt.Sprintf(DeserializationErr, string(args[3])))
			return
		}

		db.Set(key, NewString(str), ttl)
	} else if kvp.Type == ValueTypeFancyList {
		arr, ok := kvp.Value.([]string)

		if !ok {
			c.Conn().WriteError(fmt.Sprintf(DeserializationErr, string(args[3])))
			return
		}

		db.Set(key, NewListFromArr(arr), ttl)
	} else if kvp.Type == ValueTypeFancySet {
		set, ok := kvp.Value.(map[string]struct{})

		if !ok {
			c.Conn().WriteError(fmt.Sprintf(DeserializationErr, string(args[3])))
			return
		}

		db.Set(key, NewSetFromMap(set), ttl)
	} else if kvp.Type == ValueTypeFancyZSet {
		set, ok := kvp.Value.(SortedSet[string, float64, struct{}])

		if !ok {
			c.Conn().WriteError(fmt.Sprintf(DeserializationErr, string(args[3])))
			return
		}

		db.Set(key, NewZSetFromSs(set), ttl)
	}

	c.Conn().WriteString("OK")
}
