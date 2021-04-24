package provider

import (
	"database/sql"
	"time"
	"strconv"
	"strings"
	"unicode"
)

type EodRecord struct {
	Exchange string		`json:"exchange"`
	Symbol   string     `json:"symbol"`
	Date     time.Time	`json:"date"`
	Open     float64	`json:"open"`
	High     float64	`json:"high"`
	Low      float64	`json:"low"`
	Close    float64	`json:"close"`
	AdjClose float64	`json:"adj_close"`
	Volume   float64	`json:"volume"`
}

type StockProvider interface {
	RetrieveData(exchange  string, symbol string, from time.Time, to time.Time) ([]EodRecord, error)
	Setup(options string, db *sql.DB) error
}

func EodRecordFromString(exchange string, symbol string, date string, open string, high string, low string, close string, adjClose string, volume string) (*EodRecord, error) { 
	var err error	
	record := EodRecord { Exchange: exchange, Symbol: symbol }

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

func ParseOptions(options string) map[string]string {
	// to generalize:
    lastQuote := rune(0)
    f := func(c rune) bool {
        switch {
        case c == lastQuote:
            lastQuote = rune(0)
            return false
        case lastQuote != rune(0):
            return false
        case unicode.In(c, unicode.Quotation_Mark):
            lastQuote = c
            return false
        default:
            return unicode.IsSpace(c)

        }
    }

    // splitting string by space but considering quoted section
    items := strings.FieldsFunc(options, f)

    // create and fill the map
    m := make(map[string]string)
    for _, item := range items {
        x := strings.Split(item, "=")
        m[x[0]] = x[1]
	}
	return m
}