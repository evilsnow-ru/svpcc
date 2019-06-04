package main

import (
	"fmt"
	cli2 "svpcc/cli"
	"svpcc/config"
	"svpcc/rest"
)

func main() {
	fmt.Println("Loading config...")
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
			default:
				fmt.Println("Command not found")
			}
		}

	}

	fmt.Println("Exiting...")
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
	fmt.Println(" ")
}

func printBuffers(cfg config.Config) {
	data, err := rest.ReadBuffersData(cfg)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(" ")

	fmt.Println(" +--------------+-------------------------+-------------------+ ")
	fmt.Println(" |      ID      |        File Name        |    Buffer Size    | ")
	fmt.Println(" +--------------+-------------------------+-------------------+ ")

	for i := 0; i < len(data); i++ {
		entry := data[i]
		fmt.Println(fmt.Sprintf(" | %12s | %-23s | %17s |", entry.Id(), entry.FileName(), entry.BufferSize()))
	}

	fmt.Println(" +--------------+-------------------------+-------------------+ ")
	fmt.Println(" ")
}
