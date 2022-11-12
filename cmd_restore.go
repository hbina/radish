package redis

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/zavitax/sortedset-go"
)

// https://redis.io/commands/restore/
func RestoreCommand(c *Client, args [][]byte) {
	if len(args) == 0 {
		c.Conn().WriteError("no argument passed to handler. This should not be possible")
		return
	} else if len(args) < 4 {
		c.Conn().WriteError(fmt.Sprintf(WrongNumOfArgsErr, args[0]))
		return
	}

	key := string(args[1])
	ttl, err := ParseExpiryTime(string(args[2]), uint64(time.Millisecond))

	// Do not fail on time.Time{}
	if err != nil {
		c.Conn().WriteError(InvalidIntErr)
		return
	}

	var kvp Kvp
	err = json.Unmarshal(args[3], &kvp)

	if err != nil {
		log.Println(err)
		c.Conn().WriteError(fmt.Sprintf(DeserializationErr, string(args[3])))
		return
	}

	db := c.Db()

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
		set, ok := kvp.Value.(sortedset.SortedSet[string, float64, struct{}])

		if !ok {
			c.Conn().WriteError(fmt.Sprintf(DeserializationErr, string(args[3])))
			return
		}

		db.Set(key, NewZSetFromSortedSet(set), ttl)
	}

	c.Conn().WriteString("OK")
}
