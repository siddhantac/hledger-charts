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
		Theme:           types.ThemeRoma,
		BackgroundColor: "white",
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
		year            int
	)

	flag.BoolVar(&generateReports, "gen", false, "generate reports")
	flag.StringVar(&outputDir, "out", "", "output directory")
	flag.IntVar(&year, "year", time.Now().Year(), "year")
	flag.Parse()

	if outputDir == "" {
		log.Fatal("output directory must be specified")
	}

	if generateReports {
		c := NewCharts(outputDir, year)
		c.createYearlyReport()
		c.createMonthlyReport("last month")
		c.generateMonthlyReports(year)
	}

	fs := http.FileServer(http.Dir(outputDir))
	log.Println("running server at http://localhost:8081")
	log.Fatal(http.ListenAndServe("localhost:8081", logRequest(fs)))
}
