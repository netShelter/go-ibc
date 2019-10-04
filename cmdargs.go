package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"runtime"
)

func parseArgs(input *ipset) {
	var v4, v6 string
	flag.StringVar(&v4, "ipv4", "", "IPv4 Address")
	flag.StringVar(&v6, "ipv6", "", "IPv6 Address")
	flag.Parse()
	if len(os.Args) <= 1 {
		flag.PrintDefaults()
		fmt.Println("go-ibc version: ", version)
		fmt.Println("golang compile version: ", runtime.Version())
	}
	if v4 == "" && v6 == "" {
		log.Fatalln("Error: Not enough arguments given !")
	}

	input.ipv4 = net.ParseIP(v4)
	input.ipv6 = net.ParseIP(v6)

	if input.ipv4 == nil && v4 != "" {
		log.Fatalln("Error: IPv4 Address invalid")
	}

	if input.ipv6 == nil && v6 != "" {
		log.Fatalln("Error: IPv6 Address invalid")
	}
}
