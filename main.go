package main

import (
	"os"
	"log"
	"encoding/csv"
	"strconv"
	"io"
	"fmt"
	"github.com/montanaflynn/stats"
	"time"
)

//StockRecord: Stock Record element 
type StockRecord struct {
	Date time.Time
	Open float64
}

func load() []StockRecord {
		// Open the file
	csvfile, err := os.Open("AI.PA.csv")
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
			log.Fatal(err)
		}
      // Parse the grade value.
      open, err := strconv.ParseFloat(record[1], 64)
      if err != nil {
          log.Fatal(err)
      }
		records = append(records, StockRecord{t, open})
		
	}
	//fmt.Println(records)
	return records
}
func regression(records []StockRecord) {
	// Convert raw data to serie
	var serie stats.Series
	for index := range records {
		serie = append(serie, stats.Coordinate{records[index].Date.Sub(records[0].Date).Seconds(), records[index].Open})		
	}

	reg, err := stats.LinearRegression(serie) 
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(reg)
}
func calc() {
	// start with some source data to use
	data := []float64{1.0, 2.1, 3.2, 4.823, 4.1, 5.8}

	// you could also use different types like this
	// data := stats.LoadRawData([]int{1, 2, 3, 4, 5})
	// data := stats.LoadRawData([]interface{}{1.1, "2", 3})
	// etc...

	median, _ := stats.Median(data)
	fmt.Println(median) // 3.65

	roundedMedian, _ := stats.Round(median, 0)
	fmt.Println(roundedMedian) // 4
}
func main() {
	fmt.Println("hello")
	data := load()
	regression(data)
}
