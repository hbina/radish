package redis

// https://redis.io/commands/zinterstore/
// ZINTERSTORE numkeys key [key ...] [WEIGHTS weight [weight ...]] [AGGREGATE <SUM | MIN | MAX>] [WITHSCORES]
func ZinterstoreCommand(c *Client, args [][]byte) {
	implZSetSetOperationCommand(c, args, true, ZSetOperationInter, false)
}
