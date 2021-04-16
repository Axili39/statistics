package fintools
import (
	"database/sql"
	"fmt"
	"log"
	"time"
	"github.com/montanaflynn/stats"
	_ "github.com/mattn/go-sqlite3"
)

func ComputeRegression(db *sql.DB, ticker string, from time.Time, to time.Time) (stats.Series, error) {
	// load
	rows, err := db.Query("SELECT date, close FROM eod WHERE close <> \"null\" and ticker = \"" + ticker +
	 "\" and date > \"" + from.Format("2006-01-02") + 
	 "\" and date < \"" + to.Format("2006-01-02") + "\"")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()
	var serie stats.Series

	var date string
	var close float64
	for rows.Next() {
		err = rows.Scan(&date, &close)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		t, err := time.Parse("2006-01-02T00:00:00Z", date)
		if err != nil {
			fmt.Println(err)
			continue
		}
		serie = append(serie, stats.Coordinate{t.Sub(from).Seconds(), close})
	}

	// compute 
	reg, err := stats.LinearRegression(serie)
	if err != nil {
		log.Println("regression error:",err)
		return nil, err
	}
	return reg, nil
}

// Compute linear regression serie and find candidate
func CheckTicker(db *sql.DB, ticker string, after time.Time, criteria float64) {
	// load
	rows, err := db.Query("SELECT date, close FROM eod WHERE close <> \"null\" and ticker = \"" + ticker + "\" and date > \"" + after.Format("2006-01-02") + "\"")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()
	var serie stats.Series

	var date string
	var close float64
	for rows.Next() {
		err = rows.Scan(&date, &close)
		if err != nil {
			fmt.Println(err)
			return
		}

		t, err := time.Parse("2006-01-02T00:00:00Z", date)
		if err != nil {
			fmt.Println(err)
			continue
		}
		serie = append(serie, stats.Coordinate{t.Sub(after).Seconds(), close})
	}

	reg, err := stats.LinearRegression(serie)
	if err != nil {
		log.Println("regression error:",err)
		return
	}
	var sample stats.Float64Data
	for index := range reg {
		sample = append(sample, serie[index].Y-reg[index].Y)
	}
	stddev, _ := stats.StandardDeviation(sample)

	/*
	for index := range reg {
		deltaStd := (serie[index].Y-reg[index].Y)/stddev
		fmt.Println("date :", records[index].Date, " open:", serie[index].Y, " reg:", reg[index].Y, " dY:", serie[index].Y-reg[index].Y, " dStd:", deltaStd)
	}
	*/
	index := len(serie)-1
	deltaStd := (serie[index].Y-reg[index].Y)/stddev
	if deltaStd < criteria {
		fmt.Println(ticker, "date :", after.Add(time.Duration((serie[index].X))*time.Second), " open:", serie[index].Y, " reg:", reg[index].Y, " dY:", serie[index].Y-reg[index].Y, " dStd:", deltaStd, "standard dev:", stddev)
	} else {
		fmt.Println(ticker, "date :", after.Add(time.Duration((serie[index].X))*time.Second), " open:", serie[index].Y, " reg:", reg[index].Y, " dY:", serie[index].Y-reg[index].Y, " dStd:", deltaStd, "standard dev:", stddev)
	
	}
}

func FindCandidate(db *sql.DB, years int, months int, days int, criteria float64) {
	rows, err := db.Query("SELECT ticker, name  FROM tickers")
    if err != nil {
		log.Fatal(err)
	}
	var ticker string
	var name string
	for rows.Next() {
		err = rows.Scan(&ticker, &name)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("computing ", ticker, " ", name)
		CheckTicker(db, ticker, time.Now().AddDate(years,months,days), criteria)
	}
	defer rows.Close()
}