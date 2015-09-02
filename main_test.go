package main

import (
	"fmt"
	"os"
	"testing"
)

func Test_readKernelConfig(t *testing.T) {
	t.Log("Testing kernelcmdline read")
	MyReadFile = func(filename string) ([]byte, error) {
		return []byte("BOOTIMAGE=primary initrd=primary ro hostname=pickles pcd.foo=bar"), nil
	}

	err := readKernelConfig()
	if err != nil {
		t.Errorf("Error:", err.Error())
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
		t.Errorf("Error:", err.Error())
	}
}

func Test_readConfigError1(t *testing.T) {
	t.Log("Testing readConfig exec error 1")

	MyExec = func(cmd string, arg ...string) ([]byte, error) {
		return []byte{}, fmt.Errorf("Exec error")
	}

	err := readConfig()
	if err != nil {
		t.Errorf("Error:", err.Error())
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
	if err != nil {
		t.Errorf("Error:", err.Error())
	}
}
