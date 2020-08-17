package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"time"

	"github.com/sirupsen/logrus"
)

func (p *Proxy) collectMetrics() (*EnvoyMetrics, error) {
	res, err := p.client.Get("http://localhost:15000/stats?format=json")
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

	var array []EnvoyStat
	for _, metric := range metrics.Stats {
		if metric.Value != nil && validMetric(metric.Name, p.filter) {
			logrus.Debugln("Received metrics from envoy:", metric.Name, metric.Value)
			array = append(array, metric)
		}
	}
	metrics.Stats = array
	return metrics, nil
}

func (p *Proxy) routineCollectMetrics(duration time.Duration) {
	// This variable tracks the last req count
	prevValue := float64(0)

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

		// Calculate the total value
		var value float64
		for _, stat := range metrics.Stats {
			num, ok := stat.Value.(float64)
			if !ok {
				logrus.Warningln("Unable to convert value to float64:", stat.Value, reflect.TypeOf(stat.Value))
			}
			value += num
		}

		// Calculate the number of requests which occurred between subsequent requests
		count := value - prevValue
		prevValue = value

		// For active requests we need to send the active request value straight away
		if p.filter == "downstream_rq_active" {
			count = value
		}

		// Make sure count is not zero
		if count < 0 {
			count = 0
		}

		// Prepare and send proxy message
		message := &ProxyMessage{ActiveRequests: int32(count)}
		p.ch <- message
	}
}
