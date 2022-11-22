package main

import (
	"fmt"
	"log"
	"os"

	redis "github.com/hbina/radish"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Please specify the port to use")
		os.Exit(1)
	}
	log.Fatal(redis.Run(fmt.Sprintf(":%s", os.Args[1])))
	os.Exit(0)
}
