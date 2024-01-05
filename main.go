package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
)

func defaultTheme() opts.Initialization {
	return opts.Initialization{
		Theme:           types.ThemeRoma,
		BackgroundColor: "white",
	}
}

func createMonthlyReport(date string) {
	page := components.NewPage()
	page.Layout = components.PageCenterLayout
	page.AddCharts(
		expensesPieChart(date),
		expensesHorizontalBarChart(date),
		investmentsPieChart(date),
	)

	filename := fmt.Sprintf("reports/%s.html", date)
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	page.Render(f)
	log.Println("generated:", filename)
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func generateMonthlyReports(year string) {
	for i := 1; i <= 12; i++ {
		month := fmt.Sprintf("%s-%0.2d", year, i)
		createMonthlyReport(month)
	}
}

func main() {
	var (
		generateReports bool
		outputDir       string
	)

	flag.BoolVar(&generateReports, "generate", false, "generate reports")
	flag.StringVar(&outputDir, "output", "", "output directory")
	flag.Parse()

	if generateReports {
		c := NewCharts(outputDir)
		c.createYearlyReport("2023")
		// createMonthlyReport("last month")
		// generateMonthlyReports("2023")
	}

	fs := http.FileServer(http.Dir(outputDir))
	log.Println("running server at http://localhost:8081")
	log.Fatal(http.ListenAndServe("localhost:8081", logRequest(fs)))
}
