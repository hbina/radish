package redis

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/hbina/radish/internal/commands"
	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/util"
)

var started bool = false

func Run(port int, shouldLog bool) {

	if started {
		return
	}

	started = true

	if shouldLog {
		util.Logger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	} else {
		util.Logger = &util.StubLogger{}
	}

	listen, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))

	if err != nil {
		util.Logger.Fatal(err)
	}

	instance := pkg.Default(
		commands.GenerateCommands(),
		commands.GenerateBlockingCommands(),
		commands.GenerateConfigs())

	for {
		conn, err := listen.Accept()

		if err != nil {
			util.Logger.Fatal(err)
		}

		go instance.HandleClient(instance.NewClient(conn))
	}
}
