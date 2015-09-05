package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Test_runHandlers(t *testing.T) {
	// this is here because runHandlers() trips WriteDockerConfig
	MyWriteFile = func(filename string, contents []byte, mode os.FileMode) error {
		return nil
	}
	t.Log("Testing runHandlers")
	runHandlers()
}

func Test_initDaemon(t *testing.T) {
	// this is here because runHandlers() trips WriteDockerConfig
	viper.SetDefault("api.address", "127.0.0.1:8080")
	t.Log("Testing initDaemon")
	initDaemon()
}

func Test_modeDaemon(t *testing.T) {
	modeDaemon(&cobra.Command{}, []string{})
	if listener != nil {
		listener.Close()
	}
}

func Test_modeDaemonError(t *testing.T) {
	MyReadFile = func(path string) ([]byte, error) {
		return []byte{}, fmt.Errorf("This is an error from Test_modeDaemonError")
	}
	modeDaemon(&cobra.Command{}, []string{})
	if listener != nil {
		listener.Close()
	}
}
