package main

import (
	"flag"
	"log"
	"os"
)

func parseArgs(input *IP) {
	flag.StringVar(&input.ipv4, "ipv4", "", "IPv4 Address")
	flag.StringVar(&input.ipv6, "ipv6", "", "IPv6 Address")
	flag.Parse()
	if len(os.Args) <= 1 {
		flag.PrintDefaults()
	}
	if input.ipv4 == "" && input.ipv6 == "" {
		log.Fatalln("Error: Not enough arguments given !")
	}
}
