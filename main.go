package main

import "fmt"

type IP struct {
	ipv4 string
	ipv6 string
}

func main() {
	input := IP{}
	parseArgs(&input)
	fmt.Println("IPv4: ", input.ipv4)
	fmt.Println("Ipv6: ", input.ipv6)

	initDownload()
}
