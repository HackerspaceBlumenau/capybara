package commands

import (
	"fmt"
	"log"
	"net/http"

	"github.com/hackerspaceblumenau/capybara/contracts"
)

type Server interface {
	contracts.Server
	Remember(w http.ResponseWriter, r *http.Request)
}

type RunOptions struct {
	Address string
	Port    int
}

type server struct {
	storage contracts.Storage
	options RunOptions
}

func (s server) registerHandlers() {
	http.HandleFunc("/commands/remember", s.Remember)
}

func (s server) Run() error {
	s.registerHandlers()
	addr := fmt.Sprintf("%s:%d", s.options.Address, s.options.Port)
	log.Printf("Starting command server at %s", addr)
	return http.ListenAndServe(addr, nil)
}

func NewServer(st contracts.Storage, opts RunOptions) Server {
	return &server{
		storage: st,
		options: opts,
	}
}
