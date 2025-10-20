package main

import (
	"encoding/csv"
	"log"
	"strconv"
	"strings"

	"github.com/go-echarts/go-echarts/v2/opts"
)

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
		num := record[2]
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
		xdata = append(xdata, record[0])
		num := record[2]
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

func parseBudgetData(data string) ([]string, []opts.BarData) {
	r := csv.NewReader(strings.NewReader(data))

	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	var xdata []string
	var ydata []opts.BarData

	for _, record := range records[1 : len(records)-1] {
		xdata = append(xdata, record[0])
		value := getAmount(record[2])
		budget := getAmount(record[3])

		var consumedBudget int
		if budget == 0 {
			consumedBudget = 0
		} else {
			consumedBudget = int((value / budget) * 100)
		}
		ydata = append(ydata, opts.BarData{Value: consumedBudget})
	}

	return xdata, ydata
}

func getAmount(num string) float64 {
	num = strings.Replace(num, ",", "", 1)
	amt, err := strconv.ParseFloat(num, 64)
	if err != nil {
		amt = 0
		log.Println(err)
	}
	return amt
}
