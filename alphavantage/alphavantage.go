package alphavantage

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"sort"
	"strconv"
	"time"

	"github.com/martingoldstone/stockticker/stocktickerresponse"
)

type ApiResponse struct {
	Meta       ApiResponseMeta           `json:"Meta Data"`
	TimeSeries map[string]TimeSeriesData `json:"Time Series (Daily)"`
}

type ApiResponseMeta struct {
	LastRefreshed string `json:"3. Last Refreshed"`
	TimeZone      string `json:"5. Time Zone"`
}

type TimeSeriesData struct {
	Open   string `json:"1. open"`
	High   string `json:"2. high"`
	Low    string `json:"3. low"`
	Close  string `json:"4. close"`
	Volume string `json:"5. volume"`
}

func (a *ApiResponse) ToStockTickerResponse(d int) *stocktickerresponse.StockTickerResponse {
	r := stocktickerresponse.StockTickerResponse{}
	td := len(a.TimeSeries)
	if d < td {
		td = d
	}
	var days []string

	for k := range a.TimeSeries {
		if _, e := time.Parse("2006-01-02", k); e == nil {
			days = append(days, k)
		}
	}
	sort.Sort(sort.Reverse(sort.StringSlice(days)))
	var data []*stocktickerresponse.StockData
	for _, v := range days[:td] {
		rd, _ := time.Parse("2006-01-02", v)
		op, err := strconv.ParseFloat(a.TimeSeries[v].Open, 64)
		if err != nil {
			log.Panicf("Failed to parse opening price %v for date %s\n", a.TimeSeries[v].Open, v)
		}
		hp, err := strconv.ParseFloat(a.TimeSeries[v].High, 64)
		if err != nil {
			log.Panicf("Failed to parse high price %v for date %s\n", a.TimeSeries[v].High, v)
		}
		lp, err := strconv.ParseFloat(a.TimeSeries[v].Low, 64)
		if err != nil {
			log.Panicf("Failed to parse low price %v for date %s\n", a.TimeSeries[v].Low, v)
		}
		cp, err := strconv.ParseFloat(a.TimeSeries[v].Close, 64)
		if err != nil {
			log.Panicf("Failed to parse closing price %v for date %s\n", a.TimeSeries[v].Close, v)
		}
		vl, err := strconv.Atoi(a.TimeSeries[v].Volume)
		if err != nil {
			log.Panicf("Failed to parse volume %v for date %s\n", a.TimeSeries[v].Volume, v)
		}
		datum := stocktickerresponse.NewStockData(stocktickerresponse.StockTickerDate(rd), op, hp, lp, cp, vl)
		data = append(data, datum)
	}

	tc := float64(0)
	for _, v := range data {
		tc += v.Close
	}
	r.AverageClosingPrice = tc / float64(len(data))
	r.Data = data
	return &r
}

func ParseAlphaVantage(f io.Reader) (*ApiResponse, error) {
	r := ApiResponse{}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(b, &r); err != nil {
		log.Println(err)
		return nil, err
	}
	return &r, nil
}
