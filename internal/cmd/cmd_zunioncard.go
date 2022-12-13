package cmd

import "github.com/hbina/radish/internal/pkg"

// https://redis.io/commands/zunioncard/
// ZUNIONCARD numkeys key [key ...] [LIMIT limit]
func ZunioncardCommand(c *pkg.Client, args [][]byte) {
	implZSetSetOperationCommand(c, args, false, ZSetOperationUnion, true)
}
