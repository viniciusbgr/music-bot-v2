package client

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"time"

	"github.com/gorilla/websocket"
	"github.com/viniciusbgr/music-bot-v2/pkg/logger"
)

var (
	DefaultTimeOut    time.Duration = time.Second * 30
	DefaultClientName string        = "music-bot-v2"

	ErrConnectionAlreadyEstablished error = errors.New("client: connection already established")
	ErrConnectionNotEstablished     error = errors.New("client: connection not established")
	ErrEmptyClientId                error = errors.New("client: client id is don't set")
	ErrRequestNil                   error = errors.New("client: request is nil")
	ErrHostInvalid                  error = errors.New("client: host invalid")
	ErrNilLogger                    error = errors.New("client: logger is nil")
	ErrEmptyConfig                  error = errors.New("client: config nil or empty")

	regexDomain *regexp.Regexp = regexp.MustCompile(`^(?:[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}(:\d{1,5})?$`)
	regexIp     *regexp.Regexp = regexp.MustCompile(`^(?:\d{1,3}\.){3}\d{1,3}(:\d{1,5})?$`)

	// Endpoints for Client
	endpointWS string = "/v4/websocket"
)

type RawGenericSocketMessage struct {
	Operation string `json:"op"`
}

type ClientConfig struct {
	Host     string
	Password string
	ClientId string
	Tls      bool
}

type Client struct {
	httpClient       *http.Client
	socketConnection *websocket.Conn
	dialer           *websocket.Dialer
	config           *ClientConfig
	logger           *logger.Logger
	sessionId        string

	events MessageHandlers
}

func NewClient(config *ClientConfig, events MessageHandlers, log *logger.Logger) (*Client, error) {
	if config == nil || *config == (ClientConfig{}) {
		return nil, ErrEmptyConfig
	}

	if !regexDomain.MatchString(config.Host) && !regexIp.MatchString(config.Host) {
		return nil, ErrHostInvalid
	}

	if config.ClientId == "" {
		return nil, ErrEmptyClientId
	}

	if log == nil {
		return nil, ErrNilLogger
	}

	client := &http.Client{
		Timeout: DefaultTimeOut,
	}

	return &Client{
		httpClient: client,
		config:     config,
		events:     events,
		logger:     log,
	}, nil
}

func (c *Client) ConnectWS(ctx context.Context) error {
	if c.dialer == nil {
		c.dialer = &websocket.Dialer{EnableCompression: true, HandshakeTimeout: DefaultTimeOut}
	}

	header := http.Header{}
	header.Add("Authorization", c.config.Password)
	header.Add("User-Id", c.config.ClientId)
	header.Add("Client-Name", DefaultClientName)

	var (
		conn *websocket.Conn
		res  *http.Response
		err  error
	)

	if c.config.Tls {
		conn, res, err = c.dialer.DialContext(ctx, "wss://"+c.config.Host+endpointWS, header)
	} else {
		conn, res, err = c.dialer.DialContext(ctx, "ws://"+c.config.Host+endpointWS, header)
	}

	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		switch res.StatusCode {
		case 403:
			return errors.New("client: invalid password")
		case 404:
			return errors.New("client: invalid host")
		}
	}

	go c.listenWSMessages()

	c.socketConnection = conn

	return nil
}

func (c *Client) Disconnect() error {
	if c.socketConnection == nil {
		return ErrConnectionNotEstablished
	}

	return c.socketConnection.Close()
}

func (c *Client) listenWSMessages() {
	if len(c.events) == 0 {
		c.logger.Warn("client: no events provided, no listeners will be started")

		return
	}

	for {
		dataCode, data, err := c.socketConnection.ReadMessage()

		if err != nil {
			if websocket.IsCloseError(err, dataCode) {
				c.logger.Warn("client: connection closed, stopping listeners...")

				c.socketConnection = nil

				return
			} else {
				c.logger.Errorf("client: error on read message: %v", err)
			}
		}

		if dataCode == websocket.TextMessage {
			raw := RawGenericSocketMessage{}

			if err := json.Unmarshal(data, &raw); err != nil {
				c.logger.Errorf("client: error on unmarshal message: %v", err)
			}

			if handler, exists := c.events[Message(raw.Operation)]; exists {
				if err := handler(data); err != nil {
					c.logger.Errorf("client: error on handle event %s: %v", raw.Operation, err)

					continue
				}

				c.logger.Debug("client: event \"%s\" parsed", raw.Operation)
			} else {
				c.logger.Debug("client: not found handler for event -> %s", raw.Operation)
			}
		}
	}
}
