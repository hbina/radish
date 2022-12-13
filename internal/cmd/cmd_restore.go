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
		c.Conn().WriteError(fmt.Sprintf(pkg.WrongNumOfArgsErr, args[0]))
		return
	}

	key := string(args[1])
	ttl, err := util.ParseTtlFromUnitTime(string(args[2]), int64(time.Millisecond))

	// Do not fail on time.Time{}, RESTORE will simply ignore it
	if err != nil {
		c.Conn().WriteError(pkg.InvalidIntErr)
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
				c.Conn().WriteError(pkg.InvalidIntErr)
				return
			}

			ttl = newTtl
		case "idletime":

			// We need 1 more argument for the time
			if len(args) == i+1 {
				c.Conn().WriteError(pkg.SyntaxErr)
			}

			i++

			// TODO: Use the given idle time.
		case "freq":
		default:
			c.Conn().WriteError(pkg.SyntaxErr)
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
		c.Conn().WriteError(fmt.Sprintf(pkg.DeserializationErr, string(args[3])))
		return
	}

	if kvp.Type == types.ValueTypeFancyString {
		str, ok := kvp.Value.(string)

		if !ok {
			c.Conn().WriteError(fmt.Sprintf(pkg.DeserializationErr, string(args[3])))
			return
		}

		db.Set(key, types.NewString(str), ttl)
	} else if kvp.Type == types.ValueTypeFancyList {
		arr, ok := kvp.Value.([]string)

		if !ok {
			c.Conn().WriteError(fmt.Sprintf(pkg.DeserializationErr, string(args[3])))
			return
		}

		db.Set(key, types.NewListFromArr(arr), ttl)
	} else if kvp.Type == types.ValueTypeFancySet {
		set, ok := kvp.Value.(map[string]struct{})

		if !ok {
			c.Conn().WriteError(fmt.Sprintf(pkg.DeserializationErr, string(args[3])))
			return
		}

		db.Set(key, types.NewSetFromMap(set), ttl)
	} else if kvp.Type == types.ValueTypeFancyZSet {
		pair, ok := kvp.Value.(map[string]interface{})

		if !ok {
			c.Conn().WriteError(fmt.Sprintf(pkg.DeserializationErr, string(args[3])))
			return
		}

		set := types.NewZSet()

		keys, ok1 := pair["keys"].([]interface{})
		scores, ok2 := pair["scores"].([]interface{})

		if !ok1 || !ok2 || len(keys) != len(scores) {
			c.Conn().WriteError(fmt.Sprintf(pkg.DeserializationErr, string(args[3])))
			return
		}

		for i := range keys {
			key, ok3 := keys[i].(string)
			score, ok4 := scores[i].(float64)

			if !ok3 || !ok4 {
				c.Conn().WriteError(fmt.Sprintf(pkg.DeserializationErr, string(args[3])))
				return
			}

			set.Inner.AddOrUpdate(key, score)
		}

		db.Set(key, set, ttl)
	}
	c.Conn().WriteString("OK")
}
