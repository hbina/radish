package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/types"
)

// https://redis.io/commands/dump/
func DumpCommand(c *pkg.Client, args [][]byte) {
	if len(args) < 2 {
		c.Conn().WriteError(fmt.Sprintf(pkg.WrongNumOfArgsErr, args[0]))
		return
	}

	key := string(args[1])
	value, _ := c.Db().GetOrExpire(key, true)

	if value == nil {
		c.Conn().WriteNull()
		return
	}

	if value.Type() == types.ValueTypeString {
		str, err := json.Marshal(pkg.Kvp{
			Key:   key,
			Type:  value.TypeFancy(),
			Value: value.Value(),
		})

		if err != nil {
			c.Conn().WriteError(err.Error())
			return
		}

		c.Conn().WriteBulkString(string(str))

		return
	} else if value.Type() == types.ValueTypeList {
		arr := make([]string, 0)

		value.(*types.List).ForEachF(func(a string) {
			arr = append(arr, a)
		})

		str, err := json.Marshal(pkg.Kvp{
			Key:   key,
			Type:  value.TypeFancy(),
			Value: arr,
		})

		if err != nil {
			c.Conn().WriteError(err.Error())
			return
		}

		c.Conn().WriteBulkString(string(str))

		return
	} else if value.Type() == types.ValueTypeSet {
		str, err := json.Marshal(pkg.Kvp{
			Key:   key,
			Type:  value.TypeFancy(),
			Value: value.(*types.Set).Inner,
		})

		if err != nil {
			c.Conn().WriteError(err.Error())
			return
		}

		c.Conn().WriteBulkString(string(str))

		return
	} else if value.Type() == types.ValueTypeZSet {
		keys := make([]string, 0)
		scores := make([]float64, 0)

		for key, node := range value.(*types.ZSet).Inner.Dict {
			keys = append(keys, key)
			scores = append(scores, node.Score)
		}

		pair := types.SerdeZSet{
			Keys:   keys,
			Scores: scores,
		}

		str, err := json.Marshal(pkg.Kvp{
			Key:   key,
			Type:  value.TypeFancy(),
			Value: pair,
		})

		if err != nil {
			c.Conn().WriteError(err.Error())
			return
		}

		c.Conn().WriteBulkString(string(str))

		return
	}

	c.Conn().WriteError(fmt.Sprintf("Dump for %s is not yet implemented", value.TypeFancy()))
}
