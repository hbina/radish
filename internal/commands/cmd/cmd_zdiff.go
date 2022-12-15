package cmd

import "github.com/hbina/radish/internal/pkg"

// https://redis.io/commands/zdiff/
// ZDIFF numkeys key [key ...] [WITHSCORES]
func ZdiffCommand(c *pkg.Client, args [][]byte) {
	implZSetSetOperationCommand(c, args, false, ZSetOperationDiff, false)
}
