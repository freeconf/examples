package automate

import (
	"testing"

	"github.com/freeconf/gconf/c2"
	"github.com/freeconf/gconf/device"
	"github.com/freeconf/gconf/meta"
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
