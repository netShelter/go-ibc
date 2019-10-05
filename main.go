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
	bsdir, dir := initDownload()
	checkIPs(bsdir, dir, &input)
}
