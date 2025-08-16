package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	server := NewServer()

	err := server.Start()
	if err != nil {
		fmt.Printf("err: %s\n", err.Error())
		os.Exit(1)
		return
	}

	defer func() {
		server.Stop()
	}()

	messageCallback := func(topic string, message string) {
		fmt.Printf("Received message on topic '%s': %s\n", topic, message)
	}

	client := NewClient()

	err = client.Subscribe("topic", messageCallback)
	if err != nil {
		fmt.Printf("err: %s\n", err.Error())
		os.Exit(1)
		return
	}

	err = client.Publish("topic", "test")
	if err != nil {
		fmt.Printf("err: %s\n", err.Error())
		os.Exit(1)
		return
	}

	time.Sleep(100 * time.Millisecond)
}
