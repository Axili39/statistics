package main

import (
	"log"
	"time"

	"github.com/Axili39/statistics/pkg/database"
	"github.com/Axili39/statistics/pkg/fintools"
	"github.com/Axili39/statistics/pkg/provider/yahoo"
)
const dbfile string = "stocks.db"

func create() {
	err := database.Create(dbfile)
	if err != nil {
		log.Fatal(err)
	}
}

func importTickers() {
	// Open file
	db, err := database.Open(dbfile)
	if err != nil {
		log.Fatal(err)
	}
	
	err = database.ImportTickerList(db, "tickers.csv")
	if err != nil {
		log.Fatal(err)
	}
}

func update() {
// Open file
	db, err := database.Open(dbfile)
	if err != nil {
		log.Fatal(err)
	}

	p := &yahoo.YahooStockProvider{}

	err = database.UpdateAllTickers(db, p, time.Now().AddDate(0,-1,0), time.Now())
	if err != nil {
		log.Fatal(err)
	}
}

func compute() {
	// Open file
	db, err := database.Open(dbfile)
	if err != nil {
		log.Fatal(err)
	}

	fintools.FindCandidate(db)
}

func main() {
	create()
	importTickers()
	update()
	compute()
}