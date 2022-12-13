package redis

import (
	"fmt"
	"strings"
)

// https://redis.io/commands/config-get/
// https://redis.io/commands/config-set/
func ConfigCommand(c *Client, args [][]byte) {
	if len(args) < 2 {
		c.Conn().WriteError(fmt.Sprintf(WrongNumOfArgsErr, args[0]))
		return
	}

	subcommand := string(args[1])

	if strings.ToLower(subcommand) == "get" {
		if len(args) < 3 {
			c.Conn().WriteError(fmt.Sprintf(WrongNumOfArgsErr, string(args[0])))
			return
		}

		result := make([]string, 0, (len(args)-2)*2)
		for i := 2; i < len(args); i++ {
			k := string(args[i])
			v := c.Redis().GetConfigValue(k)
			if v != nil {
				result = append(result, k, *v)
			}
		}
		c.Conn().WriteArray(len(result))
		for _, v := range result {
			c.Conn().WriteBulkString(v)
		}
	} else if strings.ToLower(subcommand) == "set" {
		if len(args) < 4 {
			c.Conn().WriteError(fmt.Sprintf("Unknown subcommand or wrong number of arguments for '%s'. Try CONFIG HELP.", string(args[1])))
			return
		}

		k := string(args[2])
		v := string(args[3])

		c.Redis().SetConfigValue(k, v)

		c.Conn().WriteString("OK")
	} else {
		c.Conn().WriteError(fmt.Sprintf("Unknown subcommand '%s'. Try CONFIG HELP.", subcommand))
	}
}
