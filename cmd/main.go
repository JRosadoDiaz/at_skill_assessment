package main

/*
This is the entry point for the program.
It should contain the main function, where you'll parse command-line flags, set up your services, and start the application.
This file should be kept as clean as possible, primarily acting as a coordinator.
*/

// go run AT_Skill_Assessment/cmd/main.go www.google.com,www.reddit.com world 3

import (
	"flag" // Parses command line arguments
	"fmt"
	"strings"
	"time"

	"github.com/JRosadoDiaz/AT_Skill_Assessment/internal/pinger"
	"github.com/JRosadoDiaz/AT_Skill_Assessment/web"
)

// Command-line flags
var hostsStr string
var port string
var interval time.Duration
var count int

func main() {
	// Grabbing flags
	flag.StringVar(&hostsStr, "hosts", "", "comma-seperated list of hosts to ping")
	flag.StringVar(&port, "port", "8000", "Port number")
	flag.DurationVar(&interval, "interval", time.Second*5, "The interval between pings")
	flag.IntVar(&count, "count", 5, "Number of times the host will be pinged")
	flag.Parse()

	if hostsStr == "" {
		fmt.Printf("Error: The 'hosts' flag is required.\nDid you forget to add '--hosts='?\n")
		flag.Usage()
		return
	}
	hosts := strings.Split(hostsStr, ",")

	// Generate ping manager
	var pingStruct = pinger.NewPingerManager(hosts, interval)
	pingStruct.StartPinging()

	// Start web server
	web.StartServer(port, pingStruct)

	// select {} // Will keep the project running indefinitely
}
