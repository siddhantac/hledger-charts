package main

import (
	"fmt"
	"log"
	"os"

	"github.com/go-echarts/go-echarts/v2/components"
)

type Charts struct {
	OutputDir string
}

func NewCharts(outputDir string) Charts {
	return Charts{OutputDir: outputDir}
}

func (c Charts) createYearlyReport(date string) {
	page := components.NewPage()
	page.Layout = components.PageCenterLayout
	page.AddCharts(
		incomeStatementBarChartMonthly(date),
		expensesPieChart(date),
		expensesHorizontalBarChart(date),
		investmentsPieChart(date),
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
	for i := 1; i <= 12; i++ {
		month := fmt.Sprintf("%s-%0.2d", year, i)
		c.createMonthlyReport(month)
	}
}
