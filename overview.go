package main

import (
	"encoding/csv"
	"io"
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/siddhantac/hledger"
)

func (c Charts) yoyChart(date, endDate string) *charts.Line {
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "Overview - YoY", Subtitle: "Income, investments, expenses"}),
	)

	hlopts := c.hlopts.
		WithAccount("income|assets:investment|expenses").
		WithAccountDepth(1).
		WithStartDate(date).
		WithEndDate(endDate).
		WithPeriod(hledger.PeriodYearly).
		WithValuation(true)

	rd, err := c.hl.Balance(hlopts)
	if err != nil {
		log.Fatal(err)
	}
	out, _ := io.ReadAll(rd)
	data := string(out)
	r := csv.NewReader(strings.NewReader(data))

	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	toLineData := func(data []string) []opts.LineData {
		var lineData []opts.LineData
		for _, d := range data {
			num := strings.Replace(d, ",", "", 1)
			amt, err := strconv.ParseFloat(num, 64)
			if err != nil {
				amt = 0
				log.Println(err)
			}
			amt = math.Abs(amt)
			lineData = append(lineData, opts.LineData{Value: amt})
		}
		return lineData
	}

	investmentData := toLineData(records[1][2:])
	expenseData := toLineData(records[2][2:])
	incomeData := toLineData(records[3][2:])

	years := records[0][2:]
	line.SetXAxis(years).
		AddSeries("investment", investmentData).
		AddSeries("expense", expenseData).
		AddSeries("income", incomeData)
	return line
}
