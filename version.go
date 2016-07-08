package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"

	"github.com/spf13/cobra"
)

type versionStruct struct {
	APIVersion    string
	Version       string
	GitCommit     string
	Builddatetime string
	Arch          string
}

var version versionStruct

func init() {
	http.HandleFunc("/"+API_VERSION+"/version", handleVersion)
	cmds = append(cmds, &cobra.Command{
		Use:   "version",
		Short: "Get the client and daemon version",
		Long:  "This gets the version of the client and daemon.",
		Run:   cmdVersion,
	})
	version = versionStruct{
		APIVersion:    API_VERSION,
		Version:       "TODO",
		GitCommit:     COMMITHASH,
		Builddatetime: BUILDDATETIME,
		Arch:          runtime.GOOS + "/" + runtime.GOARCH,
	}
}

func cmdVersion(cmd *cobra.Command, args []string) {
	var daemonVersion versionStruct

	fmt.Println("Client:")
	fmt.Println(" Arch:", version.Arch)
	fmt.Println(" APIVersion:", version.APIVersion)
	fmt.Println(" Build:", version.Builddatetime)
	fmt.Println(" GitCommit:", version.GitCommit)
	fmt.Println(" Version:", version.Version)
	fmt.Println()

	daemonJSON := MyAPIGet("/version")
	json.Unmarshal([]byte(daemonJSON), &daemonVersion)
	fmt.Println("Server:")
	fmt.Println(" Arch:", daemonVersion.Arch)
	fmt.Println(" APIVersion:", daemonVersion.APIVersion)
	fmt.Println(" Build:", daemonVersion.Builddatetime)
	fmt.Println(" GitCommit:", daemonVersion.GitCommit)
	fmt.Println(" Version:", daemonVersion.Version)
}

func handleVersion(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	data, _ := json.MarshalIndent(version, "", "  ")
	w.Write(data)
}
