package main

import (
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

// Proxy is the module which collects metrics from envoy and pushes it to the autoscaler
type Proxy struct {
	addr, token, filter string

	// For communication
	c  *websocket.Conn
	ch chan *ProxyMessage
}

// New creates a new proxy instance
func New(addr, token, mode string) *Proxy {
	filter := "downstream_rq_total"
	if mode == "parallel" {
		filter = "downstream_rq_active"
	}
	return &Proxy{addr: addr, token: token, filter: filter, ch: make(chan *ProxyMessage, 1)}
}

// Start begins the metric collection operation
func (p *Proxy) Start() error {

	if err := p.connect(); err != nil {
		return err
	}

	// Start a ticker to push ping messages to server
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		for range ticker.C {
			if err := p.c.WriteMessage(websocket.PingMessage, nil); err != nil {
				logrus.Errorf("Could not write ping message to server - %s", err.Error())
				_ = p.connect()
			}
		}
	}()

	// Start the metric collection routine
	logrus.Infoln("Starting metric collection operation")
	go p.routineCollectMetrics(1 * time.Second)

	// Start infinite loop to push messages to autoscaler
	for msg := range p.ch {
		logrus.Debugln("Sending metrics to runner:", msg)
		if err := p.c.WriteJSON(msg); err != nil {
			logrus.Errorf("Could not write message to server - %s", err.Error())
			_ = p.connect()
		}
	}

	return errors.New("loop prematurely exited")
}

func (p *Proxy) connect() error {
	logrus.Debugf("Attempting websocket connection with %s", p.addr)
	u := url.URL{Scheme: "ws", Host: p.addr, Path: "/v1/runner/socket"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), http.Header{"Authorization": []string{"Bearer " + p.token}})
	if err != nil {
		return err
	}

	p.c = c
	logrus.Debugf("Established websocket connection with %s", p.addr)
	return nil
}
