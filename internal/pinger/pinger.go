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
		Count:    count,
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
		panic(err)
	}

	// Configure what the pinger can do when active
	pinger.Interval = pm.Interval
	pinger.Count = pm.Count
	pinger.SetPrivileged(true)

	// Print which host is being pinged
	pinger.OnSend = func(p *probing.Packet) {
		fmt.Printf("Sending packet to %v...\n", host)
	}

	// Print status whenever ping is recieved
	pinger.OnRecv = func(pkt *probing.Packet) {
		fmt.Printf("Recieved ping replay from %s: bytes=%d time=%v ttl=%d\n", pkt.IPAddr, pkt.Nbytes, pkt.Rtt, pkt.TTL)
	}

	// Records what happens when a ping finishes
	pinger.OnFinish = func(stats *probing.Statistics) {
		pm.mu.Lock()
		pm.Data[host] = stats
		pm.mu.Unlock()
	}

	fmt.Printf("Pinging %s with an interval of %v\n", host, pm.Interval)
	err = pinger.Run()
	if err != nil {
		panic(err)
	}
}
