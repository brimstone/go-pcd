package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"testing"
)

func Test_readKernelConfig(t *testing.T) {
	t.Log("Testing kernelcmdline read")
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Errorf("Error starting listener socket")
	}
	_, port, err := net.SplitHostPort(ln.Addr().String())
	if err != nil {
		t.Errorf("Error getting listener socket port")
	}

	server := http.NewServeMux()
	server.HandleFunc("/config.json", func(w http.ResponseWriter, r *http.Request) {
		t.Logf("Giving up the goods")
		fmt.Fprintf(w, `{"files":[{"path":"blah", "content":"foo"}]}`)
	})
	server.HandleFunc("/config.yaml", func(w http.ResponseWriter, r *http.Request) {
		t.Logf("Giving up the goods")
		fmt.Fprintf(w, `files:
  - content: foo
    path: blah
`)
	})
	go http.Serve(ln, server)

	MyReadFile = func(filename string) ([]byte, error) {
		return []byte("BOOTIMAGE=primary initrd=primary ro hostname=pickles pcd.foo=bar pcd.url=http://localhost:" + port + "/config.json"), nil
	}

	err = readKernelConfig()
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	t.Log("Kernel config read without errors")
}

func Test_readKernelConfigError(t *testing.T) {
	t.Log("Testing kernelcmdline read error")
	MyReadFile = func(filename string) ([]byte, error) {
		return []byte(""), fmt.Errorf("Panic!")
	}

	err := readKernelConfig()
	if err == nil {
		t.Errorf("readKernelConfig() did not error")
	}
}

func Test_readConfig(t *testing.T) {
	t.Log("Testing readConfig read error")

	MyWriteFile = func(filename string, contents []byte, mode os.FileMode) error {
		return nil
	}

	MyExec = func(cmd string, arg ...string) ([]byte, error) {
		return []byte{}, nil
	}

	err := readConfig()
	if err != nil {
		t.Errorf("Error: %s", err)
	}
}

func Test_readConfigError1(t *testing.T) {
	t.Log("Testing readConfig exec error 1")

	MyExec = func(cmd string, arg ...string) ([]byte, error) {
		if cmd == "mount" {
			return nil, nil
		}
		return []byte{}, fmt.Errorf("Exec error")
	}

	err := readConfig()
	if err == nil {
		t.Errorf("Expected an error and didn't receive one.")
	}
}

func Test_readConfigError2(t *testing.T) {
	t.Log("Testing readConfig exec error 2")

	MyExec = func(cmd string, arg ...string) ([]byte, error) {
		if cmd == "mount" {
			return []byte{}, nil
		}
		return []byte{}, fmt.Errorf("Exec error")
	}

	err := readConfig()
	if err == nil {
		t.Errorf("Expected an error and didn't receive one.")
	}
}

func Test_saveConfigError1(t *testing.T) {
	t.Log("Testing saveConfig exec error")

	MyExec = func(cmd string, arg ...string) ([]byte, error) {
		return []byte{}, fmt.Errorf("Exec error")
	}

	err := saveConfig()
	if err == nil {
		t.Errorf("Error: %s", err)
	}
}

func Test_saveConfigError2(t *testing.T) {
	t.Log("Testing saveConfig read error")

	MyExec = func(cmd string, arg ...string) ([]byte, error) {
		return []byte{}, nil
	}

	MyWriteFile = func(filename string, contents []byte, mode os.FileMode) error {
		return fmt.Errorf("Write error")
	}

	err := saveConfig()
	if err == nil {
		t.Errorf("Error: %s", err)
	}
}

func Test_RealReadFileError(t *testing.T) {
	t.Log("Testing RealReadFile read error")
	_, err := RealReadFile("asdf")
	if err == nil {
		t.Errorf("Error: %s", err)
	}
}

func Test_RealWriteFileError(t *testing.T) {
	t.Log("Testing RealWriteFile read error")
	err := RealWriteFile("/asdf", []byte{}, 0644)
	if err == nil {
		t.Errorf("Error: %s", err)
	}
}

func Test_RealExecError(t *testing.T) {
	t.Log("Testing RealExec read error")
	_, err := RealExec("asdf")
	if err == nil {
		t.Errorf("Error: %s", err)
	}
}

func Test_main(t *testing.T) {
	main()
}
