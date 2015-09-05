package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func APIGet(path string) string {
	resp, err := http.Get("http://" + BASE_URL + "/" + API_VERSION + "/" + path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return string(body)
}

func APIPost(path string, payload string) {
	_, err := http.Post("http://"+BASE_URL+"/"+API_VERSION+"/"+path,
		"text/plain",
		bytes.NewBufferString(payload))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
