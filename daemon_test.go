package main

import (
	"os"
	"testing"
)

func Test_runHandlers(t *testing.T) {
	// this is here because runHandlers() trips WriteDockerConfig
	MyWriteFile = func(filename string, contents []byte, mode os.FileMode) error {
		return nil
	}
	t.Log("Testing runHandlers")
	runHandlers()
}
