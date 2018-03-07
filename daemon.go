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
	inits["daemon"] = &initFunc{
		Func:   initDaemon,
		Status: false,
	}
}

func initDaemon() bool {
	if config.API != "" {
		BASE_URL = config.API
	}
	return true
}

func runHandlers() {
	flag := true
	for flag {
		flag = false
		for init := range inits {
			// If this module was already started, skip it this round
			if inits[init].Status {
				continue
			}
			log.Println("Starting", init, "handler")
			if inits[init].Func() {
				inits[init].Status = true
				flag = true
			}
		}
	}
	return
}

func modeDaemon(cmd *cobra.Command, args []string) {
	forever = make(chan bool)
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

	log.Println("Starting http server on " + BASE_URL)
	listener, err := net.Listen("tcp", BASE_URL)
	if err != nil {
		fmt.Println(err)
		return
	}
	go http.Serve(listener, nil)

	log.Println("Starting init handlers")
	runHandlers()

	<-forever
}
