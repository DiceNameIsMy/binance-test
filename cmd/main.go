package main

import (
	"flag"
	"os"
	"os/signal"

	"github.com/DiceNameIsMy/binance-test/lib/binance"
)

func main() {
	currency_symbol := flag.String(
		"symbol",
		"btcusdt",
		"currencies to track as: `btcusdt`",
	)
	params := flag.String(
		"params",
		"@1000ms",
		"params written as: `@depth20@1000ms`",
	)
	// channel to listen on code interruption
	// used to properly close connection on exit
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	socket := binance.GetCupsSocket(*currency_symbol, *params)

	socket.Connect()

	for {
		<-interrupt
		socket.Close()
		return
	}

}
