package cmd

import "github.com/hbina/radish/internal/pkg"

// https://redis.io/commands/zdiffcard/
// ZDIFFCARD numkeys key [key ...] [WITHSCORES]
func ZdiffcardCommand(c *pkg.Client, args [][]byte) {
	implZSetSetOperationCommand(c, args, false, ZSetOperationDiff, true)
}
