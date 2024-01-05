package main

import (
	"flag"
	"log"
	"net/http"

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
	)

	flag.BoolVar(&generateReports, "generate", false, "generate reports")
	flag.StringVar(&outputDir, "output", "", "output directory")
	flag.Parse()

	if outputDir == "" {
		log.Fatal("output directory must be specified")
	}

	if generateReports {
		c := NewCharts(outputDir)
		c.createYearlyReport("2023")
		c.createMonthlyReport("last month")
		c.generateMonthlyReports("2023")
	}

	fs := http.FileServer(http.Dir(outputDir))
	log.Println("running server at http://localhost:8081")
	log.Fatal(http.ListenAndServe("localhost:8081", logRequest(fs)))
}
