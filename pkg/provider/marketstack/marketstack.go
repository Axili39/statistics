package marketstack

import (
	"github.com/Axili39/statistics/pkg/provider"
	"time"
	"net/http"
	"fmt"
	"encoding/csv"
	"io"
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

	// TODO unmarshall to EodRecord
}