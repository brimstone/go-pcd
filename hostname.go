package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	http.HandleFunc("/"+API_VERSION+"/hostname", handleHostname)
	cmds = append(cmds, &cobra.Command{
		Use:   "hostname [hostname]",
		Short: "Get or Set Hostname",
		Long:  "This gets or sets the hostname for the system.",
		Run:   cmdHostname,
	})
	inits = append(inits, initHostname)
}

func cmdHostname(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Printf("%s", MyAPIGet("hostname"))
	} else if len(args) == 1 {
		MyAPIPost("hostname", args[0])
	} else {
		cmd.Help()
	}
}

func initHostname() {
	hostname := viper.GetString("hostname")
	if hostname != "" {
		SetHostname([]byte(hostname))
	}
}

func SetHostname(hostname []byte) error {
	fmt.Println("Setting hostname")
	err := MyWriteFile("/proc/sys/kernel/hostname", hostname, 0644)
	return err
}

func handleHostname(w http.ResponseWriter, r *http.Request) {
	//w.Header().Add("Content-Type", "application/json")
	if r.Method == "POST" {
		body, _ := ioutil.ReadAll(r.Body)
		err := SetHostname(body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		viper.Set("hostname", string(body))
		saveConfig()
	} else {
		hostnameBytes, err := MyReadFile("/proc/sys/kernel/hostname")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(hostnameBytes)
	}
}
