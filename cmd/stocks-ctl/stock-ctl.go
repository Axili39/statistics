package main

import (
	"log"

	"github.com/Axili39/statistics/pkg/database"
)

func main() {
	err := database.CreateDB("truc.db")
	if err != nil {
		log.Fatal(err)
	}
}