package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
	"utils/utils"
)

type ipWithRtt struct {
	ipDetails *utils.IpDetails
	ip        string
	rtt       time.Duration
}

func main() {
	file, err := os.Open("ips.txt")
	defer file.Close()

	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	counter := 1
	var ipWithRttSlice []ipWithRtt

	for scanner.Scan() {
		ipRange := scanner.Text()
		splitIpRange := strings.Split(scanner.Text(), "/")
		ip := splitIpRange[0]
		ipDetails, rtt, err := utils.IsOkIp(ip)
		fmt.Printf("%d- ", counter)
		counter++
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Printf("range: %s, country: %s, isp: %s, rtt: %s\n", ipRange, ipDetails.Country, ipDetails.Isp, rtt)
		ipWithRttSlice = append(ipWithRttSlice, ipWithRtt{
			ip:        ipRange,
			rtt:       rtt,
			ipDetails: ipDetails,
		})
	}
	sort.Slice(ipWithRttSlice, func(i, j int) bool {
		return ipWithRttSlice[i].rtt < ipWithRttSlice[j].rtt
	})

	file, err = os.OpenFile("OKRangesWithRtt.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	for _, i := range ipWithRttSlice {
		_, err := file.WriteString(fmt.Sprintf("range: %s, country: %s, isp: %s, rtt: %s\n", i.ip, i.ipDetails.Country, i.ipDetails.Isp, i.rtt))
		if err != nil {
			return
		}
	}
	file, err = os.OpenFile("OKRanges.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	for _, i := range ipWithRttSlice {
		_, err := file.WriteString(fmt.Sprintf("%s\n", i.ip))
		if err != nil {
			return
		}
	}

}
