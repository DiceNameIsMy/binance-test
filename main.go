package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"

	"github.com/sacOO7/gowebsocket"
)

type Cup struct {
	LastUpdateId int         `json:"lastUpdateId"`
	Bids         [][2]string `json:"bids"`
	Asks         [][2]string `json:"asks"`
}

func (c *Cup) cut_data() {
	c.Bids = c.Bids[:15]
	c.Asks = c.Asks[:15]
}

func main() {
	currency_symbol := "bnbbtc"
	params := "@depth20@1000ms"

	wsUrl := "wss://stream.binance.com:9443/ws/" + currency_symbol + params

	// channel to listen on code interruption
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	socket := gowebsocket.New(wsUrl)
	configure_default_socket(&socket)

	socket.OnTextMessage = func(message string, socket gowebsocket.Socket) {
		// log.Println("Recieved string message")

		var cup Cup
		if err := json.Unmarshal([]byte(message), &cup); err != nil {
			log.Fatal(err)
		}

		result := cup.parse_orders()

		log.Println(result)
	}

	socket.Connect()

	for {
		select {
		case <-interrupt:
			log.Println("interrupt")
			socket.Close()
			return
		}
	}

}

func configure_default_socket(s *gowebsocket.Socket) {

	s.OnConnected = func(socket gowebsocket.Socket) {
		log.Println("Connected to server")
	}

	s.OnConnectError = func(err error, socket gowebsocket.Socket) {
		log.Println("Recieved connect error ", err)
	}

	s.OnBinaryMessage = func(data []byte, socket gowebsocket.Socket) {
		log.Println("Recieved binary data ", data)
	}

	s.OnPingReceived = func(data string, socket gowebsocket.Socket) {
		log.Println("Recieved ping " + data)
	}

	s.OnPongReceived = func(data string, socket gowebsocket.Socket) {
		log.Println("Recieved pong " + data)
	}

	s.OnDisconnected = func(err error, socket gowebsocket.Socket) {
		log.Println("Disconnected from server ")
	}
}

// Returns averaged total
func (c Cup) parse_orders() [2]string {
	c.cut_data()

	float_total_bid := c.get_totals(c.Bids)
	str_total_bid := fmt.Sprintf("%f", float_total_bid)

	float_total_ask := c.get_totals(c.Asks)
	str_total_ask := fmt.Sprintf("%f", float_total_ask)

	return [2]string{str_total_bid, str_total_ask}
}

func (c Cup) get_totals(data [][2]string) float64 {
	var total float64
	for _, val := range data {
		q, err := strconv.ParseFloat(val[0], 64)
		if err != nil {
			log.Fatal(err)
		}
		p, err := strconv.ParseFloat(val[1], 64)
		if err != nil {
			log.Fatal(err)
		}

		total += q * p
	}
	return total
}
