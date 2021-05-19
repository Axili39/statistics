package main

//go:generate res2go -o resources.go -package main -prefix rsrc visu.template

import (
	//"encoding/json"
	//"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/Axili39/statistics/pkg/database"
	"github.com/Axili39/statistics/pkg/fintools"
	"github.com/montanaflynn/stats"
)

type DatedSample struct {
	X string
	Y float64
}
type Serie struct {
	Name string
	Values []DatedSample
}
type TemplateData struct {
	Title string
	Series []Serie
}
func Convert(start time.Time, serie stats.Series, shift float64) []DatedSample {
	ret := []DatedSample{}
	for _,e := range serie {
		t := start.AddDate(0,0,int(e.X))
		ret = append(ret, DatedSample{X: t.Format("2006-01-02"), Y: e.Y+shift})
	}
	return ret
}

func main() {
	rsrcInit()
  	http.Handle("/static/", http.FileServer(http.Dir("./")))
	http.HandleFunc("/main", mainPage)

    http.ListenAndServe("0.0.0.0:8100", nil)
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	log.Print("main")
	t, _ := template.New("main").Parse(string(rsrcFiles["visu.template"]))
	
	db, err := database.Open("stocks.db")
	if err != nil {
		log.Fatal(err)
	}
	Values, Regression, StandardDev, _, err := fintools.ComputeRegression(db, "FP", "AI", time.Now().AddDate(-10,0,0), time.Now())

	data := TemplateData{Title: "AI evolution"}
	data.Series = append(data.Series, Serie{Name: "EOD", Values: Convert(time.Now().AddDate(-10,0,0), Values, 0)})
	data.Series = append(data.Series, Serie{Name: "Reg", Values: Convert(time.Now().AddDate(-10,0,0), Regression, 0)})
	data.Series = append(data.Series, Serie{Name: "-Std", Values: Convert(time.Now().AddDate(-10,0,0), Regression, StandardDev)})
	data.Series = append(data.Series, Serie{Name: "-2Std", Values: Convert(time.Now().AddDate(-10,0,0), Regression, StandardDev*2)})
	data.Series = append(data.Series, Serie{Name: "-Std", Values: Convert(time.Now().AddDate(-10,0,0), Regression, -1*StandardDev)})
	data.Series = append(data.Series, Serie{Name: "-2Std", Values: Convert(time.Now().AddDate(-10,0,0), Regression, -1*StandardDev*2)})

 	t.ExecuteTemplate(w, "main", data)
    
}