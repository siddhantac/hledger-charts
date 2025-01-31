package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/siddhantac/hledger"
)

func (c Charts) incomeFromInvestmentsPieChart(date string) *charts.Pie {
	hlopts := c.hlopts.
		WithAccount("income:investment").
		WithAccountDrop(1).
		WithStartDate(date).
		WithInvertAmount(true)

	rd, err := c.hl.Balance(hlopts)
	if err != nil {
		log.Fatal(err)
	}
	out, _ := io.ReadAll(rd)
	data := string(out)
	pie := charts.NewPie()
	pie.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "Income from investments",
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

func (c Charts) incomePieChart(date string) *charts.Pie {
	hlopts := c.hlopts.
		WithAccount("income").
		WithAccountDrop(1).
		WithAccountDepth(2).
		WithStartDate(date).
		WithInvertAmount(true)

	rd, err := c.hl.Balance(hlopts)
	if err != nil {
		log.Fatal(err)
	}
	out, _ := io.ReadAll(rd)
	data := string(out)

	pie := charts.NewPie()
	pie.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "Income distribution",
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

func (c Charts) incomeStatementBarChartMonthly(year int) *charts.Bar {
	hlopts := c.hlopts.
		WithAccountDepth(1).
		WithStartDate(fmt.Sprintf("%d", year)).
		WithEndDate(fmt.Sprintf("%d", year+1)).
		WithPeriod(hledger.PeriodMonthly)

	rd, err := c.hl.IncomeStatement(hlopts)
	if err != nil {
		log.Fatal(err)
	}
	out, _ := io.ReadAll(rd)

	// out, err := exec.Command("hledger", "incomestatement", "-M", "--depth", "1", "-p", date, "--layout", "bare", "-O", "csv").Output()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	xdata, ydata := parseCSVIncomeStatement(string(out))

	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "Monthly Income & Expense",
		}),
		charts.WithInitializationOpts(defaultTheme()),
		charts.WithLegendOpts(opts.Legend{Show: opts.Bool(true)}),
		charts.WithTooltipOpts(opts.Tooltip{Show: opts.Bool(true)}),
	)
	bar.SetXAxis(xdata).
		AddSeries("expense", ydata[1]).
		AddSeries("income", ydata[0]).
		SetSeriesOptions(
			charts.WithLabelOpts(opts.Label{
				Show:     opts.Bool(true),
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
