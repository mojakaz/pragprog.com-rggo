package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
)

// statsFunc defines a generic statistical function
type statsFunc func(data []float64) float64

func sum(data []float64) float64 {
	sum := 0.0
	for _, v := range data {
		sum += v
	}
	return sum
}

func avg(data []float64) float64 {
	return sum(data) / float64(len(data))
}

func Min(data []float64) float64 {
	if len(data) == 0 {
		return 0.0
	}
	current := data[0]
	for _, num := range data {
		current = min(current, num)
	}
	return current
}

func Max(data []float64) float64 {
	if len(data) == 0 {
		return 0.0
	}
	current := data[0]
	for _, num := range data {
		current = max(current, num)
	}
	return current
}

func csv2float(r io.Reader, column int) ([]float64, error) {
	var data []float64

	// Create the CSV Reader used to read in data from CSV files
	cr := csv.NewReader(r)
	cr.ReuseRecord = true
	// Adjusting for 0 based index
	column--
	// Looping through all records
	for i := 0; ; i++ {
		row, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("cannot read data from file: %w", err)
		}
		if i == 0 {
			continue
		}
		// Checking number of columns in CSV file
		if len(row) <= column {
			// File does not have that many columns
			return nil, fmt.Errorf("%w: file has only %d columns", ErrInvalidColumn, len(row))
		}
		// Try to convert data read into a float number
		v, err := strconv.ParseFloat(row[column], 64)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", ErrNotNumber, err)
		}
		data = append(data, v)
	}
	return data, nil
}
