package main

import (
	"fmt"
	"strconv"
	cli2 "svpcc/cli"
	"svpcc/config"
	"svpcc/rest"
)

func main() {
	cfg := config.GetConfig()
	cli := cli2.NewReader(cfg.ServerAddress())

	loop:
	for {
		cmdArray, err := cli.Read()

		if err != nil {
			fmt.Println(err.Error())
			break loop
		}

		if len(cmdArray) >= 1 {
			cmd := cmdArray[0]

			switch cmd {
			case "quit", "exit":
				break loop
			case "get":
				if len(cmdArray) >= 2 {
					evaluateGetCommand(cmdArray, cfg)
				} else {
					fmt.Println("No commands found")
				}
			case "flush":
				evaluateFlushCommand(cmdArray, cfg)
			case "help":
				printHelp()
			default:
				fmt.Println("Command not found")
			}
		}

		fmt.Println(" ")
	}

	fmt.Println("Exiting...")
}

func printHelp() {
	var helpData = `
Commands:
---------

quit, exit - Exiting program

get buffers - Print data about filled buffers on remote application
              Mixed /data & /buffers calls and merge data

flush [all | {digit}] - Flushes remote buffers. If no args - flush all
flush, flush all - flushes all buffers
flush 1 - flush buffer with id = 1 etc.`

	fmt.Println(helpData)
}

func evaluateFlushCommand(strings []string, cfg config.Config) {
	var status string
	var err error

	if len(strings) > 1 {
		token := strings[1]

		if "all" == token {
			status, err = rest.Flush(cfg, -1)
		} else {
			var bufferID int

			bufferID, err = strconv.Atoi(token)

			if err != nil {
				fmt.Printf("Entered value must be \"all\" token or a digit. Actual: \"%s\".", token )
				return
			}

			status, err = rest.Flush(cfg, bufferID)
		}
	} else {
		status, err = rest.Flush(cfg, -1)
	}

	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Flush result: " + status)
	}

}

func evaluateGetCommand(strings []string, cfg config.Config) {
	subCommand := strings[1]

	switch subCommand {
	case "buffers":
		printBuffers(cfg)
	case "status":
		getStatus(cfg)
	default:
		fmt.Println("Command not found")
	}
}

func getStatus(cfg config.Config) {
	status, err := rest.GetStatus(cfg)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Server status: " + status)
}

func printBuffers(cfg config.Config) {
	data, err := rest.ReadBuffersData(cfg)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(" ")

	fmt.Println(" +--------------+-------------------------------+-------------------+ ")
	fmt.Println(" |      ID      |           File Name           |    Buffer Size    | ")
	fmt.Println(" +--------------+-------------------------------+-------------------+ ")

	for i := 0; i < len(data); i++ {
		entry := data[i]
		fmt.Println(fmt.Sprintf(" | %12s | %-29s | %17s |", entry.Id(), entry.FileName(), entry.BufferSize()))
	}

	fmt.Println(" +--------------+-------------------------------+-------------------+ ")
}
