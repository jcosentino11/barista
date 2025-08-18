package main

import (
	"fmt"
	"os"
	"os/signal"

	"josephcosentino.me/barista/barista"
)

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	server := barista.NewServer()
	err := server.Start()
	if err != nil {
		fmt.Printf("err: %s\n", err.Error())
		os.Exit(1)
	}

	<-interrupt

	err = server.Stop()
	if err != nil {
		fmt.Printf("failed to stop server: %s\n", err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}
