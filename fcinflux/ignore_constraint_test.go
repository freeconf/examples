package fcinflux

import (
	"testing"

	"github.com/freeconf/yang/fc"
)

func TestIgnoreConstaint(t *testing.T) {
	ignores, err := compileIgnores([]string{
		"abc/123",
	})

	fc.AssertEqual(t, nil, err)
	fc.AssertEqual(t, 1, len(ignores))
	fc.AssertEqual(t, true, ignores[0].allow("abc"))
	fc.AssertEqual(t, false, ignores[0].allow("abc/123"))
}
