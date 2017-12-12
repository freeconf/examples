package main_test

import (
	"os"
	"os/exec"
	"syscall"
	"testing"
)

func TestMain(t *testing.T) {
	t.Log("starting main...")
	p := exec.Command("go", "generate", ".")
	p.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	p.Stdout = os.Stdout
	p.Stderr = os.Stderr
	err := p.Run()
	if err != nil {
		t.Error(err)
	}
}
