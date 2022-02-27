package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/martingoldstone/stockticker/alphavantage"
)

var (
	Symbol string
	Days   int
	Key    string
	SDays  string
)

func getFavicon(w http.ResponseWriter, req *http.Request) {
	log.Printf("%s %s %s Returning 404\n", req.Method, req.RequestURI, req.RemoteAddr)
	http.Error(w, "Not found", http.StatusNotFound)
	return
}

func getStockData(w http.ResponseWriter, req *http.Request) {
	log.Printf("%s %s %s Processing request...\n", req.Method, req.RequestURI, req.RemoteAddr)

	h, err := http.Get(fmt.Sprintf("https://www.alphavantage.co/query?apikey=%s&function=TIME_SERIES_DAILY&symbol=%s", Key, Symbol))
	if err != nil {
		log.Printf("Failed connecting to alpha vantage: %v\n", err)
		http.Error(w, "Failed opening connection to Alpha Vantage", http.StatusInternalServerError)
		return
	}

	defer h.Body.Close()

	if h.StatusCode != 200 {
		etxt, e := ioutil.ReadAll(h.Body)
		if e != nil {
			etxt = []byte("Could not decode response body")
		}
		http.Error(w, string(etxt), h.StatusCode)
		return
	}

	resp, err := alphavantage.ParseAlphaVantage(h.Body)
	if err != nil {
		http.Error(w, "Failed to decode API response from Alpha Vantage", http.StatusInternalServerError)
		return
	}
	str := resp.ToStockTickerResponse(Days)
	if str == nil {
		http.Error(w, "Failed to transcode API response from Alpha Vantage into response structure", http.StatusInternalServerError)
		return
	}

	rj, err := json.MarshalIndent(str, "", "    ")
	if err != nil {
		http.Error(w, "Could not render response structure to JSON", http.StatusInternalServerError)
		return
	}
	w.Write(rj)
	return

}

func main() {
	Key = os.Getenv("APIKEY")
	Symbol = os.Getenv("SYMBOL")
	SDays = os.Getenv("NDAYS")

	if Key == "" || Symbol == "" || SDays == "" {
		panic("Must set APIKEY, SYMBOL and NDAYS environment variables!")
	}

	var e error

	Days, e = strconv.Atoi(SDays)
	if e != nil {
		panic("NDAYS must be an int!")
	}

	http.HandleFunc("/", getStockData)
	http.HandleFunc("/favicon.ico", getFavicon)

	log.Print("Listening on 8080")

	http.ListenAndServe(":8080", nil)

}
