package database

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/Axili39/statistics/pkg/provider"
	_ "github.com/mattn/go-sqlite3"
)

func anime(iterator int) {
	chars := []string{"-\b", "\\\b", "|\b", "/\b"}
	fmt.Printf(chars[iterator%4])
}

func Update(db *sql.DB, p provider.StockProvider, ticker string, from time.Time, to time.Time) error {
	// Direct download
	data, err := p.RetrieveData(ticker, from, to)
	if err != nil {
		return err
	}

	// Prepare Db
	stmt, err := db.Prepare("INSERT INTO eod(ticker, date, open, high, low, close, adj_close, volume) values(?,?,?,?,?,?,?,?)")
	if err != nil {
		return err
	}

	// Iterate through the records
	fmt.Println("\ttotal records :", len(data))
	for index, record := range data {

		res, err := stmt.Exec(record.Ticker, record.Date, record.Open, record.High, record.Low, record.Close, record.AdjClose, record.Volume)
		if err != nil {
			return err
		}
		id, err := res.LastInsertId()
		if err != nil {
			return err
		}

		if index == 0 {
			fmt.Println("\tfirst id: ", id, record)
		} else if index == len(data)-1 {
			fmt.Println("\tlast id: ", id, record)
		}
		//fmt.Printf(str(iterator))
		anime(index)
	}
	stmt.Close()
	return nil
}

func UpdateAllTickers(db *sql.DB, p provider.StockProvider, from time.Time, to time.Time) error {
	//
	rows, err := db.Query("SELECT ticker, name  FROM tickers")
	if err != nil {
		return err
	}

	var ticker string
	var name string

	type record struct {
		ticker string
		name   string
	}
	var work []record
	for rows.Next() {
		err = rows.Scan(&ticker, &name)
		if err != nil {
			return err
		}

		work = append(work, record{ticker, name})

	}
	rows.Close()

	for _, r := range work {
		fmt.Println("Loading ticker :", r.ticker, " ", r.name)
		err = Update(db, p, r.ticker, from, to)
		if err != nil {
			return err
		}
	}

	return nil
}

// ImportTickerList : import tickers list from csv file to database
// csv must follow this format : ticker,name,exchange,category name,country
func ImportTickerList(db *sql.DB, filename string) error {
	// Open the file
	csvfile, err := os.Open(filename)
	if err != nil {
		return err
	}

	// Parse the file
	r := csv.NewReader(csvfile)
	r.Comma = ';'

	// Read and ignore labels
	_, err = r.Read()
	if err == io.EOF {
		return nil
	}
	if err != nil {
		return err
	}
	// Prepare Db
	stmt, err := db.Prepare("INSERT INTO tickers('ticker','name','exchange','category name','country') values(?,?,?,?,?)")
	if err != nil {
		return err
	}

	// Iterate through the records
	for {
		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		stmt.Exec(record[0], record[1], record[2], record[3], record[4])
		if err != nil {

			return err
		}
	}
	return nil
}

func batchRequestDB(db *sql.DB, requests []string) error {
	for _, r := range requests {
		stm, err := db.Prepare(r)
		if err != nil {
			return err
		}
		_, err = stm.Exec()
		if err != nil {
			return err
		}
	}
	return nil
}

func Open(file string) (*sql.DB, error) {
	return sql.Open("sqlite3", file)
}
func Create(file string) error {
	// Open file
	db, err := Open(file)
	if err != nil {
		return err
	}
	defer db.Close()

	script := []string{
		// Clean all tables
		"DROP TABLE IF EXISTS eod",
		"DROP TABLE IF EXISTS tickers",

		// Create Tables
		`CREATE TABLE "eod" (
    		"uid" 		INTEGER PRIMARY KEY AUTOINCREMENT,
    		"ticker" 	STRING 	NULL,
    		"date" 		DATE 	NULL,
    		"open" 		FLOAT64 NULL,
    		"close" 	FLOAT64 NULL,
    		"high" 		FLOAT64 NULL,
    		"low" 		FLOAT64 NULL,
    		"adj_close" FLOAT64 NULL,
    		"volume" 	FLOAT64 NULL,
    		UNIQUE(ticker, date)
		)`,
		`CREATE TABLE "tickers"(
    		"ticker" 		STRING PRIMARY KEY,
    		"name" 			STRING NULL,
    		"exchange" 		STRING NULL, 
    		"category name" STRING NULL, 
    		"country" 		STRING NULL
		)`,
	}

	err = batchRequestDB(db, script)
	if err != nil {
		return err
	}

	return nil
}
