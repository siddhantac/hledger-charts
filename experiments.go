package main

import (
	"encoding/csv"
	"log"
	"os/exec"
	"strconv"
	"strings"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func subExpenses(subcategory string) *charts.Bar {
	out, err := exec.Command("hledger", "bal", "expenses:"+subcategory, "-M", "-p", "2023", "-O", "csv").Output()
	if err != nil {
		log.Fatal(err)
	}
	data := string(out)
	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    "Expenses 2023",
			Subtitle: "Dine out",
		}),
		charts.WithInitializationOpts(defaultTheme()),
	)

	xdata, ydata := parseCSV2(data)

	bar.SetXAxis(xdata).
		AddSeries(subcategory, ydata).
		SetSeriesOptions(
			charts.WithLabelOpts(opts.Label{
				Show:     opts.Bool(true),
				Position: "top",
			}),
		)
	return bar
}

func subExpenses2(subcategory string) *charts.Line {
	out, err := exec.Command("hledger", "bal", "expenses:"+subcategory, "-M", "-p", "2023", "-O", "csv").Output()
	if err != nil {
		log.Fatal(err)
	}
	data := string(out)
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    "Expenses 2023",
			Subtitle: "Dine out",
		}),
		charts.WithTooltipOpts(opts.Tooltip{Show: opts.Bool(true)}),
		charts.WithInitializationOpts(defaultTheme()),
	)

	xdata, ydata := parseCSV3(data)

	line.SetXAxis(xdata).
		AddSeries(subcategory, ydata).
		SetSeriesOptions(
			charts.WithLabelOpts(opts.Label{
				Show: opts.Bool(true),
			}),
			charts.WithLineChartOpts(opts.LineChart{
				ShowSymbol: opts.Bool(true),
			}),
		)
	return line
}

// parseCSV3: same as parseCSV2, but returns line data instead
func parseCSV3(data string) ([]string, []opts.LineData) {
	r := csv.NewReader(strings.NewReader(data))

	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	var xdata []string
	var ydata []opts.LineData

	xdata = records[0][1:]

	for _, record := range records[1][1:] {
		num := strings.Replace(record, "SGD$", "", 1)
		num = strings.Replace(num, ",", "", 1)
		amt, err := strconv.ParseFloat(num, 64)
		if err != nil {
			amt = 0
			log.Println(err)
		}
		ydata = append(ydata, opts.LineData{Value: amt})
	}

	return xdata, ydata
}
