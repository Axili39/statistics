package marketstack

import (
	"github.com/Axili39/statistics/pkg/provider"
	"time"
	"net/http"
	"fmt"
	"database/sql"
	"encoding/json"
	"io/ioutil"
)

const urlBase string = "http://api.marketstack.com/v1/eod"

type MarketstackProvider struct {
	ApiKey string
}

func (p *MarketstackProvider) RetrieveData(exchange  string, symbol string, from time.Time, to time.Time) ([]provider.EodRecord, error) {
	// TODO if limit exceeds Limit, makes multiple query
	url := urlBase +
		"?access_key=" + p.ApiKey +
		"&symbols=" + exchange + symbol +
		"&sort=DESC" +
		"&date_from=" + from.Format("2010-01-02") +
		"&date_to=" + to.Format("2010-01-02") +
		"&limit=50000" +
		"&offset=0"
	
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(exchange, ":", symbol, ":", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Parse body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(exchange, ":", symbol, " : ", err)
		return  nil, err
	}
	var records []provider.EodRecord
	err = json.Unmarshal(body, &records)
	if err != nil {
		fmt.Println(exchange, ":", symbol, " : ", err)
		return nil, err
	}
	
	return records, nil
	// TODO unmarshall to EodRecord
}

func (p *MarketstackProvider) Setup(options string, db *sql.DB) error {
	m := provider.ParseOptions(options)
	var ok bool
	p.ApiKey, ok = m["apikey"]
	if ok == false {
		fmt.Errorf("api key option")
	}
	return nil
}