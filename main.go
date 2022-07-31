package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	"github.com/robinj730/rtc-client-go/client"
)

func main() {
	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}

	config := client.Configuration{}

	err = json.Unmarshal(file, &config)
	if err != nil {
		panic(err)
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	client, err := client.NewClient(config)
	if err != nil {
		panic(err)
	}
	defer client.Release()

	client.JoinRoom("2022")

	<-interrupt

}
