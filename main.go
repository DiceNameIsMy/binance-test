package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/sacOO7/gowebsocket"
)

func main() {
	// url := "https://api.binance.com/api/v3/depth?symbol=BNBBTC&limit=20"
	wsUrl := "wss://stream.binance.com:9443/ws/btcusdt@depth20@1000ms"

	data_filename := "test1.json"
	first_item := true

	err := os.WriteFile(data_filename, []byte("["), 0644)
	if err != nil {
		log.Fatal(err)
	}
	file, err := os.OpenFile(data_filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	socket := gowebsocket.New(wsUrl)

	socket.OnConnected = func(socket gowebsocket.Socket) {
		log.Println("Connected to server")
	}

	socket.OnConnectError = func(err error, socket gowebsocket.Socket) {
		log.Println("Recieved connect error ", err)
	}

	socket.OnTextMessage = func(message string, socket gowebsocket.Socket) {
		log.Println("Recieved string message")
		extra_str := ""
		if !first_item {
			extra_str = ","
		}
		if _, err := file.WriteString(extra_str + message); err != nil {
			log.Fatal(err)
		}
		first_item = false
	}

	socket.OnBinaryMessage = func(data []byte, socket gowebsocket.Socket) {
		log.Println("Recieved binary data ", data)
	}

	socket.OnPingReceived = func(data string, socket gowebsocket.Socket) {
		log.Println("Recieved ping " + data)
	}

	socket.OnPongReceived = func(data string, socket gowebsocket.Socket) {
		log.Println("Recieved pong " + data)
	}

	socket.OnDisconnected = func(err error, socket gowebsocket.Socket) {
		log.Println("Disconnected from server ")
	}

	socket.Connect()

	for {
		select {
		case <-interrupt:
			log.Println("interrupt")
			file.WriteString("]")
			socket.Close()
			return
		}
	}

}
