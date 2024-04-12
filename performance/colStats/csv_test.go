package main

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
	"testing/iotest"
)

func TestOperations(t *testing.T) {
	data := [][]float64{
		{10, 20, 15, 30, 45, 50, 100, 30},
		{5.5, 8, 2.2, 9.75, 8.45, 3, 2.5, 10.25, 4.75, 6.1, 7.67, 12.287, 5.47},
		{-10, -20},
		{102, 37, 44, 57, 67, 129},
	}
	// Test cases for Operations Test
	testCases := []struct {
		name string
		op   statsFunc
		exp  []float64
	}{
		{"Sum", sum, []float64{300, 85.927, -30, 436}},
		{"Avg", avg, []float64{37.5, 6.609769230769231, -15, 72.666666666666666}},
		{name: "Min", op: Min, exp: []float64{10, 2.2, -20, 37}},
		{name: "Max", op: Max, exp: []float64{100, 12.287, -10, 129}},
	}
	// Operations Tests execution
	for _, tc := range testCases {
		for k, exp := range tc.exp {
			name := fmt.Sprintf("%sData%d", tc.name, k)
			t.Run(name, func(t *testing.T) {
				res := tc.op(data[k])
				assert.Equalf(t, exp, res, "Expected %g, got %g instead", exp, res)
			})
		}
	}
}

func TestCSV2Float(t *testing.T) {
	csvData := `IP Address,Requests,Response Time
192.168.0.199,2056,236
192.168.0.88,899,220
192.168.0.199,3054,226
192.168.0.100,4133,218
192.168.0.199,950,238
`
	// Test cases for CSV2Float Test
	testCases := []struct {
		name   string
		col    int
		exp    []float64
		expErr error
		r      io.Reader
	}{
		{
			name:   "Column2",
			col:    2,
			exp:    []float64{2056, 899, 3054, 4133, 950},
			expErr: nil,
			r:      bytes.NewBufferString(csvData),
		},
		{
			name:   "Column3",
			col:    3,
			exp:    []float64{236, 220, 226, 218, 238},
			expErr: nil,
			r:      bytes.NewBufferString(csvData),
		},
		{
			name:   "FailRead",
			col:    1,
			exp:    nil,
			expErr: iotest.ErrTimeout,
			r:      iotest.TimeoutReader(bytes.NewReader([]byte{0})),
		},
		{
			name:   "FailedNotNumber",
			col:    1,
			exp:    nil,
			expErr: ErrNotNumber,
			r:      bytes.NewBufferString(csvData),
		},
		{
			name:   "FailedInvalidColumn",
			col:    4,
			exp:    nil,
			expErr: ErrInvalidColumn,
			r:      bytes.NewBufferString(csvData),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := csv2float(tc.r, tc.col)
			if tc.expErr != nil {
				assert.Errorf(t, err, "expected error, got nil instead")
				assert.ErrorIsf(t, err, tc.expErr, "expected error %q matches %q", err, tc.expErr)
				return
			}
			assert.NoErrorf(t, err, "expected no error, got %q instead", err)

			assert.Equalf(t, tc.exp, res, "expected %g, got %g instead", tc.exp, res)
		})
	}
}
