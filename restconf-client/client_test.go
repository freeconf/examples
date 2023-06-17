package demo

import (
	"strings"
	"testing"

	"github.com/freeconf/restconf"
	"github.com/freeconf/restconf/device"
	"github.com/freeconf/restconf/testdata"
	"github.com/freeconf/yang/fc"
	"github.com/freeconf/yang/source"
)

func TestClient(t *testing.T) {

	// setup -  start a server
	pathToYangFiles := "../yang:../car"
	serverYPath := source.Path(pathToYangFiles)
	carServer := testdata.New()
	local := device.New(serverYPath)
	local.Add("car", testdata.Manage(carServer))
	s := restconf.NewServer(local)
	defer s.Close()
	cfg := `{
		"fc-restconf": {
			"debug": true,
			"web" : {
				"port": ":9998"
			}
		},
		"car" : {
			"speed": 5
		}
	}`
	fc.RequireEqual(t, nil, local.ApplyStartupConfig(strings.NewReader(cfg)))

	connectClient()
}
