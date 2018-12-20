package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/eventials/goawaysignals"
	"github.com/hackerspaceblumenau/capybara/slack"
	"github.com/hackerspaceblumenau/capybara/storage/drivers"
	"github.com/hackerspaceblumenau/capybara/contracts"
	"github.com/hackerspaceblumenau/capybara/commands"
	"github.com/hackerspaceblumenau/capybara/ticker"
)

var (
	slackToken            = flag.String("token", "", "Slack token to use for integration")
	commandsServerAddress = flag.String("address", "0.0.0.0", "Command server address listener")
	commandsServerPort    = flag.Int("port", 8080, "Command server port listener")
)

func printHelp() {
	msg := `
Capybara (c) Hackerspace Blumenau Org.
-------------------------------------`

	fmt.Println(msg)
	flag.Usage()
}

func invalidFlags() bool {
	if *slackToken == "" {
		return true
	}

	if *commandsServerAddress == "" {
		return true
	}

	if *commandsServerPort <= 0 {
		return true
	}

	return false
}

func main() {
	flag.Parse()

	if invalidFlags() {
		printHelp()
		return
	}

    slack.Setup(*slackToken)

    st := drivers.MemoryStorage{}

	servers := []contracts.Server{
		commands.NewServer(
            st,
			commands.RunOptions{
				Address: *commandsServerAddress,
				Port:    *commandsServerPort,
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
