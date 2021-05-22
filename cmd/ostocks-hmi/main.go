package main

//go:generate res2go -o resources.go -package main -prefix rsrc visu.template

import (
	//"encoding/json"
	//"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
	"database/sql"

	"github.com/Axili39/statistics/pkg/database"
	"github.com/Axili39/statistics/pkg/fintools"
	"github.com/montanaflynn/stats"
	"github.com/julienschmidt/httprouter"
)

type Server struct {
	db *sql.DB
}
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
	db, err := database.Open("stocks.db")
	if err != nil {
		log.Fatal(err)
	}
	server := Server{db}
    router := httprouter.New()
    
	router.GET("/charts/stocks/:exchange/:symbol", server.stockCharts)
  	router.ServeFiles("/static/*filepath", http.Dir("./static"))
	
    http.ListenAndServe("0.0.0.0:8101", router)
}

func (server* Server)stockCharts(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Print("main")
	t, _ := template.New("main").Parse(string(rsrcFiles["visu.template"]))
	exchange := ps.ByName("exchange")

	symbol := ps.ByName("symbol")

	Values, Regression, StandardDev, _, _ := fintools.ComputeRegression(server.db, exchange, symbol, time.Now().AddDate(-10,0,0), time.Now())

	data := TemplateData{Title: exchange + ":" + symbol + " Evolution"}
	data.Series = append(data.Series, Serie{Name: "EOD", Values: Convert(time.Now().AddDate(-10,0,0), Values, 0)})
	data.Series = append(data.Series, Serie{Name: "Reg", Values: Convert(time.Now().AddDate(-10,0,0), Regression, 0)})
	data.Series = append(data.Series, Serie{Name: "-Std", Values: Convert(time.Now().AddDate(-10,0,0), Regression, StandardDev)})
	data.Series = append(data.Series, Serie{Name: "-2Std", Values: Convert(time.Now().AddDate(-10,0,0), Regression, StandardDev*2)})
	data.Series = append(data.Series, Serie{Name: "-Std", Values: Convert(time.Now().AddDate(-10,0,0), Regression, -1*StandardDev)})
	data.Series = append(data.Series, Serie{Name: "-2Std", Values: Convert(time.Now().AddDate(-10,0,0), Regression, -1*StandardDev*2)})

 	t.ExecuteTemplate(w, "main", data)
    
}