package main

import (
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/brimstone/go-saverequest"
)

func Test_handleDockerBipPost(t *testing.T) {
	t.Log("Testing docker bip post")
	req, _ := saverequest.FakeRequest("POST", "/docker/bip", map[string]string{}, "172.16.0.1/24")
	w := httptest.NewRecorder()
	MyWriteFile = func(filename string, contents []byte, mode os.FileMode) error {
		return nil
	}
	MyExec = func(cmd string, arg ...string) ([]byte, error) {
		return []byte{}, nil
	}
	handleDockerBip(w, req)
	if w.Body.String() != "" {
		t.Errorf("Got unexpected response")
		t.Errorf("%d: %s", w.Code, w.Body.String())
		return
	}
	t.Log("Got proper response")
}

func Test_handleDockerBipGet(t *testing.T) {
	t.Log("Testing docker bip get")
	req, _ := saverequest.FakeRequest("GET", "/docker/bip", map[string]string{}, "")
	w := httptest.NewRecorder()
	handleDockerBip(w, req)
	if w.Body.String() != "172.16.0.1/24" {
		t.Errorf("Got unexpected response")
		t.Errorf("%d: %s", w.Code, w.Body.String())
		return
	}
	t.Log("Got proper bip")
}

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
