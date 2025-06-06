package main

import (
	"io"
	"log"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func (c Charts) expensesPieChart(startDate, endDate string) *charts.Pie {
	hlopts := c.hlopts.
		WithAccount("expenses").
		WithAccountDrop(1).
		WithAccountDepth(2).
		WithStartDate(startDate).
		WithEndDate(endDate)

	rd, err := c.hl.Balance(hlopts)
	if err != nil {
		log.Fatal(err)
	}

	out, _ := io.ReadAll(rd)

	data := string(out)
	pie := charts.NewPie()
	pie.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "Expense distribution",
		}),
		charts.WithInitializationOpts(defaultTheme()),
		charts.WithLegendOpts(opts.Legend{Show: opts.Bool(false)}),
		charts.WithTooltipOpts(opts.Tooltip{Show: opts.Bool(true)}),
		// alternative option
		// charts.WithTooltipOpts(opts.Tooltip{
		// 	Show:      opts.Bool(true),
		// 	Trigger:   "item",
		// 	Formatter: "{b}: ${c} ({d}%)",
		// }),
	)

	pie.AddSeries("pie", parseCSV1a(data)).
		SetSeriesOptions(
			charts.WithLabelOpts(opts.Label{
				Show:      opts.Bool(true),
				Formatter: "{b}: {d}%",
			}),
			charts.WithPieChartOpts(opts.PieChart{
				Radius: []string{"40%", "75%"},
			}),
		)
	return pie
}

func (c Charts) expensesHorizontalBarChart(startDate, endDate string) *charts.Bar {
	hlopts := c.hlopts.
		WithAccount("expenses").
		WithAccountDrop(1).
		WithAccountDepth(2).
		WithSortAmount(true).
		WithStartDate(startDate).
		WithEndDate(endDate)
		// TODO: add -S

	rd, err := c.hl.Balance(hlopts)
	if err != nil {
		log.Fatal(err)
	}
	out, _ := io.ReadAll(rd)
	data := string(out)
	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "Expenses Top List",
		}),
		charts.WithInitializationOpts(defaultTheme()),
		charts.WithLegendOpts(opts.Legend{Show: opts.Bool(false)}),
		charts.WithTooltipOpts(opts.Tooltip{Show: opts.Bool(true)}),
	)

	xdata, ydata := parseCSV1(data)
	xdata = reverse(xdata)
	ydata = reverse(ydata)

	bar.SetXAxis(xdata).
		AddSeries("expenses", ydata).
		SetSeriesOptions(
			charts.WithLabelOpts(opts.Label{
				Show:     opts.Bool(true),
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
