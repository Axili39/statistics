package main

import (
	//"flag"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Axili39/statistics/pkg/database"
	"github.com/Axili39/statistics/pkg/fintools"
	"github.com/Axili39/statistics/pkg/provider"
	"github.com/Axili39/statistics/pkg/provider/openstock"
	"github.com/Axili39/statistics/pkg/provider/yahoo"
	"github.com/Axili39/statistics/pkg/provider/dbfile"
	"github.com/Axili39/statistics/pkg/server"
)

const dbfilename string = "stocks.db"

func create() {
	err := database.Create(dbfilename)
	if err != nil {
		log.Fatal(err)
	}
}

func importTickers() {
	// Open file
	db, err := database.Open(dbfilename)
	if err != nil {
		log.Fatal(err)
	}

	err = database.ImportTickerList(db, "tickers.csv")
	if err != nil {
		log.Fatal(err)
	}
}

func update(args []string) {
	// parse command
	CommandLine := flag.NewFlagSet("update", flag.ExitOnError)
	years := CommandLine.Int("years", 0, "how many years to update")
	months := CommandLine.Int("months", 0, "how many months to update")
	days := CommandLine.Int("days", 0, "how many days to update")
	providers := CommandLine.String("provider", "yahoo", "provider to query (yahoo, xxxxx, openstocks")
	ticker := CommandLine.String("ticker", "", "ticker to update")
	err := CommandLine.Parse(args)	
	if err != nil {
		fmt.Println("error")
	}

	if len(args) > 0 && args[0] == "help" {
		fmt.Println("stocks-ctl update :")
		CommandLine.PrintDefaults()	
		os.Exit(0)
	}
	
	fmt.Println(*years)

	// Open file
	db, err := database.Open(dbfilename)
	if err != nil {
		log.Fatal(err)
	}

	var p provider.StockProvider
	switch *providers {
	case "yahoo":
		p = &yahoo.YahooStockProvider{}
	case "openstock":
		p = &openstock.OpenstockProvider{UrlBase: "http://127.0.0.1:8080"}
	case "dbfile":
		p = &dbfile.DBFileProvider{Db: db}
	default:
		fmt.Println("provider not yet implemented")
		os.Exit(1)
	}

	if *ticker == "" {
		err = database.UpdateAllTickers(db, p, time.Now().AddDate(*years*-1, *months*-1, *days*-1), time.Now())
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err = database.Update(db, p, *ticker, time.Now().AddDate(*years*-1, *months*-1, *days*-1), time.Now())
		if err != nil {
			log.Fatal(err)
		}
	}
}

func compute(args []string) {
	// parse command
	CommandLine := flag.NewFlagSet("compute", flag.ExitOnError)
	years := CommandLine.Int("years", 0, "how many years to update")
	months := CommandLine.Int("months", 0, "how many months to update")
	days := CommandLine.Int("days", 0, "how many days to update")
	criteria := CommandLine.Float64("criteria", 0, "count of std deviation under regression line to select ticker")
	ticker := CommandLine.String("ticker", "", "ticker to update")
	err := CommandLine.Parse(args)	
	if err != nil {
		fmt.Println("error")
	}

	if len(args) > 0 && args[0] == "help" {
		fmt.Println("stocks-ctl compute :")
		CommandLine.PrintDefaults()	
		os.Exit(0)
	}
	// Open file
	db, err := database.Open(dbfilename)
	if err != nil {
		log.Fatal(err)
	}

	if *ticker == "" {
		fintools.FindCandidate(db, -1*(*years), -1*(*months), -1*(*days), *criteria)
	} else {
		fintools.CheckTicker(db, *ticker, time.Now().AddDate(-1*(*years), -1*(*months), -1*(*days)), *criteria)
	}
}

func serve() {
	// Open file
	db, err := database.Open(dbfilename)
	if err != nil {
		log.Fatal(err)
	}

	server.Serve(db, ":8080")
}

func usage() {
	fmt.Println("usage: stocks-ctl help|create|import|update|compute|serve")
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	cmd := os.Args[1]

	switch cmd {
	case "help":
		usage()
	case "create":
		create()
	case "import":
		importTickers()
	case "update":
		update(os.Args[2:])
	case "compute":
		compute(os.Args[2:])
	case "server":
		serve()
	default:
		usage()
	}
}
