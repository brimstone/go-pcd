package main

import (
	"testing"
	"time"
)

func Test_RestartDocker(t *testing.T) {
	MyExec = func(cmd string, arg ...string) ([]byte, error) {
		return []byte{}, nil
	}
	t.Log("Testing RestartDocker")
	RestartDocker()
	RestartDocker()
	t.Log("Sleeping 1 second to trigger timer")
	time.Sleep(time.Second)
	t.Log("End of timer test")
}
