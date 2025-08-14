package pinger

/*
This package contains all the core logic for pinging hosts, collecting metrics, and handling concurrency.
All the functions and structs related to ICMP requests, metric storage, and the ping loop would live here.
Placing this in an internal directory signals that it's meant for use only by this project.
*/

import (
	"fmt"
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
	Count    int
}

func NewPingerManager(hosts []string, interval time.Duration, count int) *PingManager { // Constructor for PingManager
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

func (pm *PingManager) GetMetrics() map[string]*probing.Statistics {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	metricsCopy := make(map[string]*probing.Statistics, len(pm.Data))
	for k, v := range pm.Data {
		metricsCopy[k] = v
	}

	return metricsCopy
}

func (pm *PingManager) pingHost(host string) {
	pinger, err := probing.NewPinger(host)
	if err != nil {
		fmt.Printf("Error creating pinger for %s: %v\n", host, err)
		return
	}

	// Configure what the pinger will do
	pinger.Interval = pm.Interval
	pinger.Timeout = time.Second * 30
	pinger.Count = pm.Count

	// Print status whenever ping is recieved
	pinger.OnRecv = func(pkt *probing.Packet) {
		fmt.Printf("Recieved ping replay from %s: bytes=%d time=%v ttl=%d\n", pkt.IPAddr, pkt.Nbytes, pkt.Rtt, pkt.TTL)
	}

	// Records what happens when it finishes a ping
	pinger.OnFinish = func(stats *probing.Statistics) {
		fmt.Println("Done")
		pm.mu.Lock()
		pm.Data[host] = stats
		pm.mu.Unlock()
		// fmt.Println(pm.Data[host].PacketsRecv)
	}

	fmt.Printf("Pinging %s with an interval of %v\n", host, pm.Interval)
	pinger.Run()
}
