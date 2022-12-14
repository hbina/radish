package cmd

import "github.com/hbina/radish/internal/pkg"

// https://redis.io/commands/zdiffstore/
// ZDIFFSTORE numkeys key [key ...] [WITHSCORES]
func ZdiffstoreCommand(c *pkg.Client, args [][]byte) {
	implZSetSetOperationCommand(c, args, true, ZSetOperationDiff, false)
}
