package cmd

import "github.com/hbina/radish/internal/pkg"

// https://redis.io/commands/zinter/
// ZINTER numkeys key [key ...] [WEIGHTS weight [weight ...]] [AGGREGATE <SUM | MIN | MAX>] [WITHSCORES]
func ZinterCommand(c *pkg.Client, args [][]byte) {
	implZSetSetOperationCommand(c, args, false, ZSetOperationInter, false)
}
