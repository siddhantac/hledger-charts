package main

import (
	"encoding/csv"
	"log"
	"os/exec"
	"strconv"
	"strings"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
)

func expensesPieChart(date string) *charts.Pie {
	out, err := exec.Command("hledger", "bal", "expenses", "--drop", "1", "--depth", "2", "-p", date, "-O", "csv").Output()
	if err != nil {
		log.Fatal(err)
	}
	data := string(out)
	pie := charts.NewPie()
	pie.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "Expenses pie chart " + date,
		}),
		charts.WithTooltipOpts(opts.Tooltip{Show: true}),
		charts.WithInitializationOpts(
			opts.Initialization{
				Theme:           types.ThemeRoma,
				BackgroundColor: "white",
			}),
		charts.WithLegendOpts(opts.Legend{Show: false}),
	)

	pie.AddSeries("pie", parseCSV1a(data)).
		SetSeriesOptions(
			charts.WithLabelOpts(opts.Label{
				Show:      true,
				Formatter: "{b}: {c}",
			}),
			charts.WithPieChartOpts(opts.PieChart{
				Radius: []string{"40%", "75%"},
			}),
		)
	return pie
}

func expensesHorizontalBarChart(date string) *charts.Bar {
	out, err := exec.Command("hledger", "bal", "expenses", "--depth", "2", "-S", "-p", date, "-O", "csv").Output()
	if err != nil {
		log.Fatal(err)
	}
	data := string(out)
	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "Expenses Top list " + date,
		}),
		charts.WithInitializationOpts(defaultTheme()),
		charts.WithLegendOpts(opts.Legend{Show: false}),
		charts.WithTooltipOpts(opts.Tooltip{Show: true}),
	)

	xdata, ydata := parseCSV1(data)
	xdata = reverse(xdata)
	ydata = reverse(ydata)

	bar.SetXAxis(xdata).
		AddSeries("expenses", ydata).
		SetSeriesOptions(
			charts.WithLabelOpts(opts.Label{
				Show:     true,
				Position: "right",
			}),
		)
	bar.XYReversal()
	return bar
}

// parseCSV1a : same as parseCSV1 but returns pie chart data
func parseCSV1a(data string) []opts.PieData {
	r := csv.NewReader(strings.NewReader(data))

	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	var pieData []opts.PieData

	for _, record := range records[1 : len(records)-1] {
		name := record[0]
		num := strings.Replace(record[1], "SGD$", "", 1)
		num = strings.Replace(num, ",", "", 1)
		amt, err := strconv.ParseFloat(num, 64)
		if err != nil {
			amt = 0
			log.Println(err)
		}
		pieData = append(pieData, opts.PieData{Name: name, Value: amt})
	}

	return pieData
}

// parseCSV1 : each row is a X-Y value pair
func parseCSV1(data string) ([]string, []opts.BarData) {
	r := csv.NewReader(strings.NewReader(data))

	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	var xdata []string
	var ydata []opts.BarData

	for _, record := range records[1 : len(records)-1] {
		xdata = append(xdata, strings.Replace(record[0], "expenses:", "", 1))

		num := strings.Replace(record[1], "SGD$", "", 1)
		num = strings.Replace(num, ",", "", 1)
		amt, err := strconv.ParseFloat(num, 64)
		if err != nil {
			amt = 0
			log.Println(err)
		}
		ydata = append(ydata, opts.BarData{Value: amt})
	}

	return xdata, ydata
}

func reverse[T any](arr []T) []T {
	reversed := make([]T, len(arr))
	for i := range arr {
		reversed[i] = arr[len(arr)-i-1]
	}
	return reversed
}
