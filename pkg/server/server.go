package server

import (
	"time"
	"fmt"
	"net/http"
	"database/sql"
	"github.com/julienschmidt/httprouter"
	"github.com/Axili39/statistics/pkg/provider/dbfile"
)

// implements local provider to act as a stocks provider
// local provider offer a HTTP/REST API
type Server struct {
	db *sql.DB
}

// Compute linear regression serie and find candidate
func (s *Server)getTickerData(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// default period
	from := time.Now().AddDate(0,0,-1)
	to := time.Now()

	// ticker
	ticker := ps.ByName("ticker")

	// parse arguments
	query := r.URL.Query()
	fromStr, present := query["from"]
    if present && len(fromStr) > 0 {
		// convert from to time
	}

	toStr, present := query["to"]
    if present && len(toStr) > 0 {
		// convert to to time
	}
	
	
	p := dbfile.DBFile{s.db}
	data, err := p.RetrieveData(ticker, from, to)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, err.Error())
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "getting %s from:%s to: %s !\n", ticker, from, to)
	fmt.Fprintln(w, data)
}

func Serve(db *sql.DB, uri string) error {
	server := Server{db}
    router := httprouter.New()
    
    router.GET("/api/v1/stocks/:ticker", server.getTickerData)

    return http.ListenAndServe(uri, router)
}

