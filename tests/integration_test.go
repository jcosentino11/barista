//go:build integration

package integration

import (
	"fmt"
	"os"
	"testing"
	"time"

	"josephcosentino.me/barista/barista"
)

func TestBarista(t *testing.T) {
	server := barista.NewServer(barista.ServerConfig{
		Port: 8080,
	})

	err := server.Start()
	if err != nil {
		fmt.Printf("err: %s\n", err.Error())
		os.Exit(1)
		return
	}

	defer server.Stop()

	messageCallback := func(topic string, message string) {
		fmt.Printf("Received message on topic '%s': %s\n", topic, message)
	}

	client := barista.NewClient(barista.ClientConfig{
		ServerHost: "localhost",
		ServerPort: 8080,
	})

	if err := client.Connect(); err != nil {
		fmt.Printf("failed to connect: %s\n", err.Error())
		return
	}

	defer client.Close()

	err = client.Subscribe("topic", messageCallback)
	if err != nil {
		fmt.Printf("failed to subscribe: %s\n", err.Error())
		os.Exit(1)
		return
	}

	err = client.Publish("topic", "test")
	if err != nil {
		fmt.Printf("failed to publish: %s\n", err.Error())
		os.Exit(1)
		return
	}

	time.Sleep(100 * time.Millisecond)
}
