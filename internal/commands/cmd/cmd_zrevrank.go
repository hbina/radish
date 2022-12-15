package cmd

import "github.com/hbina/radish/internal/pkg"

// https://redis.io/commands/zrevrank/
// ZREVRANK key member WITHSCORE
func ZrevrankCommand(c *pkg.Client, args [][]byte) {
	implZrankCommand(c, args, true)
}
