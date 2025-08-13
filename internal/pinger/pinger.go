package pinger

/*
This package contains all the core logic for pinging hosts, collecting metrics, and handling concurrency.
All the functions and structs related to ICMP requests, metric storage, and the ping loop would live here.
Placing this in an internal directory signals that it's meant for use only by this project.
*/

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	probing "github.com/prometheus-community/pro-bing"
)

// This will handle the pinging of multiple hosts and stores the data in a map
type PingManager struct {
	mu       sync.Mutex
	Hosts    []string
	Data     map[string]*probing.Statistics
	Interval time.Duration
}

func NewPingerManager(hosts []string, interval time.Duration) *PingManager {
	return &PingManager{
		Hosts:    hosts,
		Data:     make(map[string]*probing.Statistics),
		Interval: interval,
	}
}

func (pm *PingManager) StartPinging() {
	for _, host := range pm.Hosts {
		go pm.pingHost(host)
	}
}

func (pm *PingManager) GetMetrics(w http.ResponseWriter, r *http.Request) map[string]*probing.Statistics {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	metricsCopy := make(map[string]*probing.Statistics, len(pm.Data))
	for k, v := range pm.Data {
		metricsCopy[k] = v
	}

	// w.Header().Set("Content-Type", "application/json")
	// err = json.NewDecoder(w).Encode(response)
	// if err != nil {
	// 	api.InternalErrorHandler(w)
	// 	return
	// }
	return metricsCopy
}

func (pm *PingManager) pingHost(host string) {
	pinger, err := probing.NewPinger(host)
	if err != nil {
		fmt.Printf("Error creating pinger for %s: %v\n", host, err)
		return
	}

	fmt.Println("Starting pinger")

	// Configure what the pinger will do
	pinger.Interval = pm.Interval
	pinger.Timeout = time.Second * 1
	pinger.Count = 3

	// Records what happens when it finishes a ping
	pinger.OnFinish = func(stats *probing.Statistics) {
		pm.mu.Lock()
		pm.Data[host] = stats
		pm.mu.Unlock()
	}

	fmt.Printf("Pinging %s with an interval of %v\n", host, pm.Interval)
	pinger.Run()
}
