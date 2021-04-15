package openstock

import (
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

func (p *OpenstockProvider) RetrieveData(ticker string, from time.Time, to time.Time) ([]provider.EodRecord, error) {
	url := p.UrlBase + "/api/v1/stocks/" + ticker +
		"?from=" + from.Format("2006-01-02") +
		"&to=" + to.Format("2006-01-02");
	fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(ticker, ": ", err)
		return nil, err
	}
	defer resp.Body.Close()
	
	// Parse body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(ticker, ": ", err)
		return  nil, err
	}

	var records []provider.EodRecord
	err = json.Unmarshal(body, &records)
	if err != nil {
		fmt.Println(ticker, ": ", err)
		return nil, err
	}
	
	return records, nil
}