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

func readKernelConfig() {
	cmdline, err := ioutil.ReadFile("/proc/cmdline")
	if err != nil {
		panic(fmt.Sprint("Error opening /proc/cmdline:", err.Error()))
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
}

func saveConfig() {
	cmd := exec.Command("mount", "LABEL=BOOT", "/boot")
	err := cmd.Run()
	if err != nil {
		return
	}
	b, _ := json.MarshalIndent(viper.AllSettings(), "", "  ")
	err = ioutil.WriteFile(viper.GetString("file"), b, 0644)
	if err != nil {
		panic(fmt.Sprint("Error opening file:", err.Error()))
	}
	fmt.Println("Config saved.")
	cmd = exec.Command("umount", "/boot")
	cmd.Run()

}

func runHandlers() {
	for init := range inits {
		inits[init]()
	}
}

func main() {
	MyReadFile = RealReadFile
	MyWriteFile = RealWriteFile
	viper.SetDefault("file", "/boot/config.json")
	readKernelConfig()
	viper.SetConfigFile(viper.GetString("file"))
	cmd := exec.Command("mount", "LABEL=BOOT", "/boot")
	err := cmd.Run()
	if err == nil {
		err = viper.ReadInConfig()
		cmd = exec.Command("umount", "/boot")
		cmd.Run()
		if err != nil {
			fmt.Println("Error reading config file: ", err.Error())
		}
		saveConfig()
	} else {
		fmt.Println(err.Error())
	}
	runHandlers()

	fmt.Println("Starting http server on :8080")
	http.ListenAndServe("127.0.0.1:8080", nil)

}
