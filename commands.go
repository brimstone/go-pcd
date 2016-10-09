package main

import (
	"log"
	"strings"
)

func init() {
	inits["command"] = &initFunc{
		Func:   runCommands,
		Status: false,
	}
}

func runCommands() bool {
	if inits["docker"].Status {
		return false
	}
	for _, cmd := range config.Commands {
		cmds := strings.Split(cmd, " ")
		log.Println("Running cmd:", cmd)
		MyExec(cmds[0], cmds[1:]...)
	}
	return true
}
