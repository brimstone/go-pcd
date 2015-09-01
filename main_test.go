package main

import (
	"fmt"
	"testing"
)

func Test_readKernelConfig(t *testing.T) {
	t.Log("Testing kernelcmdline read")
	MyReadFile = func(filename string) ([]byte, error) {
		return []byte("BOOTIMAGE=primary initrd=primary ro hostname=pickles pcd.foo=bar"), nil
	}

	readKernelConfig()
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
