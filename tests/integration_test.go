//go:build integration

package integration

import (
	"testing"
	"time"

	"josephcosentino.me/barista/barista"
)

func TestBarista(t *testing.T) {
	server := barista.NewServer(barista.ServerConfig{
		Port: 8080,
	})
	defer server.Stop()

	err := server.Start()
	if err != nil {
		t.Fatalf("unable to start server: %s\n", err.Error())
	}

	messageCallback := func(topic string, message string) {
		t.Logf("received message on topic '%s': %s\n", topic, message)
	}

	client := barista.NewClient(barista.ClientConfig{
		ServerHost: "localhost",
		ServerPort: 8080,
	})
	defer client.Stop()

	if err := client.Start(); err != nil {
		t.Fatalf("failed to start: %s\n", err.Error())
	}

	err = client.Subscribe("topic", messageCallback)
	if err != nil {
		t.Fatalf("failed to subscribe: %s\n", err.Error())
	}

	err = client.Publish("topic", "test")
	if err != nil {
		t.Fatalf("failed to publish: %s\n", err.Error())
	}

	time.Sleep(100 * time.Millisecond)
}
