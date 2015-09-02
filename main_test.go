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
