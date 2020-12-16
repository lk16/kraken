package websocket

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

const (
	publicWsURL = "wss://ws.kraken.com/"
)

var errBinaryMessage = errors.New("unhandled binary message")

type Client struct {
	publicWs    *websocket.Conn
	receiveChan chan interface{}
	verbose     bool
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

func (client *Client) SetVerbose(verbose bool) {
	client.verbose = verbose
}

func (client *Client) publicWsListener() {
	for {
		messageType, message, err := client.publicWs.ReadMessage()
		if err != nil {
			client.receiveChan <- err
		}
		switch messageType {
		case websocket.TextMessage:
			if client.verbose {
				log.Printf("RECV: %s", string(message))
			}
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

func (client *Client) Send(rawMessage interface{}) error {

	send := func(rawMessage interface{}) error {
		bytes, err := json.Marshal(rawMessage)
		if err != nil {
			return err
		}

		if client.verbose {
			log.Printf("SEND: %s", string(bytes))
		}
		return client.publicWs.WriteJSON(rawMessage)
	}

	switch message := rawMessage.(type) {
	case Ping:
		message.Event = "ping"
		return send(message)
	case Subscribe:
		message.Event = "subscribe"
		return send(message)
	case Unsubscribe:
		message.Event = "unsubscribe"
		return send(message)
	default:
		return fmt.Errorf("unsupported message type %T", message)
	}

}
