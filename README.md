# speedtest
[![CircleCI](https://circleci.com/gh/kylegrantlucas/speedtest.svg?style=svg)](https://circleci.com/gh/kylegrantlucas/speedtest) [![Maintainability](https://api.codeclimate.com/v1/badges/2130b46a52f698b3eaf1/maintainability)](https://codeclimate.com/github/kylegrantlucas/speedtest/maintainability) [![Test Coverage](https://api.codeclimate.com/v1/badges/2130b46a52f698b3eaf1/test_coverage)](https://codeclimate.com/github/kylegrantlucas/speedtest/test_coverage)

A golang package for running speedtests against speedtest.net.

## Usage
```
package main

import (
	"fmt"
	"time"

	"github.com/dchest/uniuri"
	"github.com/kylegrantlucas/speedtest"
	"github.com/kylegrantlucas/speedtest/http"
)

func main() {
	config := &http.SpeedtestConfig{
		ConfigURL:       "http://c.speedtest.net/speedtest-config.php?x=" + uniuri.New(),
		ServersURL:      "http://c.speedtest.net/speedtest-servers-static.php?x=" + uniuri.New(),
		AlgoType:        "max",
		NumClosest:      3,
		NumLatencyTests: 3,
		UserAgent:       "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.21 Safari/537.36",
	}

	client, err := speedtest.NewClient(
		config,
		speedtest.DefaultDLSizes,
		speedtest.DefaultULSizes,
		30*time.Second,
	)

	server, err := client.GetServer("")
	if err != nil {
		fmt.Printf("error getting server: %v", err)
	}

	dmbps, err := client.Download(server)
	if err != nil {
		fmt.Printf("error getting download: %v", err)
	}

	umbps, err := client.Upload(server)
	if err != nil {
		fmt.Printf("error getting upload: %v", err)
	}

	fmt.Printf("Ping (Lowest): %3.2f ms | Download (Max): %3.2f Mbps | Upload (Max): %3.2f Mbps\n", server.Latency, dmbps, umbps)
}
```
## Tests
`go test ./...`
## Thanks
## Contributing
