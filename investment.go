package main

import (
	"log"
	"os/exec"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
)

func investmentsPieChart(date string) *charts.Pie {
	out, err := exec.Command("hledger", "bal", "assets:invest", "--drop", "2", "-p", date, "-O", "csv").Output()
	if err != nil {
		log.Fatal(err)
	}
	data := string(out)
	pie := charts.NewPie()
	pie.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "Investments  " + date,
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
