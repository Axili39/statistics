package main

import (
	//"flag"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Axili39/statistics/pkg/database"
	"github.com/Axili39/statistics/pkg/server"
	"github.com/Axili39/statistics/pkg/fintools"
	"github.com/Axili39/statistics/pkg/provider"
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

func update(args []string) {
	fmt.Println(args)
	// parse command
	CommandLine := flag.NewFlagSet("update", flag.ExitOnError)
	years := CommandLine.Int("years", 0, "how many years to update")
	months := CommandLine.Int("months", 0, "how many months to update")
	days := CommandLine.Int("days", 0, "how many days to update")
	providers := CommandLine.String("provider", "yahoo", "provider to query (yahoo, xxxxx, openstocks")
	err := CommandLine.Parse(args)	
	if err != nil {
		fmt.Println("error")
	}

	if len(args) > 0 && args[1] == "help" {
		fmt.Println("stocks-ctl update :")
		CommandLine.PrintDefaults()	
		os.Exit(0)
	}
	
	fmt.Println(*years)

	// Open file
	db, err := database.Open(dbfile)
	if err != nil {
		log.Fatal(err)
	}

	var p provider.StockProvider
	switch *providers {
	case "yahoo":
		p = &yahoo.YahooStockProvider{}
	default:
		fmt.Println("provider not yet implemented")
		os.Exit(1)
	}

	err = database.UpdateAllTickers(db, p, time.Now().AddDate(*years*-1, *months*-1, *days*-1), time.Now())
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

func serve() {
	// Open file
	db, err := database.Open(dbfile)
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
		compute()
	case "server":
		serve()
	default:
		usage()
	}
}
