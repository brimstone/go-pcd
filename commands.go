package main

import "log"

func init() {
	inits["command"] = &initFunc{
		Func:   runCommands,
		Status: false,
	}
}

func runCommands() bool {
	if !inits["docker"].Status {
		return false
	}
	log.Println("Running cmd:", config.Command)
	result, err := MyExec("/bin/sh", "-c", config.Command)
	log.Println("Exit ", err, "Output:", string(result))
	return true
}
