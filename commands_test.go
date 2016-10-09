package main

import "testing"

func Test_RunCommands(t *testing.T) {
	config.Command = "uptime"
	MyExec = func(cmd string, arg ...string) ([]byte, error) {
		t.Log("Preventing", cmd)
		return []byte{}, nil
	}
	runCommands()
}
