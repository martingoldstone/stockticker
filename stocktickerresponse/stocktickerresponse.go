package stocktickerresponse

import (
	"fmt"
	"time"
)

type StockTickerDate time.Time

type StockTickerResponse struct {
	Data                []*StockData `json:"stockData"`
	AverageClosingPrice float64      `json:"averageClosingPrice"`
}

type StockData struct {
	date   StockTickerDate
	Date   string  `json:"date"`
	Open   float64 `json:"open"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Close  float64 `json:"close"`
	Volume int     `json:"volume"`
}

func (d StockTickerDate) String() string {
	return fmt.Sprint(time.Time(d).Format("2006-01-02"))
}

func NewStockData(date StockTickerDate, op, hp, lp, cp float64, vl int) *StockData {
	sd := StockData{Date: fmt.Sprint(date), Open: op, High: hp, Low: lp, Close: cp, Volume: vl}
	return &sd
}
