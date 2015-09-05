package main

import (
	"fmt"
	"net"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	cmds = append(cmds, &cobra.Command{
		Use:   "daemon [address]",
		Short: "Run the API daemon",
		Long:  "",
		Run:   modeDaemon,
	})
	inits = append(inits, initDaemon)
}

func initDaemon() {
	api_address := viper.GetString("api.address")
	if api_address != "" {
		BASE_URL = api_address
	}
}

func runHandlers() {
	for init := range inits {
		inits[init]()
	}
}

func modeDaemon(cmd *cobra.Command, args []string) {
	viper.SetDefault("file", "/boot/config.json")
	err := readKernelConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
	viper.SetConfigFile(viper.GetString("file"))
	readConfig()
	runHandlers()

	fmt.Println("Starting http server on " + BASE_URL)
	listener, err := net.Listen("tcp", BASE_URL)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = http.Serve(listener, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
}
