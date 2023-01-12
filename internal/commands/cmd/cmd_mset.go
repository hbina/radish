package cmd

import (
	"fmt"
	"time"

	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/types"
	"github.com/hbina/radish/internal/util"
)

// https://redis.io/commands/mset/
// MSET key value [key value ...]
func MsetCommand(c *pkg.Client, args [][]byte) {
	if len(args) < 3 || len(args)%2 != 1 {
		// If the number of arguments (excluding the command name) is not even,
		// return syntax error
		c.Conn().WriteError(fmt.Sprintf(util.WrongNumOfArgsErr, string(args[0])))
		return
	}

	db := c.Db()

	for i := 1; i < len(args); i += 2 {
		keyStr := string(args[i])
		valueStr := string(args[i+1])

		db.Set(keyStr, types.NewString(valueStr), time.Time{})
	}

	c.Conn().WriteString("OK")
}
