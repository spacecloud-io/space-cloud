package server

import (
	"log"
	"time"

	nats "github.com/nats-io/nats-server/server"
)

var DefaultNatsOptions = &nats.Options{
	Host:   "0.0.0.0",
	Port:   4222,
	NoLog:  false,
	NoSigs: true,
	Cluster: nats.ClusterOpts{
		Host: "0.0.0.0",
		Port: 4248,
	},
}

func (s *Server) RunNatsServer(opts *nats.Options) {
	s.nats = nats.New(opts)
	go s.nats.Start()
	// Wait for accept loop(s) to be started
	if !s.nats.ReadyForConnections(10 * time.Second) {
		log.Fatal("Unable to start NATS Server in Go Routine")
	}
}
