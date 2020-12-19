package main

import (
	"log"

	"github.com/lk16/kraken/websocket"
)

func main() {

	client, err := websocket.NewClient()
	if err != nil {
		panic(err)
	}

	subscribe := websocket.Subscribe{
		Pair:         []string{"XRP/EUR"},
		Subscription: websocket.Subscription{Name: "book"},
	}

	if err = client.Send(subscribe); err != nil {
		panic(err)
	}

	var book websocket.Book

	for rawMessage := range client.Listen() {
		switch message := rawMessage.(type) {
		case websocket.SubscriptionStatus, websocket.SystemStatus, websocket.HeartBeat, websocket.Pong:
			// do nothing
		case websocket.Book:
			book = message
			book.PrintTop(10)
		case websocket.BookUpdate:
			book.Update(message)
			book.PrintTop(10)
		case error:
			log.Fatalf("got err %T %s", message, message.Error())
		default:
			log.Fatalf("got %+#v", message)
		}
	}
}
