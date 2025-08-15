package pinger

/*
This package contains all the core logic for pinging hosts, collecting metrics, and handling concurrency.
All the functions and structs related to ICMP requests, metric storage, and the ping loop would live here.
Placing this in an internal directory signals that it's meant for use only by this project.
*/

import (
	"fmt"
	"log"
	"sync"
	"time"

	probing "github.com/prometheus-community/pro-bing"
)

// This will handle the pinging of multiple hosts and stores the data in a map
type PingManager struct {
	mu          sync.Mutex
	Hosts       []string
	Data        map[string]*probing.Statistics
	Interval    time.Duration
	Count       int
	RecentPings map[string][]bool
}

type HostMetrics struct {
	Stats      *probing.Statistics
	RecentLoss float64
}

func NewPingerManager(hosts []string, interval time.Duration, count int) *PingManager { // Constructor for PingManager
	return &PingManager{
		Hosts:       hosts,
		Data:        make(map[string]*probing.Statistics),
		Interval:    interval,
		Count:       count,
		RecentPings: make(map[string][]bool),
	}
}

func (pm *PingManager) StartPinging() {
	for _, host := range pm.Hosts {
		go pm.pingHost(host)
	}
}

func (pm *PingManager) GetMetrics() map[string]*HostMetrics {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	metricsCopy := make(map[string]*HostMetrics, len(pm.Data))
	for k, v := range pm.Data {
		recentLoss := pm.GetRecentLoss(k)
		metricsCopy[k] = &HostMetrics{
			Stats:      v,
			RecentLoss: recentLoss,
		}
	}

	return metricsCopy
}

func (pm *PingManager) GetRecentLoss(host string) float64 {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	history, ok := pm.RecentPings[host]
	if !ok || len(history) == 0 {
		return 0.0
	}

	lostCount := 0
	for _, success := range history {
		if !success {
			lostCount++
		}
	}
	return float64(lostCount) / float64(len(history)) * 100
}

func (pm *PingManager) pingHost(host string) {
	for {
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
			pm.mu.Lock()
			pm.RecentPings[host] = append(pm.RecentPings[host], true)
			if len(pm.RecentPings[host]) > 10 {
				pm.RecentPings[host] = pm.RecentPings[host][1:]
			}
			pm.Data[host] = pinger.Statistics()
			pm.mu.Unlock()
		}

		pinger.OnFinish = func(stats *probing.Statistics) {
			if stats.PacketLoss > 0 && stats.PacketsRecv < stats.PacketsSent {
				lostCount := stats.PacketsSent - stats.PacketsRecv
				pm.mu.Lock()
				for i := 0; i < int(lostCount); i++ {
					pm.RecentPings[host] = append(pm.RecentPings[host], false)
				}
				if len(pm.RecentPings[host]) > 10 {
					pm.RecentPings[host] = pm.RecentPings[host][len(pm.RecentPings[host])-10:]
				}
				pm.mu.Unlock()
			}
		}

		fmt.Printf("Pinging %s with an interval of %v\n", host, pm.Interval)
		err = pinger.Run()
		if err != nil {
			log.Printf("Pinger for host %s finished. Restarting in 5s...", host)
		}

		time.Sleep(5 * time.Second)
	}
}
