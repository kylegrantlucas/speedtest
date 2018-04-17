package http

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/kylegrantlucas/speedtest/coords"
	stxml "github.com/kylegrantlucas/speedtest/xml"
)

const max = "max"

// Config struct holds our config (users current ip, lat, lon and isp)
type Config struct {
	IP  string
	Lat float64
	Lon float64
	Isp string
}

// Client define a Speedtest HTTP client
type Client struct {
	Config          *Config
	SpeedtestConfig *SpeedtestConfig
	Timeout         time.Duration
	ReportChar      string
}

type SpeedtestConfig struct {
	ConfigURL       string
	ServersURL      string
	AlgoType        string
	NumClosest      int
	NumLatencyTests int
	Interface       string
	UserAgent       string
}

// NewClient define a new Speedtest client.
func NewClient(speedtestConfig *SpeedtestConfig, timeout time.Duration) (*Client, error) {
	client := &Client{
		Config:          nil,
		Timeout:         timeout,
		SpeedtestConfig: speedtestConfig,
	}

	config, err := client.GetConfig()
	if err != nil {
		return client, err
	}

	client.Config = &config
	return client, nil
}

// Server struct is a speedtest candidate server
type Server struct {
	URL      string
	Lat      float64
	Lon      float64
	Name     string
	Country  string
	CC       string
	Sponsor  string
	ID       string
	Distance float64
	Latency  float64
}

// ByDistance allows us to sort servers by distance
type ByDistance []Server

func (server ByDistance) Len() int {
	return len(server)
}

func (server ByDistance) Less(i, j int) bool {
	return server[i].Distance < server[j].Distance
}

func (server ByDistance) Swap(i, j int) {
	server[i], server[j] = server[j], server[i]
}

// ByLatency allows us to sort servers by latency
type ByLatency []Server

func (server ByLatency) Len() int {
	return len(server)
}

func (server ByLatency) Less(i, j int) bool {
	return server[i].Latency < server[j].Latency
}

func (server ByLatency) Swap(i, j int) {
	server[i], server[j] = server[j], server[i]
}

// checkHTTP tests if http response is successful (200) or not
func checkHTTP(resp *http.Response) bool {
	var ok bool
	if resp.StatusCode != 200 {
		ok = false
	} else {
		ok = true
	}
	return ok
}

// GetConfig downloads the master config from speedtest.net
func (stClient *Client) GetConfig() (c Config, err error) {
	c = Config{}

	client := &http.Client{
		Timeout: stClient.Timeout,
	}

	req, err := http.NewRequest("GET", stClient.SpeedtestConfig.ConfigURL, nil)
	if err != nil {
		return c, err
	}
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("User-Agent", stClient.SpeedtestConfig.UserAgent)

	resp, err := client.Do(req)
	if err != nil {
		return c, err
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			log.Printf("error closing body of config request: %v", err)
		}
	}()

	if !checkHTTP(resp) {
		return c, errors.New("couldn't connect to speedtest to get config")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return c, err
	}

	cx := new(stxml.XMLConfigSettings)

	err = xml.Unmarshal(body, &cx)
	if err != nil {
		return c, err
	}

	c.IP = cx.Client.IP
	c.Lat, err = strconv.ParseFloat(cx.Client.Lat, 64)
	if err != nil {
		return c, err
	}

	c.Lon, err = strconv.ParseFloat(cx.Client.Lon, 64)
	if err != nil {
		return c, err
	}

	c.Isp = cx.Client.Isp

	return c, err
}

// GetServers will get the full server list
func (stClient *Client) GetServers() (servers []Server, err error) {
	client := &http.Client{
		Timeout: stClient.Timeout,
	}

	req, err := http.NewRequest("GET", stClient.SpeedtestConfig.ServersURL, nil)
	if err != nil {
		return []Server{}, err
	}

	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("User-Agent", stClient.SpeedtestConfig.UserAgent)

	resp, err := client.Do(req)
	if err != nil {
		return []Server{}, err
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			log.Printf("error closing body of servers request: %v", err)
		}
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []Server{}, err
	}

	s := new(stxml.ServerSettings)

	err = xml.Unmarshal(body, &s)
	if err != nil {
		return []Server{}, err
	}

	for xmlServer := range s.ServersContainer.XMLServers {
		var err error
		server := new(Server)
		server.URL = s.ServersContainer.XMLServers[xmlServer].URL
		server.Lat, err = strconv.ParseFloat(s.ServersContainer.XMLServers[xmlServer].Lat, 64)
		if err != nil {
			log.Printf("error parsing lat: %v", err)
		}
		server.Lon, err = strconv.ParseFloat(s.ServersContainer.XMLServers[xmlServer].Lon, 64)
		if err != nil {
			log.Printf("error parsing lon: %v", err)
		}

		server.Name = s.ServersContainer.XMLServers[xmlServer].Name
		server.Country = s.ServersContainer.XMLServers[xmlServer].Country
		server.CC = s.ServersContainer.XMLServers[xmlServer].CC
		server.Sponsor = s.ServersContainer.XMLServers[xmlServer].Sponsor
		server.ID = s.ServersContainer.XMLServers[xmlServer].ID
		servers = append(servers, *server)
	}
	return servers, nil
}

// GetClosestServers takes the full server list and sorts by distance
func (stClient *Client) GetClosestServers(servers []Server) []Server {
	myCoords := coords.Coordinate{
		Lat: stClient.Config.Lat,
		Lon: stClient.Config.Lon,
	}
	for server := range servers {
		theirlat := servers[server].Lat
		theirlon := servers[server].Lon
		theirCoords := coords.Coordinate{Lat: theirlat, Lon: theirlon}

		servers[server].Distance = coords.HsDist(coords.DegPos(myCoords.Lat, myCoords.Lon), coords.DegPos(theirCoords.Lat, theirCoords.Lon))
	}

	sort.Sort(ByDistance(servers))

	return servers
}

// GetLatencyURL will return the proper url for the latency
func (stClient *Client) GetLatencyURL(server Server) string {
	u := server.URL
	splits := strings.Split(u, "/")
	baseURL := strings.Join(splits[1:len(splits)-1], "/")
	latencyURL := "http:/" + baseURL + "/latency.txt"

	return latencyURL
}

// GetLatency will test the latency (ping) the given server NUMLATENCYTESTS times and return either the lowest or average depending on what algorithm is set
func (stClient *Client) GetLatency(url string) (result float64, err error) {
	var latency time.Duration
	var minLatency time.Duration
	var avgLatency time.Duration

	for i := 0; i < stClient.SpeedtestConfig.NumLatencyTests; i++ {
		var failed bool
		var finish time.Time

		start := time.Now()

		client, err := stClient.getHTTPClient()
		if err != nil {
			return result, err
		}
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return 0, err
		}

		req.Header.Set("Cache-Control", "no-cache")
		req.Header.Set("User-Agent", stClient.SpeedtestConfig.UserAgent)

		resp, err := client.Do(req)

		if err != nil {
			return result, err
		}

		defer func() {
			err = resp.Body.Close()
			if err != nil {
				log.Printf("error closing body of latency request: %v", err)
			}
		}()

		finish = time.Now()
		_, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return result, err
		}

		if failed {
			latency = 1 * time.Minute
		} else {
			latency = finish.Sub(start)
		}

		if stClient.SpeedtestConfig.AlgoType == max {
			if minLatency == 0 {
				minLatency = latency
			} else if latency < minLatency {
				minLatency = latency
			}
		} else {
			avgLatency = avgLatency + latency
		}

	}

	if stClient.SpeedtestConfig.AlgoType == max {
		result = float64(time.Duration(minLatency.Nanoseconds())*time.Nanosecond) / 1000000
	} else {
		result = float64(time.Duration(avgLatency.Nanoseconds())*time.Nanosecond) / 1000000 / float64(stClient.SpeedtestConfig.NumLatencyTests)
	}

	return result, nil

}

// GetFastestServer test all servers until we find numServers that
// respond, then find the fastest of them.  Some servers show up in the
// master list but timeout or are "corrupt" therefore we bump their
// latency to something really high (1 minute) and they will drop out of
// this test
func (stClient *Client) GetFastestServer(servers []Server) (Server, error) {
	var successfulServers []Server

	for server := range servers {
		latency, err := stClient.GetLatency(stClient.GetLatencyURL(servers[server]))

		if err != nil {
			return Server{}, err
		}

		if latency < float64(1*time.Minute) {
			successfulServers = append(successfulServers, servers[server])
			successfulServers[server].Latency = latency
		}

		if len(successfulServers) == stClient.SpeedtestConfig.NumClosest {
			break
		}
	}

	sort.Sort(ByLatency(successfulServers))

	if len(successfulServers) == 0 {
		return Server{}, errors.New("no servers available")
	}

	return successfulServers[0], nil
}

// DownloadSpeed measures the mbps of downloading a URL
func (stClient *Client) DownloadSpeed(url string) (speed float64, err error) {
	start := time.Now()

	client, err := stClient.getHTTPClient()
	if err != nil {
		return 0, err
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("User-Agent", stClient.SpeedtestConfig.UserAgent)

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	finish := time.Now()
	bodyLen := len(body)

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			log.Printf("error closing body of download request: %v", err)
		}
	}()

	bits := float64(bodyLen * 8)
	megabits := bits / float64(1000) / float64(1000)
	seconds := finish.Sub(start).Seconds()
	mbps := megabits / seconds

	return mbps, err
}

// UploadSpeed measures the mbps to http.Post to a URL
func (stClient *Client) UploadSpeed(url string, mimetype string, data []byte) (speed float64, err error) {
	buf := bytes.NewBuffer(data)
	start := time.Now()

	client, err := stClient.getHTTPClient()
	if err != nil {
		return 0, err
	}
	resp, err := client.Post(url, mimetype, buf)
	finish := time.Now()
	if err != nil {
		return 0, err
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			log.Printf("error closing body of upload request: %v", err)
		}
	}()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	bits := float64(len(data) * 8)
	megabits := bits / float64(1000) / float64(1000)
	seconds := finish.Sub(start).Seconds()

	mbps := megabits / seconds
	return mbps, nil
}

func (stClient *Client) getHTTPClient() (*http.Client, error) {
	dialer := net.Dialer{
		Timeout:   stClient.Timeout,
		KeepAlive: stClient.Timeout,
	}

	transport := &http.Transport{
		Proxy:               http.ProxyFromEnvironment,
		Dial:                dialer.Dial,
		TLSHandshakeTimeout: stClient.Timeout,
	}

	client := &http.Client{
		Timeout:   stClient.Timeout,
		Transport: transport,
	}

	return client, nil
}
