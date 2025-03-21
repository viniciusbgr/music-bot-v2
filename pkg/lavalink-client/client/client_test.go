package client_test

import (
	"context"
	"errors"
	"io"
	"net"
	"os"
	"strconv"
	"testing"

	"github.com/viniciusbgr/music-bot-v2/pkg/lavalink-client/client"
	"github.com/viniciusbgr/music-bot-v2/pkg/logger"
)

func TestClient(t *testing.T) {
	t.Run("Test NewClient validations Errors", func(t *testing.T) {
		t.Run("Test Client with empty config", func(t *testing.T) {
			if _, err := client.NewClient(nil, nil, nil); err != client.ErrEmptyConfig {
				t.Errorf("Expected error %v, got %v", client.ErrEmptyConfig, err)
			}
		})

		t.Run("Test Client with invalid host", func(t *testing.T) {
			_, err := client.NewClient(&client.ClientConfig{Host: "invalid-host", ClientId: "TestClientId", Password: "TestClientWithInvalidHost"}, nil, nil)

			if err != client.ErrHostInvalid {
				t.Errorf("Expected error %v, got %v", client.ErrHostInvalid, err)
			}
		})

		t.Run("Test Client with empty client id", func(t *testing.T) {
			_, err := client.NewClient(&client.ClientConfig{"127.0.0.1", "TestClientWithEmptyClientId", "", false}, nil, nil)

			if err != client.ErrEmptyClientId {
				t.Errorf("Expected error %v, got %v", client.ErrEmptyClientId, err)
			}
		})

		t.Run("Test Client with nil logger", func(t *testing.T) {
			_, err := client.NewClient(&client.ClientConfig{"127.0.0.1", "TestClientWithNilLogger", "TestClientWithNilLogger", false}, nil, nil)

			if err != client.ErrNilLogger {
				t.Errorf("Expected error %v, got %v", client.ErrNilLogger, err)
			}
		})
	})
}

// This test, require a lavalink server running
// and the LAVALINK_URL and LAVALINK_PASSWORD environment variables set.
// for connections with TLS, set the TLS environment variable to true
func TestIntegrationWebSocketConnectionClient(t *testing.T) {
	addr, passwd, tls, clientId := os.Getenv("LAVALINK_URL"), os.Getenv("LAVALINK_PASSWORD"), os.Getenv("TLS"), "123456789"

	if addr == "" || passwd == "" {
		t.Skipf("Skipping integration test, LAVALINK_URL or LAVALINK_PASSWORD must be set")
	}
	var isTls bool

	if tls != "" {
		var err error

		isTls, err = strconv.ParseBool(tls)

		if err != nil {
			t.Fatalf("Error parsing TLS environment variable: %v", err)
		}
	} else {
		isTls = false
	}

	config := &client.ClientConfig{addr, passwd, clientId, isTls}

	log, _ := logger.New(io.Discard, 0, false)

	t.Run("Test Connection", func(t *testing.T) {
		client, err := client.NewClient(config, nil, log)

		if err != nil {
			t.Fatalf("Error creating client: %v", err)
		}

		if err := client.ConnectWS(context.Background()); err != nil {
			var opErr *net.OpError

			if errors.As(err, &opErr) && opErr.Op == "dial" && opErr.Err.Error() == "connect: connection refused" {
				t.Skip("Skipping integration test, connection refused...")

			} else {
				t.Fatalf("Error connecting client: %v", err)
			}
		} else {
			t.Log("Client connected successfully")

			if err := client.Disconnect(); err != nil {
				t.Fatalf("Error disconnecting client: %v", err)
			}
		}
	})
}
