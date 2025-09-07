package main

import (
	"io"
	"log"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func (c Charts) budgetBarChart(startDate, endDate string) *charts.Bar {
	hlopts := c.hlopts.
		WithAccount("expenses").
		WithAccountDrop(1).
		WithAccountDepth(2).
		WithSortAmount(true).
		WithStartDate(startDate).
		WithEndDate(endDate).
		WithBudget()

	rd, err := c.hl.Balance(hlopts)
	if err != nil {
		log.Fatal(err)
	}
	out, _ := io.ReadAll(rd)
	data := string(out)
	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "Budget",
		}),
		charts.WithInitializationOpts(defaultTheme()),
		charts.WithLegendOpts(opts.Legend{Show: opts.Bool(false)}),
		charts.WithTooltipOpts(opts.Tooltip{Show: opts.Bool(true)}),
	)

	xdata, ydata := parseBudgetData(data)
	xdata = reverse(xdata)
	ydata = reverse(ydata)

	bar.SetXAxis(xdata).
		AddSeries("expenses", ydata).
		SetSeriesOptions(
			charts.WithLabelOpts(opts.Label{
				Show:     opts.Bool(true),
				Position: "right",
			}),
			charts.WithMarkLineNameXAxisItemOpts(
				opts.MarkLineNameXAxisItem{
					Name:  "max",
					XAxis: 100,
				},
			),
		)
	bar.XYReversal()
	return bar
}
