package cmd

import "github.com/hbina/radish/internal/pkg"

// https://redis.io/commands/zintercard/
// ZINTERCARD numkeys key [key ...] [LIMIT limit]
func ZintercardCommand(c *pkg.Client, args [][]byte) {
	implZSetSetOperationCommand(c, args, false, ZSetOperationInter, true)
}
