package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/siddhantac/hledger"
)

type Charts struct {
	OutputDir string
	Year      int
	hl        hledger.Hledger
	hlopts    hledger.Options
}

func NewCharts(outputDir string, year int) Charts {
	dir := filepath.Join(outputDir, fmt.Sprintf("%d", year))
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	hlopts := hledger.NewOptions().
		WithLayout(hledger.LayoutBare).
		WithOutputCSV(true)

	return Charts{
		OutputDir: dir,
		Year:      year,
		hl:        hledger.New("hledger", ""),
		hlopts:    hlopts,
	}
}

func (c Charts) createYearlyReport() {
	date := fmt.Sprintf("%d", c.Year)

	page := components.NewPage()
	page.Layout = components.PageFlexLayout
	page.AddCharts(
		c.incomeStatementBarChartMonthly(c.Year),
		c.expensesHorizontalBarChart(date),
		c.expensesPieChart(date),
		c.incomePieChart(date),
		c.investmentsPieChart(date),
		c.incomeFromInvestmentsPieChart(date),
	)
	page.PageTitle = fmt.Sprintf("%d Yearly Report", c.Year)

	filename := fmt.Sprintf("%s/%s.html", c.OutputDir, date)
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	page.Render(f)
	log.Println("generated:", filename)
}

func (c Charts) createMonthlyReport(t time.Time) {
	date := fmt.Sprintf("%d-%0.2d", t.Year(), t.Month())
	page := components.NewPage()
	page.Layout = components.PageFlexLayout
	page.AddCharts(
		c.expensesPieChart(date),
		c.expensesHorizontalBarChart(date),
		c.investmentsPieChart(date),
		c.incomePieChart(date),
	)

	filename := fmt.Sprintf("%s/%s.html", c.OutputDir, date)
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	page.Render(f)
	log.Println("generated:", filename)
}

func (c Charts) generateMonthlyReports(year int) {
	var wg sync.WaitGroup
	workCh := make(chan time.Time, 12)

	endMonth := 12
	now := time.Now()
	if year == now.Year() {
		endMonth = int(now.Month()) - 1
	}
	for i := time.January; i <= time.Month(endMonth); i++ {
		t := time.Date(year, time.Month(i), 1, 0, 0, 0, 0, time.UTC)
		workCh <- t
	}
	close(workCh)

	numWorkers := 3
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			for date := range workCh {
				c.createMonthlyReport(date)
			}
		}(&wg)
	}
	wg.Wait()
}
