# Book example

This example shows how to:
* connect to the kraken websocket
* subscribe to a book currency topic
* update a `Book` state as updates come in

For show, we print the top 10 to stdout on each update.

In a real setting, you probably want to process messages from `Listen()` in a separate go-routine.
