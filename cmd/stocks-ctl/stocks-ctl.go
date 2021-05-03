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
	"github.com/Axili39/statistics/pkg/provider/dbfile"
	"github.com/Axili39/statistics/pkg/provider/marketstack"
	"github.com/Axili39/statistics/pkg/provider/openstock"
	"github.com/Axili39/statistics/pkg/provider/yahoo"
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

	err = database.ImportExchanges(db, "gbe.csv")
	if err != nil {
		log.Fatal(err)
	}
	err = database.ImportStocksList(db, "FP", "FP.csv")
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
	exchange := CommandLine.String("exchange", "", "exchange to update")
	symbol := CommandLine.String("symbol", "", "symbol to update")
	options := CommandLine.String("options", "", "option passed to provider")
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
		// Option rate limiter
	case "openstock":
		p = &openstock.OpenstockProvider{}
		// UrlBase: option "http://127.0.0.1:8080" mandatory
	case "dbfile":
		p = &dbfile.DBFileProvider{}
		// Filename : option mandatory
	case "marketstack":
		p = &marketstack.MarketstackProvider{}
		// ApiKey : option mandatory		
	default:
		fmt.Println("provider not yet implemented")
		os.Exit(1)
	}

	err = p.Setup(*options, db)
	if err != nil {
		fmt.Println("error when try to setup provider :", err)
	}
	fmt.Println("using provider :", *providers, "with configuration :", *options)

	if *exchange == "" || *symbol == "" {
		err = database.UpdateAll(db, p, time.Now().AddDate(*years*-1, *months*-1, *days*-1), time.Now())
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err = database.Update(db, p, *exchange, *symbol, time.Now().AddDate(*years*-1, *months*-1, *days*-1), time.Now())
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
	exchange := CommandLine.String("exchange", "", "exchange to update")
	symbol := CommandLine.String("symbol", "", "symbol to check")
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

	if *exchange == "" || *symbol == "" {
		fintools.FindCandidate(db, -1*(*years), -1*(*months), -1*(*days), *criteria)
	} else {
		fintools.CheckTicker(db, *exchange, *symbol, time.Now().AddDate(-1*(*years), -1*(*months), -1*(*days)), *criteria)
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
