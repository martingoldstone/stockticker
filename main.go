package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/martingoldstone/stockticker/alphavantage"
	"github.com/martingoldstone/stockticker/stocktickerresponse"
)

var (
	Symbol    string
	Days      int
	Key       string
	SDays     string
	stckr     *stocktickerresponse.StockTickerResponse
	cachetime time.Time
	mu        sync.Mutex
)

func getFavicon(w http.ResponseWriter, req *http.Request) {
	log.Printf("%s %s %s Returning 404\n", req.Method, req.RequestURI, req.RemoteAddr)
	http.Error(w, "Not found", http.StatusNotFound)
	return
}

func timenow() string {
	return time.Now().Format(time.RFC3339)
}

func updateCache(s *stocktickerresponse.StockTickerResponse) {
	//Lock updating the cache variable in case multiple requests try to run before cache is initialised or after cache has expired
	mu.Lock()
	defer mu.Unlock()
	stckr = s
	cachetime = time.Now()
}

func getStockData(w http.ResponseWriter, req *http.Request) {
	log.Printf("%s %s %s Processing request...\n", req.RemoteAddr, req.Method, req.RequestURI)

	//Data is refreshed at end of trading time (16:00 US/Eastern), but fluctuations observed after end of trading so only cache for 1 hour
	//Alpha Vantage quotas requests per API key, caching necessary to avoid running out of quota
	if stckr == nil || time.Now().After(cachetime.Add(1*time.Hour)) {

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
		updateCache(resp.ToStockTickerResponse(Days))
		if stckr == nil {
			http.Error(w, "Failed to transcode API response from Alpha Vantage into response structure", http.StatusInternalServerError)
			return
		}

	} else {
		log.Printf("%s %s %s Serving from cache\n", req.RemoteAddr, req.Method, req.RequestURI)
	}

	rj, err := json.MarshalIndent(stckr, "", "    ")
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
