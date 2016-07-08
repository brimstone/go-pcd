package main

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

func init() {
	http.HandleFunc("/version", handleVersion)
	cmds = append(cmds, &cobra.Command{
		Use:   "version",
		Short: "Get the client and daemon version",
		Long:  "This gets the version of the client and daemon.",
		Run:   cmdVersion,
	})
}

func cmdVersion(cmd *cobra.Command, args []string) {
	fmt.Println("Client:")
	fmt.Println(" Version:", COMMITHASH)
	fmt.Println()
	daemonVersion := MyAPIGet("version")
	fmt.Println("Server:")
	fmt.Println(" Version:", daemonVersion)
}

func handleVersion(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(COMMITHASH))
}
