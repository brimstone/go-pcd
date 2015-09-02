package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/viper"
)

var API_VERSION = "/1"

var inits []func()

var MyReadFile func(string) ([]byte, error)

func RealReadFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

var MyWriteFile func(string, []byte, os.FileMode) error

func RealWriteFile(filename string, contents []byte, mode os.FileMode) error {
	return ioutil.WriteFile(filename, contents, mode)
}

var MyExec func(string, ...string) ([]byte, error)

func RealExec(cmd string, arg ...string) ([]byte, error) {
	return exec.Command(cmd, arg...).CombinedOutput()
}

func readKernelConfig() error {
	cmdline, err := MyReadFile("/proc/cmdline")
	if err != nil {
		return err
	}
	options := strings.Split(strings.TrimSpace(string(cmdline)), " ")
	for _, option := range options {
		kv := strings.SplitN(option, "=", 2)
		if len(kv) < 2 {
			continue
		}
		if kv[0][0:4] == "pcd." {
			viper.Set(kv[0], kv[1])
		} else if kv[0] == "hostname" {
			viper.Set(kv[0], kv[1])
		}
	}
	return nil
}

func saveConfig() error {
	_, err := MyExec("mount", "LABEL=BOOT", "/boot")
	if err != nil {
		return err
	}
	b, _ := json.MarshalIndent(viper.AllSettings(), "", "  ")
	err = MyWriteFile(viper.GetString("file"), b, 0644)
	if err != nil {
		return err
	}
	fmt.Println("Config saved.")
	MyExec("umount", "/boot")
	return nil
}

func runHandlers() {
	for init := range inits {
		inits[init]()
	}
}

func readConfig() error {
	_, err := MyExec("mount", "LABEL=BOOT", "/boot")
	if err == nil {
		err = viper.ReadInConfig()
		_, err = MyExec("umount", "/boot")
		if err != nil {
			fmt.Println("Error reading config file: ", err.Error())
		}
		saveConfig()
	} else {
		fmt.Println(err.Error())
	}
	return nil
}

func main() {
	MyReadFile = RealReadFile
	MyWriteFile = RealWriteFile
	MyExec = RealExec
	viper.SetDefault("file", "/boot/config.json")
	err := readKernelConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	viper.SetConfigFile(viper.GetString("file"))
	readConfig()
	runHandlers()

	fmt.Println("Starting http server on :8080")
	http.ListenAndServe("127.0.0.1:8080", nil)

}
