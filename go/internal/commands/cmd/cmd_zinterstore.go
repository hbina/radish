package cmd

import "github.com/hbina/radish/internal/pkg"

// https://redis.io/commands/zinterstore/
// ZINTERSTORE numkeys key [key ...] [WEIGHTS weight [weight ...]] [AGGREGATE <SUM | MIN | MAX>] [WITHSCORES]
func ZinterstoreCommand(c *pkg.Client, args [][]byte) {
	implZSetSetOperationCommand(c, args, true, ZSetOperationInter, false)
}
