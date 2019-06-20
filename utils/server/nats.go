package server

import (
	"log"
	"net/url"
	"strings"
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

func (s *Server) RunNatsServer(seeds string, port, clusterPort int) error {
	// TODO read nats config from the yaml file if it exists
	if seeds != "" {
		array := strings.Split(seeds, ",")
		urls := []*url.URL{}
		for _, v := range array {
			if v != "" {
				u, err := url.Parse("nats://" + v)
				if err != nil {
					return err
				}
				urls = append(urls, u)
			}
		}
		DefaultNatsOptions.Routes = urls
	}
	DefaultNatsOptions.Port = port
	DefaultNatsOptions.Cluster.Port = clusterPort

	s.nats = nats.New(DefaultNatsOptions)

	go s.nats.Start()
	// Wait for accept loop(s) to be started
	if !s.nats.ReadyForConnections(10 * time.Second) {
		log.Fatal("Unable to start NATS Server in Go Routine")
	}
	return nil
}
