package main

import "strings"

func init() {
	inits = append(inits, runCommands)
}

func runCommands() {
	for _, cmd := range config.Commands {
		cmds := strings.Split(cmd, " ")
		MyExec(cmds[0], cmds[1:]...)
	}
}
