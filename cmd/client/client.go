package main

import (
	"bufio"
	"fmt"
	"go-redis"
	"log"
	"net"
	"os"
	"os/signal"
)

var (
	InfoLogger *log.Logger
)

func main() {

	InfoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

	osArgs := os.Args[:]

	InfoLogger.Println(osArgs)

	// Only positional arguments right now
	if len(osArgs) < 2 {
		InfoLogger.Println("Please provide the port")
		os.Exit(1)
	}

	port := osArgs[1]
	// strEcho := "*1\r\n$4\r\nPING\r\n"
	servAddr := fmt.Sprintf("localhost:%s", port)
	tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)

	if err != nil {
		InfoLogger.Println("ResolveTCPAddr failed:", err.Error())
		os.Exit(1)
	}

	tcpConn, err := net.DialTCP("tcp", nil, tcpAddr)

	if err != nil {
		InfoLogger.Println("Dial failed:", err.Error())
		os.Exit(1)
	}

	{
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		go func() {
			for range c {
				tcpConn.Close()
				os.Exit(1)
			}
		}()
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		inputStr, err := reader.ReadString('\n')

		if err != nil {
			InfoLogger.Println("ReadString failed:", err.Error())
			continue
		}

		// Remove the last byte which is the newline
		inputStr = inputStr[:len(inputStr)-1]

		InfoLogger.Println("arg", inputStr)

		if len(inputStr) == 0 {
			continue
		}

		args, valid := redis.SplitStringIntoArgs(inputStr)

		InfoLogger.Println("args", args)

		if !valid || len(args) == 0 {
			continue
		}

		respInput := redis.ConvertCommandArgToResp(args)

		InfoLogger.Println("respStr", redis.EscapeString(respInput))

		_, err = tcpConn.Write([]byte(respInput))

		if err != nil {
			InfoLogger.Println("Write to server failed:", err.Error())
			os.Exit(1)
		}

		respOutput := make([]byte, 1024)

		// TODO: Fix this so that it actually reads all the output by looping it until EOF
		readCount, err := tcpConn.Read(respOutput)

		InfoLogger.Printf("respOutput '%s'", redis.EscapeString(string(respOutput)))

		if err != nil {
			InfoLogger.Println("Write to server failed:", err.Error())
			os.Exit(1)
		}

		res, leftover := redis.CreateRespReply(respOutput[:readCount])

		if len(leftover) != 0 {
			InfoLogger.Println("CreateRespReply have leftover", leftover)
		}

		fmt.Println(res)
	}
}
