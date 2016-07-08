package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/spf13/cobra"
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
	if config.API != "" {
		BASE_URL = config.API
	}
}

func runHandlers() {
	for init := range inits {
		log.Println("Starting", init, "handler")
		inits[init]()
	}
}

func modeDaemon(cmd *cobra.Command, args []string) {
	log.Println("Reading kernel config")
	err := readKernelConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Println("Reading local config")
	err = readConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Println("Starting init handlers")
	runHandlers()

	log.Println("Starting http server on " + BASE_URL)
	listener, err := net.Listen("tcp", BASE_URL)
	if err != nil {
		fmt.Println(err)
		return
	}
	go http.Serve(listener, nil)
	forever = make(chan bool)
	<-forever
}
