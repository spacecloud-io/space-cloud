package main

import (
	"log"
	"time"

	nats "github.com/nats-io/nats-server/server"
)

var defaultNatsOptions = &nats.Options{
	Host:   "0.0.0.0",
	Port:   4222,
	NoLog:  true,
	NoSigs: true,
	Cluster: nats.ClusterOpts{
		Host: "0.0.0.0",
		Port: 4248,
	},
}

func (s *server) runNatsServer(opts *nats.Options) {
	s.nats = nats.New(opts)
	go s.nats.Start()
	// Wait for accept loop(s) to be started
	if !s.nats.ReadyForConnections(10 * time.Second) {
		log.Fatal("Unable to start NATS Server in Go Routine")
	}
}
