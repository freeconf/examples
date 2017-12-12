package automate

import (
	"fmt"

	"github.com/freeconf/examples/car"
	"github.com/freeconf/examples/garage"

	"github.com/freeconf/c2g/device"
	"github.com/freeconf/c2g/meta"
)

type GoSystem struct {
	Map      device.Map
	YangPath meta.StreamSource
	deviceId int
}

func (self *GoSystem) nextDeviceId(role string) string {
	id := fmt.Sprintf("%s%d", role, self.deviceId)
	self.deviceId++
	return id
}

func (self *GoSystem) New(role string) (*Handle, error) {
	d := device.New(self.YangPath)
	switch role {
	case "car":
		if err := d.Add("car", car.Manage(car.New())); err != nil {
			return nil, err
		}
	case "garage":
		g := garage.NewGarage()
		defer func() {
			garage.ManageCars(g, self.Map)
		}()
		if err := d.Add("garage", garage.Manage(g)); err != nil {
			return nil, err
		}
	default:
		panic(role)
	}
	deviceId := self.nextDeviceId(role)
	self.Map.Add(deviceId, d)
	return &Handle{
		Id:     deviceId,
		Device: d,
		Close:  func() {},
	}, nil
}
