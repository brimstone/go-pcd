package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/spf13/cobra"
)

type ConfigType struct {
	Hostname string `json:"hostname"`
	API      string `json:"api"`
	Docker   struct {
		Bip string
	}
	Files []struct {
		Path    string `json:"path"`
		Content string `json:"content"`
	} `json:"files"`
	Commands []string `json:"commands"`
}

type initFunc struct {
	Status bool
	Func   func() bool
}

var (
	COMMITHASH    = "dev"
	BUILDDATETIME = "today"
	API_VERSION   = "1"
	BASE_URL      = "127.0.0.1:8080"
	MyAPIGet      func(string) string
	MyAPIPost     func(string, string)
	MyExec        func(string, ...string) ([]byte, error)
	MyReadFile    func(string) ([]byte, error)
	MyWriteFile   func(string, []byte, os.FileMode) error
	cmds          []*cobra.Command
	configfile    = "/boot/config.yaml"
	inits         = make(map[string]*initFunc)
	listener      net.Listener
	forever       chan bool
	config        ConfigType
)

func RealReadFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

func RealWriteFile(filename string, contents []byte, mode os.FileMode) error {
	return ioutil.WriteFile(filename, contents, mode)
}

func RealExec(cmd string, arg ...string) ([]byte, error) {
	return exec.Command(cmd, arg...).CombinedOutput()
}

func getUrl(url string) *http.Response {
	for {
		resp, err := http.Get(url)
		if err == nil {
			return resp
		}
		log.Println("Error getting config from url:", err)
		log.Println("Waiting to retry")
		time.Sleep(time.Second * 2)
	}
}

func readKernelConfig() error {
	cmdline, err := MyReadFile("/proc/cmdline")
	if err != nil {
		return err
	}
	kernel := make(map[string]string)
	options := strings.Split(strings.TrimSpace(string(cmdline)), " ")
	for _, option := range options {
		kv := strings.SplitN(option, "=", 2)
		if len(kv) < 2 {
			continue
		}
		if kv[0] == "pcd.url" {
			log.Println("Got a url", kv[1])
			resp := getUrl(kv[1])
			configcontents, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Println("Error getting config from url:", err)
			}
			err = yaml.Unmarshal(configcontents, &config)
			if err != nil {
				log.Println("Error parsing config:", err)
			}
			resp.Body.Close()

		} else if kv[0][0:4] == "pcd." {
			kernel[kv[0]] = kv[1]
		} else if kv[0] == "hostname" {
			kernel[kv[0]] = kv[1]
		}
	}
	for k, v := range kernel {
		if k == "hostname" {
			config.Hostname = v
		}
	}

	return nil
}

func saveConfig() error {
	_, err := MyExec("mount", "LABEL=BOOT", "/boot")
	if err != nil {
		return err
	}
	b, _ := yaml.Marshal(config)
	err = MyWriteFile(configfile, b, 0644)
	if err != nil {
		return err
	}
	fmt.Println("Config saved.")
	MyExec("umount", "/boot")
	return nil
}

func readConfig() error {
	_, err := MyExec("mount", "LABEL=BOOT", "/boot")
	if err != nil {
		return nil
	}
	// we have our defaults, kernel config, any url config. Now apply what's on disk
	// The file on disk shouldn't overwrite anything we already have.
	configcontents, err := ioutil.ReadFile(configfile)
	err = yaml.Unmarshal(configcontents, &config)
	_, err = MyExec("umount", "/boot")
	if err != nil {
		return err
	}
	saveConfig()
	return nil
}

func init() {
	MyAPIGet = APIGet
	MyAPIPost = APIPost
	MyReadFile = RealReadFile
	MyWriteFile = RealWriteFile
	MyExec = RealExec
}

func main() {
	var rootCmd = &cobra.Command{
		Use:  os.Args[0],
		Long: "Pancake Crop Deli Control Program",
	}

	for cmd := range cmds {
		rootCmd.AddCommand(cmds[cmd])
	}
	rootCmd.PersistentFlags().StringVarP(&BASE_URL, "address", "a", "127.0.0.1:8080", "Address for API server")
	rootCmd.Execute()
}
