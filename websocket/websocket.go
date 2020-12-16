package websocket

import (
	"errors"

	"github.com/gorilla/websocket"
)

const (
	publicWsURL = "wss://ws.kraken.com/"
)

var errBinaryMessage = errors.New("unhandled binary message")

type Client struct {
	publicWs    *websocket.Conn
	receiveChan chan interface{}
}

func NewClient() (*Client, error) {
	client := &Client{
		receiveChan: make(chan interface{}),
	}
	var err error

	if client.publicWs, _, err = websocket.DefaultDialer.Dial(publicWsURL, nil); err != nil {
		return nil, err
	}

	go client.publicWsListener()

	return client, nil
}

func (client *Client) publicWsListener() {
	for {
		messageType, message, err := client.publicWs.ReadMessage()
		if err != nil {
			client.receiveChan <- err
		}
		switch messageType {
		case websocket.TextMessage:
			model, err := unmarshalReceivedMessage(message)
			if err != nil {
				client.receiveChan <- err
				break
			}
			client.receiveChan <- model
		default:
			client.receiveChan <- errBinaryMessage
		}
	}
}

func (client *Client) Listen() <-chan interface{} {
	return client.receiveChan
}

func (client *Client) send(message interface{}) error {
	return client.publicWs.WriteJSON(message)
}

func (client *Client) SendPing(ping Ping) error {
	ping.Event = "ping"
	return client.send(ping)
}

func (client *Client) SendSubscribe(subscribe Subscribe) error {
	subscribe.Event = "subscribe"
	return client.send(subscribe)
}
