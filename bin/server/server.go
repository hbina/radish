package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	redis "github.com/hbina/radish"
)

const (
	PORT = 6381
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Please specify the port to use")
		os.Exit(1)
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
					fmt.Println("port must be a number")
					os.Exit(1)
				}

				port = int(port64)
			}
		case "--no-log":
			{
				logging = false
			}
		default:
			{
				fmt.Printf("Unknown parameter '%s'\n", string(arg))
				os.Exit(1)
			}
		}
	}

	fmt.Printf("Starting server at port %d\n", port)
	redis.Run(port, logging)
}
