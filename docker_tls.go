package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/cobra"
)

func init() {
	http.HandleFunc("/"+API_VERSION+"/docker/tls", handleDockerTLS)
	cmds = append(cmds, &cobra.Command{
		Use:   "docker/tls",
		Short: "Get or Set Docker TLS settings",
		Long:  "This gets or sets the TLS settings for the docker daemon.",
		Run:   cmdDockerTLS,
	})
}

func cmdDockerTLS(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println(MyAPIGet("docker/tls"))
	} else if len(args) == 1 {
		MyAPIPost("docker/tls", args[0])
	} else {
		cmd.HelpFunc()
	}
}

func handleDockerTLS(w http.ResponseWriter, r *http.Request) {
	//w.Header().Add("Content-Type", "application/json")
	if r.Method == "POST" {
		body, _ := ioutil.ReadAll(r.Body)
		tls := false
		if string(body) == "true" {
			tls = true
		}
		config.Docker.TLS = tls
		WriteDockerConfig()
		RestartDocker()
		saveConfig()
	} else {
		tls := config.Docker.TLS
		if tls {
			w.Write([]byte("true"))
		} else {
			w.Write([]byte("false"))
		}
	}
}
