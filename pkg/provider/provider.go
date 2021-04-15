package provider

import (
	"time"
	"strconv"
)

type EodRecord struct {
	Ticker   string		`json:"ticker"`
	Date     time.Time	`json:"date"`
	Open     float64	`json:"open"`
	High     float64	`json:"high"`
	Low      float64	`json:"low"`
	Close    float64	`json:"close"`
	AdjClose float64	`json:"adj_close"`
	Volume   float64	`json:"volume"`
}

type StockProvider interface {
	RetrieveData(ticker string, from time.Time, to time.Time) ([]EodRecord, error)
}

func EodRecordFromString(ticker string, date string, open string, high string, low string, close string, adjClose string, volume string) (*EodRecord, error) { 
	var err error	
	record := EodRecord { Ticker: ticker }

	record.Date, err = time.Parse("2006-01-02", date)
	if err != nil {
		return nil, err
	}
		 
	record.Open, err = strconv.ParseFloat(open, 64)
	if err != nil {
		return nil, err
	}

	record.High, err = strconv.ParseFloat(high, 64)
	if err != nil {
		return nil, err
	}
	
	record.Low, err = strconv.ParseFloat(low, 64)
	if err != nil {
		return nil, err
	}
	
	record.Close, err = strconv.ParseFloat(close, 64)
	if err != nil {
		return nil, err
	}

	record.AdjClose, err = strconv.ParseFloat(adjClose, 64)
	if err != nil {
		return nil, err
	}
	
	record.Volume, err = strconv.ParseFloat(volume, 64)
	if err != nil {
		return nil, err
	}
	
	return &record, nil
}

