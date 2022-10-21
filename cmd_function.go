package redis

import (
	"fmt"

	"github.com/tidwall/redcon"
)

func FunctionCommand(c *Client, cmd redcon.Command) {
	for _, v := range cmd.Args {
		fmt.Printf("%s ", string(v))
	}
	fmt.Println()
	c.Conn().WriteNull()
}
