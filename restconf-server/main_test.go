package main_test

import (
	"os"
	"os/exec"
	"testing"

	"github.com/freeconf/yang/fc"
)

func TestMain(t *testing.T) {
	p := exec.Command("go", "run", "main.go", "-test=true")
	p.Stderr = os.Stderr
	p.Stdout = os.Stdout
	fc.AssertEqual(t, nil, p.Run())
}
