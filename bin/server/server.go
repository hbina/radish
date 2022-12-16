package main

import (
	"log"
	"os"
	"strconv"
	"strings"

	redis "github.com/hbina/radish"
	"github.com/hbina/radish/internal/util"
)

const (
	PORT = 6381
)

func main() {

	if len(os.Args) < 2 {
		util.Logger.Fatal("Please specify the port to use")
	}

	args := os.Args[1:]

	port := PORT
	logging := true

	for idx := 0; idx < len(args); idx++ {
		arg := strings.ToLower(string(args[idx]))

		switch arg {
		case "--port":
			{
				idx++
				if idx >= len(args) {
					log.Fatal("Need to provide the port")
				}

				port64, err := strconv.ParseInt(string(args[idx]), 10, 32)

				if err != nil {
					log.Fatal("port must be a number")
				}

				port = int(port64)
			}
		case "--no-log":
			{
				logging = false
			}
		default:
			{
				log.Fatalf("Unknown parameter '%s'\n", string(arg))
			}
		}
	}

	redis.Run(port, logging)
}
