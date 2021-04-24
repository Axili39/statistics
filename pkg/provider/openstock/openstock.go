package openstock

import (
	"database/sql"
	"github.com/Axili39/statistics/pkg/provider"
	"time"
	"net/http"
	"fmt"
	"encoding/json"
	"io/ioutil"
)

type OpenstockProvider struct {
	UrlBase string
}

func (p *OpenstockProvider) RetrieveData(exchange string, symbol string, from time.Time, to time.Time) ([]provider.EodRecord, error) {
	url := p.UrlBase + "/api/v1/stocks/" + exchange + "/" + symbol +
		"?from=" + from.Format("2006-01-02") +
		"&to=" + to.Format("2006-01-02");
	fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(exchange, ":", symbol, " : ", err)
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
}

func (p *OpenstockProvider) Setup(options string, db *sql.DB) error {
	m := provider.ParseOptions(options)
	var ok bool
	p.UrlBase, ok = m["url"]
	if ok == false {
		fmt.Errorf("missing filename option")
	}
	return nil
}
