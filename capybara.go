package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/eventials/goawaysignals"
	"github.com/hackerspaceblumenau/capybara/commands"
	"github.com/hackerspaceblumenau/capybara/contracts"
	"github.com/hackerspaceblumenau/capybara/slack"
	"github.com/hackerspaceblumenau/capybara/storage/drivers"
	"github.com/hackerspaceblumenau/capybara/ticker"
)

func printHelp() {
	msg := `Capybara (c) Hackerspace Blumenau Org.
-------------------------------------
Environment:
 SLACK_TOKEN: token used to connect to slack api
 SERVER_ADDR: address where server will listen to (default localhost)
 SERVER_PORT: port where server will listen to (default 8080)`

	fmt.Println(msg)
}

func main() {
	slackToken := os.Getenv("SLACK_TOKEN")
	if slackToken == "" {
		printHelp()
		return
	}

	commandsServerAddress := os.Getenv("SERVER_ADDR")
	if commandsServerAddress == "" {
		commandsServerAddress = "localhost"
	}

	commandsServerPort, _ := strconv.Atoi(os.Getenv("SERVER_PORT"))
	if commandsServerPort == 0 {
		commandsServerPort = 8080
	}

	slack.Setup(slackToken)

	st := drivers.MemoryStorage{}

	servers := []contracts.Server{
		commands.NewServer(
			st,
			commands.RunOptions{
				Address: commandsServerAddress,
				Port:    commandsServerPort,
			}),
		ticker.NewServer(st),
	}

	for _, srv := range servers {
		go srv.Run()
	}

	log.Println("Starting...")
	goawaysignals.Wait()
	log.Println("Exiting...")
}
