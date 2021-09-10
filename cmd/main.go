package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"

	"github.com/DiceNameIsMy/binance-test/lib/binance"

	"github.com/sacOO7/gowebsocket"
)

func main() {
	currency_symbol := "btcusdt"
	params := "@depth20@1000ms"

	wsUrl := "wss://stream.binance.com:9443/ws/" + currency_symbol + params

	// channel to listen on code interruption
	// used to properly close connection on exit
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	socket := gowebsocket.New(wsUrl)
	configure_default_socket(&socket)

	socket.OnTextMessage = func(message string, socket gowebsocket.Socket) {
		var mcup binance.Cup
		if err := json.Unmarshal([]byte(message), &mcup); err != nil {
			log.Fatal(err)
		}
		mcup.Cut_data()

		total_bids := mcup.Get_total_bids()
		total_asks := mcup.Get_total_asks()

		log.Printf(" Bids: %s, Asks: %s", total_bids, total_asks)
	}

	socket.Connect()

	for {
		<-interrupt
		socket.Close()
		return
	}

}

func configure_default_socket(s *gowebsocket.Socket) {

	s.OnConnected = func(socket gowebsocket.Socket) {
		log.Println("Connected to server")
	}

	s.OnConnectError = func(err error, socket gowebsocket.Socket) {
		log.Println("Recieved connect error ", err)
	}

	s.OnDisconnected = func(err error, socket gowebsocket.Socket) {
		log.Println("Disconnected from server ")
	}
}
