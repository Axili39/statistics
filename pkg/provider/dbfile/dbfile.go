package dbfile

import (
	"github.com/Axili39/statistics/pkg/provider"
	"time"
	"fmt"
	"database/sql"
)

const urlBase string = "https://query1.finance.yahoo.com/v7/finance/download/"

type DBFile struct {
	Db *sql.DB
}

func (p *DBFile) RetrieveData(ticker string, from time.Time, to time.Time) ([]provider.EodRecord, error) {
	// load
	rows, err := p.Db.Query("SELECT ticker, date, open, high, low, close, adj_close, volume FROM eod WHERE ticker = \"" + ticker + "\" and date >= \"" + from.Format("2006-01-02") + "\" and date <= \"" + to.Format("2006-01-02") + "\"")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()


	var records []provider.EodRecord
	
	for rows.Next() {
		var record provider.EodRecord
		err = rows.Scan(&record.Ticker, &record.Date, &record.Open, &record.High, &record.Low, &record.Close, &record.AdjClose, &record.Volume)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		records = append(records, record)
	}
	return records, nil
}