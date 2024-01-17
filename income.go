package main

import (
	"encoding/csv"
	"io"
	"log"
	"os/exec"
	"strconv"
	"strings"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
)

func (c Charts) incomeFromInvestmentsPieChart(date string) *charts.Pie {
	hlopts := c.hlopts.
		WithAccount("income:investment").
		WithAccountDrop(1).
		WithStartDate(date).
		WithInvertAmount()

	rd, err := c.hl.Balance(hlopts)
	if err != nil {
		log.Fatal(err)
	}
	out, _ := io.ReadAll(rd)
	data := string(out)
	pie := charts.NewPie()
	pie.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "Income from investments  " + date,
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

func (c Charts) incomePieChart(date string) *charts.Pie {
	hlopts := c.hlopts.
		WithAccount("income").
		WithAccountDrop(1).
		WithStartDate(date).
		WithInvertAmount()

	rd, err := c.hl.Balance(hlopts)
	if err != nil {
		log.Fatal(err)
	}
	out, _ := io.ReadAll(rd)
	data := string(out)

	pie := charts.NewPie()
	pie.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "Income  " + date,
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

func incomeStatementBarChartMonthly(date string) *charts.Bar {
	incomeStmt, err := exec.Command("hledger", "incomestatement", "-M", "--depth", "1", "-p", date, "--layout", "bare", "-O", "csv").Output()
	if err != nil {
		log.Fatal(err)
	}
	xdata, ydata := parseCSVIncomeStatement(string(incomeStmt))

	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "Monthly Income & Expense " + date,
		}),
		charts.WithInitializationOpts(defaultTheme()),
		charts.WithLegendOpts(opts.Legend{Show: true}),
		charts.WithTooltipOpts(opts.Tooltip{Show: true}),
	)
	bar.SetXAxis(xdata).
		AddSeries("income", ydata[0]).
		AddSeries("expense", ydata[1]).
		SetSeriesOptions(
			charts.WithLabelOpts(opts.Label{
				Show:     true,
				Position: "top",
			}),
		)
	return bar
}

// parseCSV2: only 1 row of data. Columnname is X, row is Y
func parseCSV2(data string) ([]string, []opts.BarData) {
	r := csv.NewReader(strings.NewReader(data))

	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	var xdata []string
	var ydata []opts.BarData

	xdata = records[0][1:]

	for _, record := range records[1][1:] {
		num := strings.Replace(record, "SGD$", "", 1)
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

// parseCSVIncomeStatement: same as parseCSV2 but returns specific rows
func parseCSVIncomeStatement(data string) ([]string, [][]opts.BarData) {
	r := csv.NewReader(strings.NewReader(data))

	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	var xdata []string

	xdata = records[1][2:]

	parseYData := func(rowNum int, records [][]string) []opts.BarData {
		var ydata []opts.BarData
		for _, record := range records[rowNum][2:] {
			num := strings.Replace(record, ",", "", 1)
			amt, err := strconv.ParseFloat(num, 64)
			if err != nil {
				amt = 0
				log.Println(err)
			}
			ydata = append(ydata, opts.BarData{Value: amt})
		}
		return ydata
	}

	incomeData := parseYData(3, records)
	expenseData := parseYData(6, records)

	return xdata, [][]opts.BarData{incomeData, expenseData}
}
