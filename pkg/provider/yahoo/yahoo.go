package yahoo

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/Axili39/statistics/pkg/provider"
)

const urlBase string = "https://query1.finance.yahoo.com/v7/finance/download/"

type YahooStockProvider struct {
	// TODO add a rate limiter
	lastRequest time.Time
	period int
	xchgMap map [string]string
}

func (p *YahooStockProvider) getTicker(exchange string, symbol string) string {
	return symbol + p.xchgMap[exchange]
}

func (p *YahooStockProvider) RetrieveData(exchange string, symbol string, from time.Time, to time.Time) ([]provider.EodRecord, error) {
	if  time.Now().Second() - p.lastRequest.Second() < p.period {
		// wait
		time.Sleep(time.Duration(p.period)*time.Second)
	}
	ticker := p.getTicker(exchange, symbol)
	url := urlBase + ticker +
		"?period1=" + fmt.Sprint(from.Unix()) +
		"&period2=" + fmt.Sprint(to.Unix()) +
		"&interval=1d" +
		"&events=history" +
		"&includeAdjustedClose=true"
	p.lastRequest = time.Now()

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(ticker, ": ", err)
		return nil, err
	}
	defer resp.Body.Close()
	

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

		eodRecord, err := provider.EodRecordFromString(exchange, symbol, record[0], record[1], record[2], record[3], record[4], record[5], record[6])
		if err != nil {
			continue
		}
		records = append(records, *eodRecord)
	}
	return records, nil
}

func (p *YahooStockProvider) Setup(options string, db *sql.DB) error {
	m := provider.ParseOptions(options)
	period, ok := m["period"]
	if ok == true {
		var err error
		p.period, err = strconv.Atoi(period)
	
		if err != nil {
			fmt.Println("bad value for period ", period)
		}
	}
	p.lastRequest = time.Now().Add(time.Duration(-1*p.period)*time.Second)

	// fill exchange map
	rows, err := db.Query("SELECT BEC, EOD FROM exchanges")
	if err != nil {
		return err
	}
	p.xchgMap = make(map[string]string)
	for rows.Next() {
		var bec string
		var eod string
		err = rows.Scan(&bec, &eod)
		if err != nil {
			fmt.Println(err)
			return err
		}
		p.xchgMap[bec] = "." + eod
	}
	return nil
}