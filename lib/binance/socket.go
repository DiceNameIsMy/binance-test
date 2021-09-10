package binance

import (
	"encoding/json"
	"log"

	"github.com/sacOO7/gowebsocket"
)

const (
	url      = "wss://stream.binance.com:9443/ws/"
	cupParam = "@depth20"
)

func GetCupsSocket(currency string, params string) *gowebsocket.Socket {
	wsUrl := url + currency + cupParam + params

	socket := gowebsocket.New(wsUrl)
	configure_default_socket(&socket)

	socket.OnTextMessage = func(message string, socket gowebsocket.Socket) {
		var cup Cup
		if err := json.Unmarshal([]byte(message), &cup); err != nil {
			log.Fatal(err)
		}
		cup.CutData()

		total_bids := cup.GetTotalBids()
		total_asks := cup.GetTotalAsks()

		log.Printf(" Bids: %s, Asks: %s", total_bids, total_asks)
	}

	return &socket
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
