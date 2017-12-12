package automate

import (
	"testing"

	"github.com/freeconf/c2g/c2"
	"github.com/freeconf/c2g/device"
	"github.com/freeconf/c2g/meta"
)

func Test_GoSystem(t *testing.T) {
	for _, role := range []string{"car", "garage"} {
		sys := &GoSystem{
			YangPath: &meta.FileStreamSource{Root: "../" + role},
			Map:      device.NewMap(),
		}
		if _, err := sys.New(role); err != nil {
			t.Error(err)
		}
		c2.AssertEqual(t, 1, sys.Map.Len())
	}
}
