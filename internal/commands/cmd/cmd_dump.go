package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/types"
	"github.com/hbina/radish/internal/util"
)

// https://redis.io/commands/dump/
func DumpCommand(c *pkg.Client, args [][]byte) {
	if len(args) < 2 {
		c.Conn().WriteError(fmt.Sprintf(util.WrongNumOfArgsErr, args[0]))
		return
	}

	key := string(args[1])
	value, _ := c.Db().Get(key)

	if value == nil {
		c.Conn().WriteNull()
		return
	}

	if value.Type() == types.ValueTypeString {
		data, err := value.(*types.String).Marshal()

		if err != nil {
			c.Conn().WriteError(err.Error())
			return
		}

		str, err := json.Marshal(pkg.Kvp{
			Key:  key,
			Type: value.TypeFancy(),
			Data: data,
		})

		if err != nil {
			c.Conn().WriteError(err.Error())
			return
		}

		c.Conn().WriteBulkString(string(str))

		return
	} else if value.Type() == types.ValueTypeList {
		data, err := value.(*types.List).Marshal()

		if err != nil {
			c.Conn().WriteError(err.Error())
			return
		}

		str, err := json.Marshal(pkg.Kvp{
			Key:  key,
			Type: value.TypeFancy(),
			Data: data,
		})

		if err != nil {
			c.Conn().WriteError(err.Error())
			return
		}

		c.Conn().WriteBulkString(string(str))

		return
	} else if value.Type() == types.ValueTypeSet {
		data, err := value.(*types.Set).Marshal()

		if err != nil {
			c.Conn().WriteError(err.Error())
			return
		}

		str, err := json.Marshal(pkg.Kvp{
			Key:  key,
			Type: value.TypeFancy(),
			Data: data,
		})

		if err != nil {
			c.Conn().WriteError(err.Error())
			return
		}

		c.Conn().WriteBulkString(string(str))

		return
	} else if value.Type() == types.ValueTypeZSet {
		data, err := value.(*types.ZSet).Marshal()

		if err != nil {
			c.Conn().WriteError(err.Error())
			return
		}

		str, err := json.Marshal(pkg.Kvp{
			Key:  key,
			Type: value.TypeFancy(),
			Data: data,
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
