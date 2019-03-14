package main

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

var restartTimer *time.Timer

func init() {
	inits["docker"] = &initFunc{
		Func: func() bool {
			if !inits["files"].Status {
				return false
			}
			WriteDockerConfig()
			RestartDocker()
			return true
		},
		Status: false,
	}
}

func RestartDocker() {
	if restartTimer != nil {
		restartTimer.Stop()
	}
	restartTimer = time.AfterFunc(time.Second, func() {
		if _, err := os.Stat("/service/docker"); err != nil {
			if os.IsNotExist(err) {
				// file does not exist
				log.Println("Enabling docker")
				err = os.Rename("/service.disable/docker", "/service/docker")
				if err != nil {
					log.Printf("Error enabling docker service: %s\n", err)
				}
			} else {
				// other error
				log.Println("Error checking docker service status")
			}
			return
		}
		log.Println("Restarting docker")
		MyExec("sv", "restart", "/service/docker")
	})
}

func WriteDockerConfig() {
	log.Println("Writing docker config")

	daemon, _ := json.Marshal(config.Docker)
	os.Mkdir("/etc/docker", 0755)
	MyWriteFile("/etc/docker/daemon.json", daemon, 0644)
}
