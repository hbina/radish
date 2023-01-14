package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"runtime"

	"github.com/hbina/radish/internal/util"
)

func main() {
	osArgs := os.Args[:]

	// Only positional arguments right now
	if len(osArgs) != 3 || osArgs[1] != "--port" {
		fmt.Println("Please provide the port")
		os.Exit(1)
	}

	port := osArgs[2]
	servAddr := fmt.Sprintf("localhost:%s", port)
	tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)

	if err != nil {
		fmt.Println("ResolveTCPAddr failed:", err.Error())
		os.Exit(1)
	}

	tcpConn, err := net.DialTCP("tcp", nil, tcpAddr)

	if err != nil {
		fmt.Println("Dial failed:", err.Error())
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
			fmt.Println("ReadString failed:", err.Error())
			continue
		}

		if len(inputStr) == 0 {
			continue
		}

		// Remove the last byte which is the newline
		// TODO: Check for other runtimes?
		// TODO: Validate the range here
		if runtime.GOOS == "windows" {
			if len(inputStr) <= 2 {
				continue
			}

			inputStr = inputStr[:len(inputStr)-2]
		} else {
			if len(inputStr) <= 1 {
				continue
			}

			inputStr = inputStr[:len(inputStr)-1]
		}

		args, valid := util.SplitStringIntoArgs(inputStr)

		if !valid || len(args) == 0 {
			continue
		}

		respInput := util.ConvertCommandArgToResp(args)

		_, err = tcpConn.Write([]byte(respInput))

		if err != nil {
			fmt.Println("Write to server failed:", err.Error())
			os.Exit(1)
		}

		response := make([]byte, 0, 1024)

		for {
			buffer := make([]byte, 1024)
			readCount, err := tcpConn.Read(buffer)

			if err != nil {
				if err != io.EOF {
					fmt.Println("Read from server failed:", err.Error())
					os.Exit(1)
				}
				break
			}

			buffer = buffer[:readCount]

			response = append(response, buffer...)

			responseDisplay, ok, _ := util.StringifyRespBytes(response)

			if ok {
				fmt.Println(util.EscapeString(string(response)))
				fmt.Println(responseDisplay)
				break
			}
		}
	}
}
