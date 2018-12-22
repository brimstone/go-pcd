package main

import (
	"io/ioutil"
	"log"
)

func init() {
	inits["files"] = &initFunc{
		Func:   initFiles,
		Status: false,
	}
}

func initFiles() bool {
	for _, f := range config.Files {
		log.Printf("Writing %s\n", f.Path)
		err := ioutil.WriteFile(f.Path, []byte(f.Content), 0777)
		if err != nil {
			return false
		}
	}
	return true
}
