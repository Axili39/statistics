package yahoo

import (
	"github.com/Axili39/statistics/pkg/provider"
	"time"
	"net/http"
	"fmt"
	"encoding/csv"
	"io"
)

const urlBase string = "https://query1.finance.yahoo.com/v7/finance/download/"

type YahooStockProvider struct {
	// TODO add a rate limiter
	lastRequest time.Time
}

func (p *YahooStockProvider) RetrieveData(ticker string, from time.Time, to time.Time) ([]provider.EodRecord, error) {
	if p.lastRequest.Add(time.Second*10).Second() > time.Now().Second() {
		// wait
		time.Sleep(10*time.Second)
	}
	url := urlBase + ticker +
		"?period1=" + fmt.Sprint(from.Unix()) +
		"&period2=" + fmt.Sprint(to.Unix()) +
		"&interval=1d" +
		"&events=history" +
		"&includeAdjustedClose=true"

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(ticker, ": ", err)
		return nil, err
	}
	defer resp.Body.Close()
	p.lastRequest = time.Now()

	// Parse the file
	r := csv.NewReader(resp.Body)

	// Read and ignore labels : Date,Open,High,Low,Close,Adj Close,Volume
	_, err = r.Read()
	if err == io.EOF {
		return nil, err
	}

	// Iterate through the records
	var records []provider.EodRecord
	for {
		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		eodRecord, err := provider.EodRecordFromString(ticker, record[0], record[1], record[2], record[3], record[4], record[5], record[6])
		if err != nil {
			continue
		}
		records = append(records, *eodRecord)
	}
	return records, nil
}