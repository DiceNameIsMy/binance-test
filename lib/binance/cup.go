package binance

import (
	"fmt"
	"log"
	"strconv"
)

type Cup struct {
	LastUpdateId int         `json:"lastUpdateId"`
	Bids         [][2]string `json:"bids"`
	Asks         [][2]string `json:"asks"`
}

func (c *Cup) CutData() {
	c.Bids = c.Bids[:15]
	c.Asks = c.Asks[:15]
}

func (c *Cup) GetTotalAsks() string {
	float_total := c.get_totals(c.Asks)
	total_ask := fmt.Sprintf("%f", float_total)
	return total_ask
}

func (c *Cup) GetTotalBids() string {
	float_total := c.get_totals(c.Bids)
	total_ask := fmt.Sprintf("%f", float_total)
	return total_ask
}

func (c *Cup) get_totals(data [][2]string) float64 {
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
