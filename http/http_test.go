package http

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strconv"
	"testing"
	"time"

	stxml "github.com/kylegrantlucas/speedtest/xml"
)

func TestCheckHTTPSuccess(t *testing.T) {
	resp := http.Response{}
	resp.StatusCode = 200
	r := checkHTTP(&resp)
	if r != true {
		t.Fail()
	}
}

func TestCheckHTTPFail(t *testing.T) {
	resp := http.Response{}
	resp.StatusCode = 404
	r := checkHTTP(&resp)
	if r != false {
		t.Fail()
	}
}

func TestGetLatencyURL(t *testing.T) {
	s := Server{}
	stc := Client{}
	s.URL = "http://example.com/speedtest/"
	u := stc.GetLatencyURL(s)
	if u != "http://example.com/speedtest/latency.txt" {
		t.Logf("Got latency URL: %s\n", u)
		t.Fail()
	}
}

func TestNewClient(t *testing.T) {
	x, err := ioutil.ReadFile("sthttp_test_config.xml")
	if err != nil {
		t.Logf("Cannot read sthttp_test_config.xml")
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, string(x))
	}))
	defer ts.Close()

	type args struct {
		speedtestConfig *SpeedtestConfig
		timeout         time.Duration
	}
	tests := []struct {
		name    string
		args    args
		want    *Client
		wantErr bool
	}{
		{
			name: "new client test",
			args: args{
				speedtestConfig: &SpeedtestConfig{
					ConfigURL:       ts.URL,
					NumClosest:      1,
					NumLatencyTests: 1,
				},
			},
			want: &Client{
				Timeout: 0 * time.Second,
				SpeedtestConfig: &SpeedtestConfig{
					ConfigURL:       ts.URL,
					NumClosest:      1,
					NumLatencyTests: 1,
				},
				Config: &Config{
					IP:  "23.124.0.25",
					Lat: 32.5155,
					Lon: -90.1118,
					Isp: "AT&T U-verse",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewClient(tt.args.speedtestConfig, tt.args.timeout)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestByDistance_Len(t *testing.T) {
	tests := []struct {
		name   string
		server ByDistance
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.server.Len(); got != tt.want {
				t.Errorf("ByDistance.Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestByDistance_Less(t *testing.T) {
	type args struct {
		i int
		j int
	}
	tests := []struct {
		name   string
		server ByDistance
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.server.Less(tt.args.i, tt.args.j); got != tt.want {
				t.Errorf("ByDistance.Less() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestByDistance_Swap(t *testing.T) {
	type args struct {
		i int
		j int
	}
	tests := []struct {
		name   string
		server ByDistance
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.server.Swap(tt.args.i, tt.args.j)
		})
	}
}

func TestByLatency_Len(t *testing.T) {
	tests := []struct {
		name   string
		server ByLatency
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.server.Len(); got != tt.want {
				t.Errorf("ByLatency.Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestByLatency_Less(t *testing.T) {
	type args struct {
		i int
		j int
	}
	tests := []struct {
		name   string
		server ByLatency
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.server.Less(tt.args.i, tt.args.j); got != tt.want {
				t.Errorf("ByLatency.Less() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestByLatency_Swap(t *testing.T) {
	type args struct {
		i int
		j int
	}
	tests := []struct {
		name   string
		server ByLatency
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.server.Swap(tt.args.i, tt.args.j)
		})
	}
}

func Test_checkHTTP(t *testing.T) {
	type args struct {
		resp *http.Response
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkHTTP(tt.args.resp); got != tt.want {
				t.Errorf("checkHTTP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetConfig(t *testing.T) {
	x, err := ioutil.ReadFile("sthttp_test_config.xml")
	if err != nil {
		t.Logf("Cannot read sthttp_test_config.xml")
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, string(x))
	}))
	defer ts.Close()

	tests := []struct {
		name     string
		stClient *Client
		wantC    Config
		wantErr  bool
	}{
		{
			name: "basic config test",
			stClient: &Client{
				SpeedtestConfig: &SpeedtestConfig{ConfigURL: ts.URL},
				Timeout:         (15 * time.Second),
			},
			wantC: Config{
				IP:  "23.124.0.25",
				Lat: 32.5155,
				Lon: -90.1118,
				Isp: "AT&T U-verse",
			},
			wantErr: false,
		},
		{
			name: "basic config failure",
			stClient: &Client{
				SpeedtestConfig: &SpeedtestConfig{ConfigURL: "bad url"},
				Timeout:         (15 * time.Second),
			},
			wantC:   Config{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotC, err := tt.stClient.GetConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotC, tt.wantC) {
				t.Errorf("Client.GetConfig() = %v, want %v", gotC, tt.wantC)
			}
		})
	}
}

func TestClient_GetServers(t *testing.T) {
	x, err := ioutil.ReadFile("sthttp_test_servers.xml")
	if err != nil {
		t.Logf("Cannot read sthttp_test_servers.xml")
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, string(x))
	}))
	defer ts.Close()

	tests := []struct {
		name        string
		stClient    *Client
		wantServers []Server
		wantErr     bool
	}{
		{
			name: "basic latency test",
			stClient: &Client{
				SpeedtestConfig: &SpeedtestConfig{ServersURL: ts.URL},
				Timeout:         (15 * time.Second),
			},
			wantServers: []Server{
				{
					URL: "http://88.84.191.230/speedtest/upload.php",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotServers, err := tt.stClient.GetServers()
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetServers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotServers[0].URL != tt.wantServers[0].URL {
				t.Errorf("Client.GetServers() = %v, want %v", gotServers[0].URL, tt.wantServers[0].URL)
			}
		})
	}
}

func TestClient_GetClosestServers(t *testing.T) {
	x, err := ioutil.ReadFile("sthttp_test_servers.xml")
	if err != nil {
		t.Logf("Cannot read sthttp_test_servers.xml")
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, string(x))
	}))
	defer ts.Close()

	client := &Client{
		SpeedtestConfig: &SpeedtestConfig{ServersURL: ts.URL},
		Config: &Config{
			Lat: 32.5155,
			Lon: -90.1118,
		},
		Timeout: (15 * time.Second),
	}

	servers, err := client.GetServers()
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		servers []Server
	}
	tests := []struct {
		name     string
		stClient *Client
		args     args
		want     []Server
	}{
		{
			name:     "basic closest servers test",
			stClient: client,
			args: args{
				servers: servers,
			},
			want: []Server{
				{
					ID: "2630",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.stClient.GetClosestServers(tt.args.servers); !reflect.DeepEqual(got[0].ID, tt.want[0].ID) {
				t.Errorf("Client.GetClosestServers() = %v, want %v", got[0].ID, tt.want[0].ID)
			}
		})
	}
}

func TestClient_GetLatencyURL(t *testing.T) {
	type args struct {
		server Server
	}
	tests := []struct {
		name     string
		stClient *Client
		args     args
		want     string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.stClient.GetLatencyURL(tt.args.server); got != tt.want {
				t.Errorf("Client.GetLatencyURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetLatency(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		fmt.Fprintln(w, "Hello World")
	}))
	defer ts.Close()

	type args struct {
		url string
	}
	tests := []struct {
		name       string
		stClient   *Client
		args       args
		wantResult float64
		wantErr    bool
	}{
		{
			name: "basic latency test",
			stClient: &Client{
				SpeedtestConfig: &SpeedtestConfig{NumLatencyTests: 5},
				Timeout:         (15 * time.Second),
			},
			args: args{
				url: ts.URL,
			},
			wantResult: 100,
			wantErr:    false,
		},
		{
			name: "basic max latency test",
			stClient: &Client{
				SpeedtestConfig: &SpeedtestConfig{NumLatencyTests: 5, AlgoType: "max"},
				Timeout:         (15 * time.Second),
			},
			args: args{
				url: ts.URL,
			},
			wantResult: 100,
			wantErr:    false,
		},
		{
			name: "basic latency failure",
			stClient: &Client{
				SpeedtestConfig: &SpeedtestConfig{NumLatencyTests: 1},
				Timeout:         (15 * time.Second),
			},
			args: args{
				url: "planned bad request",
			},
			wantResult: 0,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := tt.stClient.GetLatency(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetLatency() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResult < tt.wantResult {
				t.Errorf("Client.GetLatency() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func convertServers(body []byte) ([]Server, error) {
	var servers []Server

	s := new(stxml.ServerSettings)

	err := xml.Unmarshal(body, &s)
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

func TestClient_GetFastestServer(t *testing.T) {
	x, err := ioutil.ReadFile("sthttp_test_servers.xml")
	if err != nil {
		t.Logf("Cannot read sthttp_test_servers.xml")
	}

	servers, err := convertServers(x)
	if err != nil {
		t.Fatalf("error converting server fixture: %v", err)
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		fmt.Fprintln(w, "Hello World")
	}))
	defer ts.Close()

	client := &Client{
		SpeedtestConfig: &SpeedtestConfig{ServersURL: ts.URL, NumClosest: 2, NumLatencyTests: 2},
		Config: &Config{
			Lat: 32.5155,
			Lon: -90.1118,
		},
		Timeout: (15 * time.Second),
	}

	closest := client.GetClosestServers(servers)

	type args struct {
		servers []Server
	}
	tests := []struct {
		name     string
		stClient *Client
		args     args
		wantNil  bool
		wantErr  bool
	}{
		{
			name:     "basic fastsest server",
			stClient: client,
			args: args{
				servers: closest,
			},
			wantNil: false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.stClient.GetFastestServer(tt.args.servers)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetFastestServer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if reflect.DeepEqual(got, Server{}) && tt.wantNil != true {
				t.Errorf("Client.GetFastestServer() = %v", got)
			}
		})
	}
}

func TestClient_DownloadSpeed(t *testing.T) {
	f, err := os.Open("random750x750.jpg")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, b)
	}))
	defer ts.Close()

	type args struct {
		url string
	}
	tests := []struct {
		name      string
		stClient  *Client
		args      args
		wantSpeed float64
		wantErr   bool
	}{
		{
			name: "basic download test",
			stClient: &Client{
				SpeedtestConfig: &SpeedtestConfig{},
				Timeout:         (15 * time.Second),
			},
			args: args{
				url: ts.URL,
			},
			wantSpeed: 0,
			wantErr:   false,
		},
		{
			name: "basic download failure",
			stClient: &Client{
				SpeedtestConfig: &SpeedtestConfig{},
				Timeout:         (15 * time.Second),
			},
			args: args{
				url: "bad request",
			},
			wantSpeed: 0,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSpeed, err := tt.stClient.DownloadSpeed(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.DownloadSpeed() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotSpeed < tt.wantSpeed {
				t.Errorf("Client.DownloadSpeed() = %v, want %v", gotSpeed, tt.wantSpeed)
			}
		})
	}
}

func TestClient_UploadSpeed(t *testing.T) {
	f, err := os.Open("random750x750.jpg")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, b)
	}))
	defer ts.Close()

	type args struct {
		url      string
		mimetype string
		data     []byte
	}
	tests := []struct {
		name      string
		stClient  *Client
		args      args
		wantSpeed float64
		wantErr   bool
	}{
		{
			name: "basic upload test",
			stClient: &Client{
				SpeedtestConfig: &SpeedtestConfig{},
				Timeout:         (15 * time.Second),
			},
			args: args{
				url: ts.URL,
			},
			wantSpeed: 0,
			wantErr:   false,
		},
		{
			name: "basic upload failure",
			stClient: &Client{
				SpeedtestConfig: &SpeedtestConfig{},
				Timeout:         (15 * time.Second),
			},
			args: args{
				url: "bas request",
			},
			wantSpeed: 0,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSpeed, err := tt.stClient.UploadSpeed(tt.args.url, tt.args.mimetype, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.UploadSpeed() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotSpeed < tt.wantSpeed {
				t.Errorf("Client.UploadSpeed() = %v, want %v", gotSpeed, tt.wantSpeed)
			}
		})
	}
}

func TestClient_getHTTPClient(t *testing.T) {
	tests := []struct {
		name     string
		stClient *Client
		want     *http.Client
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.stClient.getHTTPClient()
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.getHTTPClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.getHTTPClient() = %v, want %v", got, tt.want)
			}
		})
	}
}
