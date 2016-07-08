package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/cobra"
)

func init() {
	http.HandleFunc("/"+API_VERSION+"/docker/bip", handleDockerBip)
	cmds = append(cmds, &cobra.Command{
		Use:   "docker/bip [address]",
		Short: "Get or Set Docker Bridge IP",
		Long:  "This gets or sets the IP for the internal docker bridge.",
		Run:   cmdDockerBip,
	})
}

func cmdDockerBip(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println(MyAPIGet("docker/bip"))
	} else if len(args) == 1 {
		MyAPIPost("docker/bip", args[0])
	} else {
		cmd.Help()
	}
}

func handleDockerBip(w http.ResponseWriter, r *http.Request) {
	//w.Header().Add("Content-Type", "application/json")
	if r.Method == "POST" {
		body, _ := ioutil.ReadAll(r.Body)
		config.Docker.Bip = string(body)
		WriteDockerConfig()
		RestartDocker()
		saveConfig()
	} else {
		bip := config.Docker.Bip
		w.Write([]byte(bip))
	}
}
