package dbfile

import (
	"database/sql"
	"fmt"
	"time"


	"github.com/Axili39/statistics/pkg/database"
	"github.com/Axili39/statistics/pkg/provider"
)

type DBFileProvider struct {
	Db *sql.DB
}

func (p *DBFileProvider) RetrieveData(exchange  string, symbol string, from time.Time, to time.Time) ([]provider.EodRecord, error) {
	// load
	rows, err := p.Db.Query("SELECT stocks.exchange, stocks.symbol, eod.date, eod.open, eod.high, eod.low, eod.close, eod.adj_close, eod.volume FROM eod JOIN stocks " + 
	"WHERE eod.stock_id = stocks.stock_id and stocks.exchange = \"" + exchange + "\" and stocks.symbol =\"" + symbol + 
	"\" and eod.date >= \"" + from.Format("2006-01-02") + "\" and eod.date <= \"" + to.Format("2006-01-02") + "\"")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	var records []provider.EodRecord

	for rows.Next() {
		var record provider.EodRecord
		err = rows.Scan(&record.Exchange, &record.Symbol, &record.Date, &record.Open, &record.High, &record.Low, &record.Close, &record.AdjClose, &record.Volume)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		records = append(records, record)
	}
	return records, nil
}

func (p *DBFileProvider) Setup(options string, db *sql.DB) error {
	m := provider.ParseOptions(options)
	file, ok := m["file"]
	if ok == true {
		var err error
		p.Db, err = database.Open(file)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("missing filename option")
	}
	return nil
}
