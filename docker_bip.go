package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	http.HandleFunc("/"+API_VERSION+"/docker/bip", handleDockerBip)
	cmds = append(cmds, &cobra.Command{
		Use:   "docker/bip [address]",
		Short: "Get or Set Docker Bridge IP",
		Long:  "This gets or sets the IP for the internal docker bridge.",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println(APIGet("docker/bip"))
			} else if len(args) == 1 {
				APIPost("docker/bip", args[0])
			} else {
				cmd.Help()
			}
		},
	})
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
