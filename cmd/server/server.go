package main

import (
	"net"
	"os"
	"strings"

	redis "github.com/hbina/radish"
)

const (
	HOST = "localhost"
	PORT = "9001"
	TYPE = "tcp"
)

func main() {

	if len(os.Args) < 2 {
		redis.Logger.Fatal("Please specify the port to use")
	}

	args := os.Args[1:]

	port := PORT

	for idx := 0; idx < len(args); idx++ {
		arg := strings.ToLower(string(args[idx]))

		switch arg {
		case "--port":
			{
				if idx+1 >= len(args) {
					redis.Logger.Fatal("Need to provide the port")
				}

				port = string(string(args[idx+1]))
				idx++
			}
		default:
			{
				redis.Logger.Fatalf("Unknown parameter '%s'\n", string(arg))
			}
		}
	}

	listen, err := net.Listen(TYPE, HOST+":"+port)

	if err != nil {
		redis.Logger.Fatal(err)
	}

	instance := redis.Default()

	for {
		conn, err := listen.Accept()

		if err != nil {
			redis.Logger.Fatal(err)
		}

		go handleIncomingRequest(instance, instance.NewClient(conn))
	}
}

func handleIncomingRequest(r *redis.Redis, client *redis.Client) {
	buffer := make([]byte, 0, 1024)
	tmp := make([]byte, 1024)
	count, err := client.Read(tmp)

	if err != nil {
		redis.Logger.Fatal(err)
	}

	for {
		buffer = append(buffer, tmp[:count]...)

		// Try to parse the current buffer as a RESP
		resp, leftover := redis.ConvertBytesToRespType(buffer)

		if resp != nil {
			redis.Logger.Println(redis.EscapeString(string(buffer)))
			buffer = leftover
			r.HandleRequest(client, redis.ConvertRespToArgs(resp))
		}

		count, err = client.Read(tmp)

		if err != nil || count == 0 {
			return
		}
	}
}
