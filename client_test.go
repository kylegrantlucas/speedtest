package speedtest

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"time"

	sthttp "github.com/kylegrantlucas/speedtest/http"
)

func EmptyTest(t *testing.T) {
	t.Logf("Empty test...\n")
}

func TestNewClient(t *testing.T) {
	type args struct {
		config  *sthttp.SpeedtestConfig
		dlsizes []int
		ulsizes []int
		timeout time.Duration
	}
	tests := []struct {
		name    string
		args    args
		want    *Client
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewClient(tt.args.config, tt.args.dlsizes, tt.args.ulsizes, tt.args.timeout)
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

func TestClient_Download(t *testing.T) {
	f, err := os.Open("http/random750x750.jpg")
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
		server sthttp.Server
	}
	tests := []struct {
		name    string
		client  *Client
		args    args
		want    float64
		wantErr bool
	}{
		{
			name: "download test",
			client: &Client{
				HTTPClient: &sthttp.Client{
					SpeedtestConfig: &sthttp.SpeedtestConfig{
						ServersURL:      ts.URL,
						NumClosest:      1,
						NumLatencyTests: 1,
					},
					Timeout: (15 * time.Second),
				},
				DLSizes: []int{350, 500, 750, 1000, 1500, 2000, 2500, 3000, 3500, 4000},
			},
			args: args{
				server: sthttp.Server{
					URL: ts.URL + "/",
				},
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "download max test",
			client: &Client{
				HTTPClient: &sthttp.Client{
					SpeedtestConfig: &sthttp.SpeedtestConfig{
						ServersURL:      ts.URL,
						NumClosest:      1,
						NumLatencyTests: 1,
						AlgoType:        "max",
					},
					Timeout: (15 * time.Second),
				},
				DLSizes: []int{350, 500, 750, 1000, 1500, 2000, 2500, 3000, 3500, 4000},
			},
			args: args{
				server: sthttp.Server{
					URL: ts.URL + "/",
				},
			},
			want:    0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.client.Download(tt.args.server)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Download() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got < tt.want {
				t.Errorf("Client.Download() = %v, want greater than %v", got, tt.want)
			}
		})
	}
}

func TestClient_Upload(t *testing.T) {
	f, err := os.Open("http/random750x750.jpg")
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
		server sthttp.Server
	}
	tests := []struct {
		name    string
		client  *Client
		args    args
		want    float64
		wantErr bool
	}{
		{
			name: "upload test",
			client: &Client{
				HTTPClient: &sthttp.Client{
					SpeedtestConfig: &sthttp.SpeedtestConfig{
						ServersURL:      ts.URL + "/",
						NumClosest:      1,
						NumLatencyTests: 1,
					},
					Timeout: (15 * time.Second),
				},
				ULSizes: []int{int(0.25 * 1024 * 1024), int(0.5 * 1024 * 1024), int(1.0 * 1024 * 1024), int(1.5 * 1024 * 1024), int(2.0 * 1024 * 1024)},
			},
			args: args{
				server: sthttp.Server{
					URL: ts.URL + "/",
				},
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "upload max test",
			client: &Client{
				HTTPClient: &sthttp.Client{
					SpeedtestConfig: &sthttp.SpeedtestConfig{
						ServersURL:      ts.URL + "/",
						NumClosest:      1,
						NumLatencyTests: 1,
						AlgoType:        "max",
					},
					Timeout: (15 * time.Second),
				},
				ULSizes: []int{int(0.25 * 1024 * 1024), int(0.5 * 1024 * 1024), int(1.0 * 1024 * 1024), int(1.5 * 1024 * 1024), int(2.0 * 1024 * 1024)},
			},
			args: args{
				server: sthttp.Server{
					URL: ts.URL + "/",
				},
			},
			want:    0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.client.Upload(tt.args.server)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Upload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got < tt.want {
				t.Errorf("Client.Upload() = %v, want greater than %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetServer(t *testing.T) {
	x, err := ioutil.ReadFile("http/sthttp_test_servers.xml")
	if err != nil {
		t.Logf("Cannot read sthttp_test_servers.xml")
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, string(x))
	}))
	defer ts.Close()

	client := &sthttp.Client{
		SpeedtestConfig: &sthttp.SpeedtestConfig{ServersURL: ts.URL, NumClosest: 1, NumLatencyTests: 1},
		Config: &sthttp.Config{
			Lat: 32.5155,
			Lon: -90.1118,
		},
		Timeout: (15 * time.Second),
	}

	// servers, err := client.GetServers()
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		serverID string
	}
	tests := []struct {
		name    string
		client  *Client
		args    args
		want    sthttp.Server
		wantErr bool
	}{
		{
			name: "no server provided",
			client: &Client{
				HTTPClient: client,
			},
			args: args{
				serverID: "",
			},
			want: sthttp.Server{
				ID: "2630",
			},
			wantErr: false,
		},
		{
			name: "server provided",
			client: &Client{
				HTTPClient: client,
			},
			args: args{
				serverID: "2630",
			},
			want: sthttp.Server{
				ID: "2630",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.client.GetServer(tt.args.serverID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetServer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.ID != tt.want.ID {
				t.Errorf("Client.GetServer() = %v, want %v", got.ID, tt.want.ID)
			}
		})
	}
}

func TestClient_FindServer(t *testing.T) {
	type args struct {
		id          string
		serversList []sthttp.Server
	}
	tests := []struct {
		name   string
		client *Client
		args   args
		want   http.Server
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.client.FindServer(tt.args.id, tt.args.serversList); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.FindServer() = %v, want %v", got, tt.want)
			}
		})
	}
}
