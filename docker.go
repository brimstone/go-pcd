package main

import (
	"fmt"
	"os"
	"time"
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
	opts := fmt.Sprintf("BIP=\"%s\"\n", config.Docker.Bip)
	os.Mkdir("/etc/config", 0755)
	MyWriteFile("/etc/config/docker", []byte(opts), 0644)
}
