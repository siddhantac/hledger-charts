package main

import (
	"io"
	"log"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func (c Charts) investmentsPieChart(date, endDate string) *charts.Pie {
	hlopts := c.hlopts.
		WithAccount("assets:investment").
		WithAccountDrop(2).
		WithAccountDepth(3).
		WithStartDate(date).
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
			Title: "Investments",
		}),
		charts.WithTooltipOpts(opts.Tooltip{Show: opts.Bool(true)}),
		charts.WithInitializationOpts(defaultTheme()),
		charts.WithLegendOpts(opts.Legend{Show: opts.Bool(false)}),
	)

	pie.AddSeries("pie", parseCSV1a(data)).
		SetSeriesOptions(
			charts.WithLabelOpts(opts.Label{
				Show:      opts.Bool(true),
				Formatter: "{b}: {c}",
			}),
			charts.WithPieChartOpts(opts.PieChart{
				Radius: []string{"40%", "75%"},
			}),
		)
	return pie
}
