package main

import (
	"fmt"
	"net"
)

type ipset struct {
	ipv4 net.IP
	ipv6 net.IP
}

func main() {
	input := ipset{}
	parseArgs(&input)
	fmt.Println("IPv4: ", input.ipv4)
	fmt.Println("Ipv6: ", input.ipv6)
	initDownload()
}
