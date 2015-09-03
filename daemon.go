package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	cmds = append(cmds, &cobra.Command{
		Use:   "daemon [address]",
		Short: "Run the API daemon",
		Long:  "",
		Run:   mode_daemon,
	})
}

func mode_daemon(cmd *cobra.Command, args []string) {
	viper.SetDefault("file", "/boot/config.json")
	err := readKernelConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	viper.SetConfigFile(viper.GetString("file"))
	readConfig()
	runHandlers()

	fmt.Println("Starting http server on :8080")
	http.ListenAndServe("127.0.0.1:8080", nil)
}
