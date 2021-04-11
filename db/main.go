package main

import (
	"net/http"
	"database/sql"
	"encoding/csv"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"log"
	//"os"
)
func anime(iterator int) {
	chars := []string{"-\b", "\\\b", "|\b", "/\b" }
	fmt.Printf(chars[iterator % 4])
}
func load(db *sql.DB, ticker string) {
	// direct download from yahoo
	url_base := "https://query1.finance.yahoo.com/v7/finance/download/"
	url := url_base + ticker + "?period1=946857600&period2=1618012800&interval=1d&events=history&includeAdjustedClose=true"
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

func main() {
	db, err := sql.Open("sqlite3", "./db.sql")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
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
