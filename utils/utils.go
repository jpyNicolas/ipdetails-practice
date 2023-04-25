package utils

import (
	"encoding/json"
	"fmt"
	probing "github.com/prometheus-community/pro-bing"
	"io/ioutil"
	"net/http"
	"time"
)

type IpDetails struct {
	Query, Status, Country, Isp string
}

func GetIpDetails(ip string) (*IpDetails, error) {
	resp, err := http.Get(fmt.Sprintf("http://ip-api.com/json/%s?fields=status,country,query,isp", ip))
	defer resp.Body.Close()

	statusCode := resp.StatusCode

	if statusCode != 200 {
		return nil, fmt.Errorf("%d", statusCode)
	}

	if err != nil {
		return nil, fmt.Errorf("faild to call API")
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("faild to read response body")
	}

	var ipDetails IpDetails

	if err := json.Unmarshal(body, &ipDetails); err != nil {
		return nil, fmt.Errorf("faild to read Ip Details")
	}

	return &ipDetails, nil

}

func PingRtt(ip string) (time.Duration, error) {
	pinger, err := probing.NewPinger(ip)
	pinger.SetPrivileged(true)
	if err != nil {
		return 0, err
	}
	pinger.Timeout = time.Duration(time.Second)
	pinger.Count = 2

	err = pinger.Run() // Blocks until finished.
	if err != nil {
		return 0, err
	}
	return pinger.Statistics().AvgRtt, nil

}

func IsOkIp(ip string) (*IpDetails, time.Duration, error) {
	rtt, err := PingRtt(ip)

	OKPing := func() bool {
		if rtt == 0 || err != nil {
			return false
		}
		return true
	}
	if !OKPing() {
		return nil, 0, fmt.Errorf("time out %s", ip)
	}
foo:
	ipDetails, err := GetIpDetails(ip)

	if err != nil {
		if err.Error() == "429" {
			fmt.Println("wait a minute")
			time.Sleep(time.Minute)
			goto foo
		}
		return nil, 0, err
	}
	return ipDetails, rtt, err
}
