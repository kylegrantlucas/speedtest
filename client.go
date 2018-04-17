package speedtest

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/kylegrantlucas/speedtest/http"
	"github.com/kylegrantlucas/speedtest/util"
)

var (
	// DefaultDLSizes defines the default download sizes
	DefaultDLSizes = []int{350, 500, 750, 1000, 1500, 2000, 2500, 3000, 3500, 4000}
	// DefaultULSizes defines the default upload sizes
	DefaultULSizes = []int{int(0.25 * 1024 * 1024), int(0.5 * 1024 * 1024), int(1.0 * 1024 * 1024), int(1.5 * 1024 * 1024), int(2.0 * 1024 * 1024)}
)

const max = "max"

// Client defines a Speedtester client tester
type Client struct {
	HTTPClient *http.Client
	DLSizes    []int
	ULSizes    []int
}

// Config define Speedtest settings
type Config struct {
	ConfigURL       string
	ServersURL      string
	AlgoType        string
	NumClosest      int
	NumLatencyTests int
	Interface       string
	Blacklist       []string
	UserAgent       string
}

func NewClient(config *http.SpeedtestConfig, dlsizes []int, ulsizes []int, timeout time.Duration) (*Client, error) {
	httpClient, err := http.NewClient(config, timeout)
	if err != nil {
		return &Client{}, err
	}

	return &Client{
		HTTPClient: httpClient,
		DLSizes:    dlsizes,
		ULSizes:    ulsizes,
	}, nil
}

// Download will perform the "normal" speedtest download test
func (client *Client) Download(server http.Server) (float64, error) {
	var urls []string
	var maxSpeed float64
	var avgSpeed float64

	// http://speedtest1.newbreakcommunications.net/speedtest/speedtest/
	for size := range client.DLSizes {
		url := server.URL

		splits := strings.Split(url, "/")
		baseURL := strings.Join(splits[1:len(splits)-1], "/")

		randomImage := fmt.Sprintf("random%dx%d.jpg", client.DLSizes[size], client.DLSizes[size])
		downloadURL := "http:/" + baseURL + "/" + randomImage
		urls = append(urls, downloadURL)
	}

	for u := range urls {
		dlSpeed, err := client.HTTPClient.DownloadSpeed(urls[u])
		if err != nil {
			return 0, err
		}

		if client.HTTPClient.SpeedtestConfig.AlgoType == max {
			if dlSpeed > maxSpeed {
				maxSpeed = dlSpeed
			}
		} else {
			avgSpeed = avgSpeed + dlSpeed
		}
	}

	if client.HTTPClient.SpeedtestConfig.AlgoType != max {
		return avgSpeed / float64(len(urls)), nil
	}
	return maxSpeed, nil

}

// Upload runs a "normal" speedtest upload test
func (client *Client) Upload(server http.Server) (float64, error) {
	// https://github.com/sivel/speedtest-cli/blob/master/speedtest-cli
	var ulsize []int
	var maxSpeed float64
	var avgSpeed float64

	for size := range client.ULSizes {
		ulsize = append(ulsize, client.ULSizes[size])
	}

	for i := 0; i < len(ulsize); i++ {
		r := util.Urandom(ulsize[i])
		ulSpeed, err := client.HTTPClient.UploadSpeed(server.URL, "text/xml", r)
		if err != nil {
			return 0, err
		}

		if client.HTTPClient.SpeedtestConfig.AlgoType == max {
			if ulSpeed > maxSpeed {
				maxSpeed = ulSpeed
			}
		} else {
			avgSpeed = avgSpeed + ulSpeed
		}

	}

	if client.HTTPClient.SpeedtestConfig.AlgoType != max {
		return avgSpeed / float64(len(ulsize)), nil
	}
	return maxSpeed, nil
}

func (client *Client) GetServer(serverID string) (http.Server, error) {
	server := http.Server{}

	allServers, err := client.HTTPClient.GetServers()
	if err != nil {
		return server, err
	}

	if serverID != "" {
		server = client.FindServer(serverID, allServers)
		server.Latency, err = client.HTTPClient.GetLatency(client.HTTPClient.GetLatencyURL(server))
		if err != nil {
			return server, err
		}
	} else {
		closestServers := client.HTTPClient.GetClosestServers(allServers)
		server, err = client.HTTPClient.GetFastestServer(closestServers)
		if err != nil {
			return server, err
		}
	}

	return server, nil
}

// FindServer will find a specific server in the servers list
func (client *Client) FindServer(id string, serversList []http.Server) http.Server {
	var foundServer http.Server
	for s := range serversList {
		if serversList[s].ID == id {
			foundServer = serversList[s]
		}
	}
	if foundServer.ID == "" {
		log.Printf("cannot locate server id '%s' in our list of speedtest servers", id)
	}
	return foundServer
}
