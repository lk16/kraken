package websocket

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/lk16/kraken/rest"
)

const (
	publicWsURL       = "wss://ws.kraken.com/"
	privateWsURL      = "wss://ws-auth.kraken.com"
	keepAliveDuration = 10 * time.Second
	readDeadline      = keepAliveDuration
)

var errBinaryMessage = errors.New("unhandled binary message")

type Client struct {
	publicWs     *websocket.Conn
	receiveChan  chan interface{}
	verbose      bool
	privateToken string
	privateWs    *websocket.Conn
}

func NewClient() (*Client, error) {
	client := &Client{
		receiveChan: make(chan interface{}),
	}
	var err error

	if err = client.ConnectWs("public"); err != nil {
		return nil, err
	}

	go client.keepAliveLoop()

	return client, nil
}

func (client *Client) ConnectWs(publicPrivate string) error {
	var (
		connPtr **websocket.Conn
		url     string
	)

	if publicPrivate == "public" {
		connPtr = &client.publicWs
		url = publicWsURL
	} else {
		connPtr = &client.privateWs
		url = privateWsURL
	}

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return fmt.Errorf("could not connect to %s websocket: %w", publicPrivate, err)
	}

	*connPtr = conn

	go client.wsListener(publicPrivate, conn)
	go client.privateKeepAliveLoop()
	return nil
}

func (client *Client) SetVerbose(verbose bool) {
	client.verbose = verbose
}

func (client *Client) LoadWebsocketToken(key string, secret string) error {

	restClient := rest.NewClient()
	if err := restClient.SetAuth(key, secret); err != nil {
		return err
	}

	token, err := restClient.GetWebSocketsToken()
	if err != nil {
		return err
	}

	client.privateToken = token.Token
	return nil
}

func (client *Client) keepAliveLoop() {
	ticker := time.NewTicker(keepAliveDuration)
	for range ticker.C {
		if err := client.Send(Ping{}); err != nil {
			client.receiveChan <- fmt.Errorf("keep alive failed: %w", err)
			return
		}
	}
}

// TODO consider refactoring this
func (client *Client) privateKeepAliveLoop() {
	ticker := time.NewTicker(keepAliveDuration)
	for range ticker.C {
		if err := client.SendPrivate(Ping{}); err != nil {
			client.receiveChan <- fmt.Errorf("keep alive failed: %w", err)
			return
		}
	}
}

type DisconnectError struct {
	PublicPrivate string
	error
}

func (client *Client) wsListener(publicPrivate string, ws *websocket.Conn) {
	log.Printf("listening on %s websocket", publicPrivate)
	for {
		messageType, message, err := ws.ReadMessage()

		if err != nil {
			if _, ok := err.(*websocket.CloseError); ok {
				log.Printf("RECV %7s: dicconnect %s", publicPrivate, err.Error())

				client.receiveChan <- DisconnectError{error: err, PublicPrivate: publicPrivate}
				return
			}
			log.Printf("RECV %7s: error %s", publicPrivate, err.Error())

			client.receiveChan <- err
			continue
		}

		if messageType != websocket.TextMessage {
			client.receiveChan <- errBinaryMessage
			continue
		}

		if client.verbose {
			log.Printf("RECV %7s: %s", publicPrivate, string(message))
		}

		model, err := unmarshalReceivedMessage(message)

		if err != nil {
			client.receiveChan <- err
			continue
		}

		client.receiveChan <- model
	}
}

func (client *Client) Listen() <-chan interface{} {
	return client.receiveChan
}

func (client *Client) Send(rawMessage interface{}) error {
	return client.send(rawMessage, "public")
}

func (client *Client) SendPrivate(rawMessage interface{}) error {
	return client.send(rawMessage, "private")
}

func (client *Client) send(rawMessage interface{}, privatePublic string) error {

	doSend := func(message interface{}) error {
		bytes, err := json.Marshal(message)
		if err != nil {
			return err
		}

		if client.verbose {
			log.Printf("SEND %7s: %s", privatePublic, string(bytes))
		}

		if privatePublic == "public" {
			return client.publicWs.WriteJSON(message)
		}
		return client.privateWs.WriteJSON(message)
	}

	switch message := rawMessage.(type) {
	case Ping:
		message.Event = "ping"
		return doSend(message)
	case Subscribe:
		message.Event = "subscribe"
		if privatePublic == "private" {
			message.Subscription.Token = client.privateToken
		}
		return doSend(message)
	case Unsubscribe:
		message.Event = "unsubscribe"
		if privatePublic == "private" {
			message.Subscription.Token = client.privateToken
		}
		return doSend(message)
	case AddOrder:
		message.Event = "addOrder"
		message.Token = client.privateToken
		return doSend(message)
	case CancelOrder:
		message.Event = "cancelOrder"
		message.Token = client.privateToken
		return doSend(message)
	case CancelAll:
		message.Event = "cancelAll"
		message.Token = client.privateToken
		return doSend(message)
	default:
		return fmt.Errorf("unsupported message type %T", message)
	}
}

func Round(number float64, decimals int) string {
	format := fmt.Sprintf("%%.%df", decimals)
	return fmt.Sprintf(format, number)
}
