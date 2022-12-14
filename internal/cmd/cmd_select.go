package cmd

import (
	"fmt"
	"strconv"

	"github.com/hbina/radish/internal/pkg"
)

// https://redis.io/commands/select/
func SelectCommand(c *pkg.Client, args [][]byte) {
	if len(args) == 1 {
		c.Conn().WriteError(fmt.Sprintf(pkg.WrongNumOfArgsErr, args[0]))
		return
	}

	index, err := strconv.ParseUint(string(args[1]), 10, 32)

	if err != nil {
		c.Conn().WriteError(pkg.InvalidIntErr)
	} else {
		c.SetDb(index)
		c.Conn().WriteString("OK")
	}
}
