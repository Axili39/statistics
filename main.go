package main

import (
	"encoding/csv"
	"fmt"
	"github.com/montanaflynn/stats"
	"io"
	"log"
	"os"
	"strconv"
	"time"
	"flag"
)

//StockRecord: Stock Record element
type StockRecord struct {
	Date time.Time
	Open float64
}

func load(ticker string, after time.Time) []StockRecord {
	// Open the file
	csvfile, err := os.Open(ticker + ".csv")
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}

	// Parse the file
	r := csv.NewReader(csvfile)

	// Read labels
	_, err = r.Read()
	if err == io.EOF {
		return nil
	}
	var records []StockRecord

	// Iterate through the records
	for {
		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		/*
			for index := range labels {
				fmt.Printf(" %s: %s", labels[index], record[index]);
			}
			fmt.Println()
		*/
		layout := "2006-01-02"
		t, err := time.Parse(layout, record[0])
		if err != nil {
			continue
		}
		// Parse the grade value.
		open, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			continue
		}
		if t.Sub(after).Seconds() > 0 {
			records = append(records, StockRecord{t, open})

		}

	}
	//fmt.Println(records)
	return records
}
func compute(ticker string, records []StockRecord) {
	// Convert raw data to serie
	var serie stats.Series
	for index := range records {
		serie = append(serie, stats.Coordinate{records[index].Date.Sub(records[0].Date).Seconds(), records[index].Open})
	}

	reg, err := stats.LinearRegression(serie)
	if err != nil {
		fmt.Println(err)
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
	index := len(records)-1
	deltaStd := (serie[index].Y-reg[index].Y)/stddev
	fmt.Println(ticker, "date :", records[index].Date, " open:", serie[index].Y, " reg:", reg[index].Y, " dY:", serie[index].Y-reg[index].Y, " dStd:", deltaStd, "standard dev:", stddev)
}

func main() {
	var ticker = flag.String("t", "ORA.PA", "ticker to check")
	flag.Parse()
	data := load(*ticker, time.Now().AddDate(-10, 0, 0))
	compute(*ticker, data)
}
