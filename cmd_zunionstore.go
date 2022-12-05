package redis

// https://redis.io/commands/zunionstore/
// ZUNIONSTORE destination numkeys key [key ...] [WEIGHTS weight [weight ...]] [AGGREGATE <SUM | MIN | MAX>]
func ZunionstoreCommand(c *Client, args [][]byte) {
	implZSetSetOperationCommand(c, args, true, ZSetOperationUnion, false)
}
