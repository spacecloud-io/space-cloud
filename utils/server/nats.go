package server

import (
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"

	nats "github.com/nats-io/nats-server/v2/server"
)

// DefaultNatsOptions are the default setting to start nats with
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

// RunNatsServer starts a nats server in a separate goroutine
func (s *Server) RunNatsServer(seeds string, port, clusterPort int) error {
	// TODO read nats config from the yaml file if it exists
	if seeds != "" {
		array := strings.Split(seeds, ",")
		urls := []*url.URL{}
		for _, v := range array {
			if v != "" {
				u, err := url.Parse("nats://" + v + ":" + strconv.Itoa(clusterPort))
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

	fmt.Println("Starting Nats Server")
	go s.nats.Start()
	// Wait for accept loop(s) to be started
	if !s.nats.ReadyForConnections(10 * time.Second) {
		log.Fatal("Unable to start NATS Server in Go Routine")
	}
	return nil
}
