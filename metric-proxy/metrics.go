package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func (p *Proxy) collectMetrics() (*EnvoyMetrics, error) {
	logrus.Debugln("Pulling metrics from envoy...")
	res, err := http.Get("http://localhost:15000/stats?filter=(?=.*downstream_rq_total)(?=.*http.inbound)&format=json")
	if err != nil {
		return nil, err
	}
	defer closeReaderCloser(res.Body)

	data, _ := ioutil.ReadAll(res.Body)
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid response (status code: %d; body: %s) received from envoy", res.StatusCode, string(data))
	}

	metrics := new(EnvoyMetrics)
	if err := json.Unmarshal(data, metrics); err != nil {
		return nil, err
	}

	logrus.Debugln("Received metrics from envoy:", metrics)
	return metrics, nil
}

func (p *Proxy) routineCollectMetrics(duration time.Duration) {
	// This variable tracks the last req count
	prevValue := uint64(0)

	ticker := time.NewTicker(duration)
	for range ticker.C {
		metrics, err := p.collectMetrics()
		if err != nil {
			logrus.Errorln("Could not pull metrics from envoy:", err)
			continue
		}

		if len(metrics.Stats) == 0 {
			logrus.Errorln("Could not pull metrics. Is something wrong with envoy?")
			continue
		}

		// Calculate the number of requests which occurred between subsequent requests
		count := metrics.Stats[0].Value - prevValue
		prevValue = metrics.Stats[0].Value

		// Prepare and send proxy message
		message := &ProxyMessage{ActiveRequests: int32(count)}
		p.ch <- message
	}
}
