package main

import (
	"net"
)

type ipset struct {
	ipv4 net.IP
	ipv6 net.IP
}

func main() {
	input := ipset{}
	parseArgs(&input)
	_, dir := initDownload()
	checkIPs(dir, &input)
}
