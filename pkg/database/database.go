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

func Update(db *sql.DB, p provider.StockProvider, exchange string, symbol string, from time.Time, to time.Time) error {
	//
	rows, err := db.Query("SELECT stock_id, description FROM stocks WHERE exchange=\""+ exchange +"\" AND symbol=\"" + symbol +"\"")
	if err != nil {
		return err
	}

	defer rows.Close()
	
	var	stockID int
	var desc   string	
	
	if rows.Next() {

		err = rows.Scan(&stockID, &desc)
		if err != nil {
			return err
		}		
	}
	rows.Close()

	fmt.Println("Loading stocks EX:", exchange, " SYMB :", symbol, " ", desc)
	err = updateInternal(db, p, stockID, exchange, symbol, from, to)
	if err != nil {
		return err
	}
	
	return nil
}

func updateInternal(db *sql.DB, p provider.StockProvider, stockID int, exchange string, symbol string, from time.Time, to time.Time) error {
	// Direct download
	data, err := p.RetrieveData(exchange, symbol, from, to)
	if err != nil {
		return err
	}

	// Prepare Db
	stmt, err := db.Prepare("INSERT INTO eod(stock_id, date, open, high, low, close, adj_close, volume) values(?,?,?,?,?,?,?,?)")
	if err != nil {
		return err
	}

	// Iterate through the records
	fmt.Println("\ttotal records :", len(data))
	for index, record := range data {

		res, err := stmt.Exec(stockID, record.Date, record.Open, record.High, record.Low, record.Close, record.AdjClose, record.Volume)
		if err != nil {
			fmt.Println(err, "  ...ignored")
			continue
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

func UpdateAll(db *sql.DB, p provider.StockProvider, from time.Time, to time.Time) error {
	//
	rows, err := db.Query("SELECT stock_id, exchange, symbol, description FROM stocks")
	if err != nil {
		return err
	}

	type record struct {
		stockID int
		exchange   string
		symbol string
		desc   string
	}
	var work []record
	for rows.Next() {
		var elem record
		err = rows.Scan(&elem.stockID, &elem.exchange, &elem.symbol, &elem.desc)
		if err != nil {
			return err
		}

		work = append(work, elem)

	}
	rows.Close()

	for _, r := range work {
		fmt.Println("Loading stocks EX:", r.exchange, " SYMB :", r.symbol, " ", r.desc)
		err = updateInternal(db, p, r.stockID, r.exchange, r.symbol, from, to)
		if err != nil {
			fmt.Println("Error when updating ", r.exchange, ":", r.stockID, err)
		}
	}

	return nil
}

// ImportExchanges : Import Stock Exchanges from Stock Market MBA csv format
func ImportExchanges(db *sql.DB, filename string) error {
	// Open the file
	csvfile, err := os.Open(filename)
	if err != nil {
		return err
	}

	defer csvfile.Close()

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
	stmt, err := db.Prepare("INSERT INTO exchanges ('BEC', 'BCC', 'country', 'description', 'MIC', 'GOOGLE', 'EOD') values(?,?,?,?,?,?,?)")
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
		stmt.Exec(record[0], record[1], record[2], record[3], record[4], record[5], record[6])
		if err != nil {

			return err
		}
	}
	return nil
}

// ImportStocksList : import stocks list from csv file to database
// csv must follow this format : ticker,name,exchange,category name,country
func ImportStocksList(db *sql.DB, exchange string, filename string) error {
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
	stmt, err := db.Prepare("INSERT INTO stocks('BBsymbol', 'description', 'exchange', 'symbol', 'IPO', 'ISIN', 'SEDOL') values(?,?,?,?,?,?,?)")
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
		stmt.Exec(record[0], record[1], exchange, record[2], record[3], record[8], record[9])
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
		"DROP TABLE IF EXISTS stocks",
		"DROP TABLE IF EXISTS exchanges",

		// Create Tables

		`CREATE TABLE "exchanges" (
    		"BEC"		STRING PRIMARY KEY,
    		"BCC" 		STRING 	NULL,
    		"country" 	STRING NULL,
    		"description" STRING NULL,
    		"MIC" 		STRING NULL,
    		"GOOGLE" 	STRING NULL,
    		"EOD" 		STRING NULL
		)`,

		`CREATE TABLE "eod" (
			"stock_id"  INTEGER,
    		"date" 		DATE,
    		"open" 		FLOAT64 NULL,
    		"close" 	FLOAT64 NULL,
    		"high" 		FLOAT64 NULL,
    		"low" 		FLOAT64 NULL,
    		"adj_close" FLOAT64 NULL,
    		"volume" 	FLOAT64 NULL,
    		PRIMARY KEY (stock_id, date)
		)`,

		`CREATE TABLE "stocks"(
			"stock_id"        INTEGER PRIMARY KEY AUTOINCREMENT,
			"BBsymbol"		STRING,
			"description"   STRING NULL,
			"exchange"		STRING,
			"symbol"    	STRING,
			"IPO"           DATE NULL,
			"ISIN" 			STRING NULL,
			"SEDOL"         STRING NULL,
			UNIQUE (exchange, symbol)
		)`,
	}

	err = batchRequestDB(db, script)
	if err != nil {
		return err
	}

	return nil
}
