package fintools
import (
	"database/sql"
	"fmt"
	"log"
	"time"
	"github.com/montanaflynn/stats"
	_ "github.com/mattn/go-sqlite3"
)

func ComputeRegression(db *sql.DB, exchange string, symbol string, from time.Time, to time.Time) (stats.Series, stats.Series, float64, *time.Time, error) {
	// load
	rows, err := db.Query("SELECT eod.date, eod.close FROM eod JOIN stocks WHERE eod.stock_id = stocks.stock_id and eod.close <> \"null\" and stocks.exchange = \"" + exchange + "\" and stocks.symbol=\"" + symbol +
	 "\" and eod.date >= \"" + from.Format("2006-01-02") + 
	 "\" and eod.date <= \"" + to.Format("2006-01-02") + "\"")
	if err != nil {
		fmt.Println(err)
		return nil, nil, 0.0, nil, err
	}
	defer rows.Close()
	var serie stats.Series

	var start* time.Time
	var date time.Time
	var close float64
	for rows.Next() {
		err = rows.Scan(&date, &close)
		if err != nil {
			fmt.Println(err)
			return nil, nil, 0.0, nil, err
		}
/*
		t, err := time.Parse("2006-01-02T00:00:00Z", date)
		if err != nil {
			fmt.Println(err)
			continue
		}
		*/
		if start == nil {
			start = &date
		}
		serie = append(serie, stats.Coordinate{date.Sub(from).Seconds()/(24*3600), close})
	}

	// compute 
	reg, err := stats.LinearRegression(serie)
	if err != nil {
		log.Println("regression error:",err)
		return nil, nil, 0.0, nil, err
	}

	// Compute Standard Deviation
	var sample stats.Float64Data
	for index := range reg {
		sample = append(sample, serie[index].Y-reg[index].Y)
	}
	stddev, _ := stats.StandardDeviation(sample)
	
	return serie, reg, stddev, start, nil
}

// Compute linear regression serie and find candidate
func CheckTicker(db *sql.DB, exchange string, symbol string, after time.Time, criteria float64) {
	
	// Compute regression
	serie, reg, stddev, start, err := ComputeRegression(db, exchange, symbol, after, time.Now())
	if err != nil {
		fmt.Println(err)
		return
	}
	index := len(serie)-1
	deltaStd := (serie[index].Y-reg[index].Y)/stddev
	if deltaStd < criteria {
		fmt.Println(exchange, ":", symbol, "date :", *start, " open:", serie[index].Y, " reg:", reg[index].Y, " dY:", serie[index].Y-reg[index].Y, " dStd:", deltaStd, "standard dev:", stddev)
	} else {
		//fmt.Println(exchange, ":", symbol, "date :", after.Add(time.Duration((serie[index].X))*time.Second), " open:", serie[index].Y, " reg:", reg[index].Y, " dY:", serie[index].Y-reg[index].Y, " dStd:", deltaStd, "standard dev:", stddev)	
	}
}

func FindCandidate(db *sql.DB, years int, months int, days int, criteria float64) {
	rows, err := db.Query("SELECT exchange, symbol, description FROM stocks")
    if err != nil {
		log.Fatal(err)
	}
	var exchange string
	var symbol string
	var description string
	for rows.Next() {
		err = rows.Scan(&exchange, &symbol, &description)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("computing ", exchange, ":", symbol, " :", description)
		CheckTicker(db, exchange, symbol, time.Now().AddDate(years,months,days), criteria)
	}
	defer rows.Close()
}