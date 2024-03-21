package cmd

import (
	"fmt"
	"strings"

	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/util"
)

// https://redis.io/commands/client-getname/
// https://redis.io/commands/client-setname/
func ClientCommand(c *pkg.Client, args [][]byte) {
	if len(args) < 2 {
		c.Conn().WriteError(fmt.Sprintf(util.WrongNumOfArgsErr, args[0]))
		return
	}

	subcommand := string(args[1])

	if strings.ToLower(subcommand) == "getname" {
		// Requires an extra argument for the name
		if len(args) != 2 {
			c.Conn().WriteError(fmt.Sprintf(util.WrongNumOfArgsErr, string(args[0])))
			return
		}

		if c.Name == nil {
			c.Conn().WriteNull()
			return
		} else {
			c.Conn().WriteBulkString(*c.Name)
			return
		}
	} else if strings.ToLower(subcommand) == "setname" {
		if len(args) < 3 {
			c.Conn().WriteError(fmt.Sprintf("Unknown subcommand or wrong number of arguments for '%s'. Try CONFIG HELP.", string(args[1])))
			return
		}

		newName := string(args[2])
		c.Name = &newName

		c.Conn().WriteString("OK")
		return
	} else {
		c.Conn().WriteError(fmt.Sprintf("Unknown subcommand '%s'. Try CONFIG HELP.", subcommand))
		return
	}
}
