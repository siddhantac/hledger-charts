package main

import (
	"log"
	"os/exec"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
)

func expensesPieChart(date string) *charts.Pie {
	out, err := exec.Command("hledger", "bal", "expenses", "--drop", "1", "--depth", "2", "-p", date, "--layout", "bare", "-O", "csv").Output()
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

func reverse[T any](arr []T) []T {
	reversed := make([]T, len(arr))
	for i := range arr {
		reversed[i] = arr[len(arr)-i-1]
	}
	return reversed
}
