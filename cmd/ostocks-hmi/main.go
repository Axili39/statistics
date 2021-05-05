package main

//go:generate res2go -o resources.go -package main -prefix rsrc visu.template

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/Axili39/statistics/pkg/database"
	"github.com/Axili39/statistics/pkg/fintools"
)

func main() {
	rsrcInit()
	http.HandleFunc("/", mainPage)

    http.ListenAndServe("0.0.0.0:8100", nil)
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	log.Print("main")
	t, _ := template.New("main").Parse(string(rsrcFiles["visu.template"]))
	
	db, err := database.Open("")
	if err != nil {
		log.Fatal(err)
	}
	Serie, Regression, StandardDev, _, err := fintools.ComputeRegression(db, "FP", "AI", time.Now().AddDate(-1,0,0), time.Now())
	fmt.Print(Serie)
	fmt.Print(Regression)
	fmt.Print(StandardDev)
 	t.ExecuteTemplate(w, "main", nil)
    
}