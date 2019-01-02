package main

import (
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

func main() {
	slackToken := os.Getenv("SLACK_TOKEN")
	if slackToken == "" {
		log.Fatal("SLACK_TOKEN environment variable is empty")
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
