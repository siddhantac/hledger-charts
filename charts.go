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
		WithOutputCSV()

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
	page.Layout = components.PageCenterLayout
	page.AddCharts(
		c.incomeStatementBarChartMonthly(date),
		c.expensesPieChart(date),
		c.expensesHorizontalBarChart(date),
		c.investmentsPieChart(date),
		c.incomePieChart(date),
		c.incomeFromInvestmentsPieChart(date),
	)
	page.PageTitle = "Yearly Report 2023"

	filename := fmt.Sprintf("%s/%s.html", c.OutputDir, date)
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	page.Render(f)
	log.Println("generated:", filename)
}

func (c Charts) createMonthlyReport(date string) {
	page := components.NewPage()
	page.Layout = components.PageCenterLayout
	page.AddCharts(
		c.expensesPieChart(date),
		c.expensesHorizontalBarChart(date),
		c.investmentsPieChart(date),
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
	workCh := make(chan string, 12)

	endMonth := 12
	now := time.Now()
	if year == now.Year() {
		endMonth = int(now.Month())
	}
	for i := 1; i < endMonth; i++ {
		month := fmt.Sprintf("%d-%0.2d", year, i)
		workCh <- month
	}
	close(workCh)

	numWorkers := 3
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			for month := range workCh {
				c.createMonthlyReport(month)
			}
		}(&wg)
	}
	wg.Wait()
}
