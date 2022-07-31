package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/robinj730/rtc-client-go/client"
)

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	client, err := client.NewClient("localhost", 4000, true)
	if err != nil {
		panic(err)
	}
	defer client.Release()

	client.JoinRoom("2022")

	<-interrupt

}
