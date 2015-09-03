package main

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

var restartTimer *time.Timer

func init() {
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
