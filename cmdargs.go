package main

import (
	"fmt"
	"log"
	"net"
	"runtime"
)

func parseArgs(args []string) (argset argumentSet) {
	var incomingIP string

	for index := 0; index < len(args)-1; index++ {
		incomingIP = args[index+1]
		argset.inputIP = net.ParseIP(incomingIP)

		if argset.inputIP != nil {
			return
		}
	}

	if len(args) <= 1 {
		fmt.Println("Use of go-ibc:\n\ngo-ibc IP-ADDRESS")
		fmt.Println("IP-ADDRESS can be an valid IPv4 or IPv6 Address\nCIDR Notation is not supported")
		fmt.Printf("\ngo-ibc: %s build: %s\n", version, runtime.Version())
	}

	if incomingIP == "" {
		log.Fatalln("Error: Not enough arguments given !")
	}

	return
}
