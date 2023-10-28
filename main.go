package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"text/tabwriter"
)

var resultsWriter *tabwriter.Writer

func main() {
	// Define flags
	target := flag.String("target", "", "Checks a single IP, a range of IPs given CIDR notation, or a hostname")
	outputFile := flag.String("o", "", "Write output to file")
	// Parse flags
	flag.Parse()

	// Validate parameters
	if *target == "" {
		fmt.Println("target is required.")
		flag.PrintDefaults()
		return
	}
	var output *os.File
	if *outputFile == "" {
		output = os.Stdout
	} else {
		var err error
		output, err = os.Create(*outputFile)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	fmt.Printf("Scanning Target: %s\n", *target)
	fmt.Printf("Output File: %s\n", *outputFile)
	// Init tab writter for results
	resultsWriter = tabwriter.NewWriter(output, 0, 8, 0, '\t', 0)
	fmt.Fprintln(resultsWriter, "Target\tHost\tIP")
	// Check target type
	if strings.Contains(*target, "/") {
		CheckCIDR(*target)
	} else if ip := net.ParseIP(*target); ip != nil {
		CheckIp(ip)
	} else {
		CheckHost(*target)
	}
	resultsWriter.Flush()
}
func CheckHost(host string) {
	addrs, err := net.LookupHost(host)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, addr := range addrs {
		if ip := net.ParseIP(addr); ip != nil {
			CheckIp(ip)
		}
	}
}
func CheckIp(ip net.IP) {
	// fmt.Printf("Checking %s\n", ip.String())
	hosts, err := net.LookupAddr(ip.String())
	if err != nil {
		return
	}
	var foundIps string
	for i, host := range hosts {
		ips, err := net.LookupHost(host)
		if err != nil {
			return
		}
		if i == 0 {
			foundIps += strings.Join(ips, ", ")
		} else {
			foundIps += ", " + strings.Join(ips, ", ")
		}
	}
	fmt.Fprintf(resultsWriter, "%s\t%s\t%s\n", ip.String(), strings.Join(hosts, ", "), foundIps)
}
func CheckCIDR(target string) {
	start, network, err := net.ParseCIDR(target)
	if err != nil {
		fmt.Println(err)
		return
	}
	for ip := start; network.Contains(ip); ip = NextIp(ip) {
		CheckIp(ip)
	}
}
func NextIp(ip net.IP) net.IP {
	for i := len(ip) - 1; i >= 0; i-- {
		ip[i]++
		if ip[i] != 0 {
			break
		}
	}
	return ip
}
