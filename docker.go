package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/spf13/viper"
)

var restartTimer *time.Timer

func init() {
	http.HandleFunc(API_VERSION+"/docker/bip", handleDockerBip)
	inits = append(inits, func() {
		WriteDockerConfig()
		RestartDocker()
	})
}

func RestartDocker() {
	if restartTimer != nil {
		restartTimer.Stop()
	}
	restartTimer = time.AfterFunc(time.Second, func() {
		MyExec("sv", "restart", "/service/docker")
	})
}

func WriteDockerConfig() {
	fmt.Println("Writing docker config")
	config := fmt.Sprintf("BIP=\"%s\"\n", viper.GetString("docker.bip"))
	os.Mkdir("/etc/config", 0755)
	MyWriteFile("/etc/config/docker", []byte(config), 0644)
}

func handleDockerBip(w http.ResponseWriter, r *http.Request) {
	//w.Header().Add("Content-Type", "application/json")
	if r.Method == "POST" {
		body, _ := ioutil.ReadAll(r.Body)
		viper.Set("docker.bip", string(body))
		WriteDockerConfig()
		RestartDocker()
		saveConfig()
	} else {
		bip := viper.GetString("docker.bip")
		w.Write([]byte(bip))
	}
}
