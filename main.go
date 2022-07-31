package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	"github.com/robinj730/rtc-client-go/client"
)

type Config struct {
	RemoteAddr string `json:"server"`
	RemotePort int    `json:"port"`
	Secure     bool   `json:"secure"`
}

func main() {
	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}

	c := Config{}

	err = json.Unmarshal(file, &c)
	if err != nil {
		panic(err)
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	client, err := client.NewClient(c.RemoteAddr, c.RemotePort, c.Secure)
	if err != nil {
		panic(err)
	}
	defer client.Release()

	client.JoinRoom("2022")

	<-interrupt

}
