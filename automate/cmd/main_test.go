package main_test

import (
	"os"
	"os/exec"
	"syscall"
	"testing"
	"time"
)

func TestMain(t *testing.T) {
	t.Log("starting main...")
	p := exec.Command("go", "run", "./main.go", "-verbose")
	p.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	p.Stdout = os.Stdout
	p.Stderr = os.Stderr
	wait := make(chan error)
	go func() {
		wait <- p.Run()
	}()
	go func() {
		<-time.After(2 * time.Second)
		syscall.Kill(-p.Process.Pid, syscall.SIGINT)
		wait <- nil
	}()
	if err := <-wait; err != nil {
		t.Error(err)
	}
	// wait for web port to free up when running tests back to back
	<-time.After(1 * time.Second)
}
