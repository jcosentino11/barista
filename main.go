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

	server := barista.NewServer(barista.ServerConfig{
		Port: 8080,
	})
	err := server.Start()
	if err != nil {
		fmt.Printf("err: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Printf("server started on port: %d\n", server.Config.Port)

	<-interrupt

	fmt.Printf("stopping server\n")

	err = server.Stop()
	if err != nil {
		fmt.Printf("failed to stop server: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Printf("server stopped\n")

	os.Exit(0)
}
