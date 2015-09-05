package main

import (
	"net/http/httptest"
	"os"
	"testing"

	"github.com/brimstone/go-saverequest"
	"github.com/spf13/cobra"
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

func Test_cmdDockerBipBlank(t *testing.T) {
	t.Log("Testing cmdDockerBip")
	MyAPIGet = func(path string) string {
		return ""
	}
	cmdDockerBip(&cobra.Command{}, []string{})
}

func Test_cmdDockerBipSetting(t *testing.T) {
	t.Log("Testing cmdDockerBip with argument")
	MyAPIPost = func(path string, payload string) {
		if payload != "pickles" {
			t.Errorf("Didn't get expected payload 'pickles': %s", payload)
		}
		return
	}
	cmdDockerBip(&cobra.Command{}, []string{"pickles"})
}

func Test_cmdDockerBipInvalid(t *testing.T) {
	t.Log("Testing cmdDockerBip with invalid arguments")
	cmdDockerBip(&cobra.Command{}, []string{"invalid", "invalid"})
}
