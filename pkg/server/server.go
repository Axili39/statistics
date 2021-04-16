package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Axili39/statistics/pkg/fintools"
	"github.com/Axili39/statistics/pkg/provider/dbfile"
	"github.com/julienschmidt/httprouter"
)

// implements local provider to act as a stocks provider
// local provider offer a HTTP/REST API
type Server struct {
	db *sql.DB
}

// Compute linear regression serie and find candidate
func (s *Server)getTickerData(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var err error
	// default period
	from := time.Now().AddDate(0,0,-1)
	to := time.Now()

	// ticker
	ticker := ps.ByName("ticker")

	// parse arguments
	query := r.URL.Query()
	fromStr := query.Get("from")
    if fromStr != "" {
		// convert from to time
		from, err = time.Parse("2006-01-02", fromStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	toStr := query.Get("to")
    if toStr != "" {
		// convert to to time
		to, err = time.Parse("2006-01-02", toStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	
	p := dbfile.DBFileProvider{s.db}
	data, err := p.RetrieveData(ticker, from, to)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, err.Error())
	}

	body, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, err.Error())
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, string(body))
}

func (s *Server)RegressionSeries(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var err error
	// default period
	from := time.Now().AddDate(0,0,-1)
	to := time.Now()

	// ticker
	ticker := ps.ByName("ticker")

	// parse arguments
	query := r.URL.Query()
	fromStr := query.Get("from")
    if fromStr != "" {
		// convert from to time
		from, err = time.Parse("2006-01-02", fromStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	toStr := query.Get("to")
    if toStr != "" {
		// convert to to time
		to, err = time.Parse("2006-01-02", toStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	reg, err = fintools.ComputeRegression(s.db, ticker, from, to)
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		fmt.FPrintln(w,err)
		return
	}
}

func Serve(db *sql.DB, uri string) error {
	server := Server{db}
    router := httprouter.New()
    
	router.GET("/api/v1/stocks/:ticker", server.getTickerData)
	router.GET("api/v1/regression/:ticker", server.RegressionSeries)

    return http.ListenAndServe(uri, router)
}

