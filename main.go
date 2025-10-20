package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
)

func defaultTheme() opts.Initialization {
	return opts.Initialization{
		Theme: types.ThemeVintage,
		// BackgroundColor: "white",
	}
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func main() {
	var (
		generateReports bool
		outputDir       string
		year, month     int
		server          bool
	)

	flag.BoolVar(&generateReports, "gen", false, "generate reports")
	flag.StringVar(&outputDir, "out", "", "output directory")
	flag.IntVar(&year, "year", time.Now().Year(), "year")
	flag.IntVar(&month, "month", -1, "month. If not provided, reports for all months of the year will be generated")
	flag.BoolVar(&server, "server", false, "start file server")
	flag.Parse()

	if outputDir == "" {
		log.Fatal("output directory must be specified")
	}

	if generateReports {
		c := NewCharts(outputDir, year)
		c.createYoYReport()
		c.createYearlyReport()

		if month == -1 {
			c.generateMonthlyReports(year)
		} else {
			t := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
			c.createMonthlyReport(t)
		}
	}

	if server {
		fs := http.FileServer(http.Dir(outputDir))
		log.Println("running server at http://localhost:8081")
		log.Fatal(http.ListenAndServe("localhost:8081", logRequest(noCacheWrapper(fs))))
	}
}

func noCacheWrapper(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Disable caching
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		w.Header().Set("Surrogate-Control", "no-store")

		// Serve the file
		h.ServeHTTP(w, r)
	})
}
