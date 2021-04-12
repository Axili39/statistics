package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
	"github.com/montanaflynn/stats"
	_ "github.com/mattn/go-sqlite3"
	//"os"
)
func anime(iterator int) {
	chars := []string{"-\b", "\\\b", "|\b", "/\b" }
	fmt.Printf(chars[iterator % 4])
}
func load(db *sql.DB, ticker string) {
	// direct download from yahoo
	urlBase := "https://query1.finance.yahoo.com/v7/finance/download/"
	url := urlBase + ticker + "?period1=946857600&period2=1618012800&interval=1d&events=history&includeAdjustedClose=true"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(ticker, ": ", err)
		return
	}
	defer resp.Body.Close()

	// Parse the file
	r := csv.NewReader(resp.Body)

	// Read labels : Date,Open,High,Low,Close,Adj Close,Volume
	_, err = r.Read()
	if err == io.EOF {
		return
	}

	// Prepare Db

	stmt, err := db.Prepare("INSERT INTO eod(ticker, date, open, high, low, close, adj_close, volume) values(?,?,?,?,?,?,?,?)")
	if err != nil {
		fmt.Println("error db.Prepare")
		log.Fatal(err)
	}
	
	// Iterate through the records
	var record []string
	var last []string
	var id int64
	iterator := 0
	for {
		// Read each record from csv
		last = record
		record, err = r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		res, err := stmt.Exec(ticker, record[0], record[1], record[2], record[3], record[4], record[5], record[6])
		if err != nil {
			log.Fatal(err)
		}
		id, err = res.LastInsertId()
		if err != nil {
			log.Fatal(err)
		}

		if iterator == 0 {
			fmt.Println("\tfirst id: ", id, record)
		}
		//fmt.Printf(str(iterator))
		anime(iterator)
		iterator++
	}
	fmt.Println("\ttotal records :", iterator)
	fmt.Println("\tlast id: ", id, last)
}

func initdb(db *sql.DB) {

	//
	rows, err := db.Query("SELECT ticker, name  FROM tickers")
    if err != nil {
		log.Fatal(err)
	}

	var ticker string
	var name string

	type record struct {
		ticker string
		name string
	}
	var work []record
    for rows.Next() {
		err = rows.Scan(&ticker, &name)
		if err != nil {
			log.Fatal(err)
		}
		
		work = append(work, record{ticker, name})
		
	}
	rows.Close()

	for _,r := range work {
		fmt.Println("Loading ticker :", r.ticker, " ", r.name)
		load(db, r.ticker)
	}
}

// Compute linear regression serie and find candidate
func compute(db *sql.DB, ticker string, after time.Time, criteria float64) {
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
		fmt.Println("regression error:",err)
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
	}
}

func findCandidate(db *sql.DB) {
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
		fmt.Println("computing ", ticker)
		compute(db, ticker, time.Now().AddDate(-10,0,0), 20)
		
	}
	defer rows.Close()
}

func main() {
	db, err := sql.Open("sqlite3", "./db.sql")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	
	//initdb(db)
	findCandidate(db)
}
