package cmd

import "github.com/hbina/radish/internal/pkg"

// https://redis.io/commands/zunionstore/
// ZUNIONSTORE destination numkeys key [key ...] [WEIGHTS weight [weight ...]] [AGGREGATE <SUM | MIN | MAX>]
func ZunionstoreCommand(c *pkg.Client, args [][]byte) {
	implZSetSetOperationCommand(c, args, true, ZSetOperationUnion, false)
}
