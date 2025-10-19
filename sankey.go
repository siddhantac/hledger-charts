package main

import (
	"encoding/csv"
	"log"
	"math"
	"strconv"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/siddhantac/hledger"
)

type Data struct {
	Nodes []opts.SankeyNode
	Links []opts.SankeyLink
}

func (c Charts) sankeyDataForSubAccount(hlopts hledger.Options, account, replaceTotal string) Data {
	rd, err := c.hl.Balance(hlopts)
	if err != nil {
		log.Fatal(err)
	}
	reader := csv.NewReader(rd)

	// Read and skip header row
	_, err = reader.Read()
	if err != nil {
		log.Fatal(err)
	}

	// Parse CSV data
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	data := Data{
		Nodes: []opts.SankeyNode{},
		Links: []opts.SankeyLink{},
	}

	// Parse each record
	for _, record := range records {
		if len(record) >= 3 {
			subAccount := record[0]
			balanceStr := record[2]

			if subAccount == account {
				continue
			}

			if subAccount == "Total:" {
				if replaceTotal == "" {
					continue
				}
				subAccount = replaceTotal
			}

			if subAccount == "assets" {
				subAccount = "assets:investment"
			}

			amt, err := strconv.ParseFloat(balanceStr, 64)
			if err != nil {
				amt = 0
				log.Println(err)
			}

			data.Nodes = append(data.Nodes, opts.SankeyNode{Name: subAccount})
			data.Links = append(data.Links, opts.SankeyLink{
				Source: account,
				Target: subAccount,
				Value:  float32(math.Abs(amt)),
			})
		}
	}

	return data
}

func mergeData(a Data, all ...Data) Data {
	for _, b := range all {
		a.Nodes = append(a.Nodes, b.Nodes...)
		a.Links = append(a.Links, b.Links...)
	}
	return a
}

func (c Charts) sankeyChart(date, endDate string) *charts.Sankey {
	hlopts := c.hlopts.
		WithAccount("expenses").
		WithAccountDepth(2).
		WithAccountDrop(1).
		WithStartDate(date).
		WithEndDate(endDate).
		WithValuation(true)
	expenseData := c.sankeyDataForSubAccount(hlopts, "expenses", "")

	hlopts = c.hlopts.
		WithAccount("assets:investment").
		WithAccountDepth(3).
		WithAccountDrop(1).
		WithStartDate(date).
		WithEndDate(endDate).
		WithValuation(true)
	investmentData := c.sankeyDataForSubAccount(hlopts, "assets:investment", "")

	hlopts = c.hlopts.
		WithAccount("income|expenses|assets:investment").
		WithAccountDepth(1).
		WithAccountDrop(0).
		WithStartDate(date).
		WithEndDate(endDate).
		WithValuation(true)
	incomeData := c.sankeyDataForSubAccount(hlopts, "income", "savings")
	incomeData.Nodes = append(incomeData.Nodes, opts.SankeyNode{Name: "income"})

	data := mergeData(investmentData, expenseData, incomeData)

	sankey := charts.NewSankey()
	sankey.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "Where is our money going",
		}),
		charts.WithTooltipOpts(opts.Tooltip{Show: opts.Bool(true)}),
		charts.WithInitializationOpts(defaultTheme()),
	)

	sankey.AddSeries("sankey", data.Nodes, data.Links).
		SetSeriesOptions(
			charts.WithLineStyleOpts(opts.LineStyle{
				// Color:     "source",
				Curveness: 0.3,
			}),
			charts.WithLabelOpts(opts.Label{
				Show:      opts.Bool(true),
				Formatter: "{b}: ${c}",
			}),
		)
	return sankey
}
