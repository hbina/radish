package cmd

import (
	"fmt"
	"strconv"

	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/util"
)

// https://redis.io/commands/select/
func SelectCommand(c *pkg.Client, args [][]byte) {
	if len(args) == 1 {
		c.Conn().WriteError(fmt.Sprintf(util.WrongNumOfArgsErr, args[0]))
		return
	}

	index, err := strconv.ParseUint(string(args[1]), 10, 32)

	if err != nil {
		c.Conn().WriteError(util.InvalidIntErr)
	} else {
		c.Db().Unlock()
		c.SetDb(index)
		c.Db().Lock()
		c.Conn().WriteString("OK")
	}
}
