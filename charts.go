package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/go-echarts/go-echarts/v2/components"
)

type Charts struct {
	OutputDir string
	Year      int
}

func NewCharts(outputDir string, year int) Charts {
	dir := filepath.Join(outputDir, fmt.Sprintf("%d", year))
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		log.Fatal(err)
	}
	return Charts{
		OutputDir: dir,
		Year:      year,
	}
}

func (c Charts) createYearlyReport() {
	date := fmt.Sprintf("%d", c.Year)

	page := components.NewPage()
	page.Layout = components.PageCenterLayout
	page.AddCharts(
		incomeStatementBarChartMonthly(date),
		expensesPieChart(date),
		expensesHorizontalBarChart(date),
		investmentsPieChart(date),
		incomePieChart(date),
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
		expensesPieChart(date),
		expensesHorizontalBarChart(date),
		investmentsPieChart(date),
	)

	filename := fmt.Sprintf("%s/%s.html", c.OutputDir, date)
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	page.Render(f)
	log.Println("generated:", filename)
}

func (c Charts) generateMonthlyReports(year string) {
	workCh := make(chan string, 12)
	for i := 1; i <= 12; i++ {
		month := fmt.Sprintf("%s-%0.2d", year, i)
		workCh <- month
	}
	close(workCh)

	numWorkers := 3
	for i := 0; i < numWorkers; i++ {
		go func() {
			for month := range workCh {
				c.createMonthlyReport(month)
			}
		}()
	}
}
