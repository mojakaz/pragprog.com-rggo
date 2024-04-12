package main

import (
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	// Test cases for Run Tests
	testCases := []struct {
		name   string
		col    int
		op     string
		exp    string
		files  []string
		expErr error
	}{
		{
			name:   "RunAvg1File",
			col:    3,
			op:     "avg",
			exp:    "227.6\n",
			files:  []string{"./testdata/example.csv"},
			expErr: nil,
		},
		{
			name:   "RunAvgMultiFiles",
			col:    3,
			op:     "avg",
			exp:    "233.84\n",
			files:  []string{"./testdata/example.csv", "./testdata/example2.csv"},
			expErr: nil,
		},
		{
			name:   "RunFailRead",
			col:    2,
			op:     "avg",
			exp:    "",
			files:  []string{"./testdata/example.csv", "./testdata/fakefile.csv"},
			expErr: os.ErrNotExist,
		},
		{
			name:   "RunFailColumn",
			col:    0,
			op:     "avg",
			exp:    "",
			files:  []string{"./testdata/example.csv"},
			expErr: ErrInvalidColumn,
		},
		{
			name:   "RunFailNoFiles",
			col:    2,
			op:     "avg",
			exp:    "",
			files:  []string{},
			expErr: ErrNoFiles,
		},
		{
			name:   "RunFailOperation",
			col:    2,
			op:     "invalid",
			exp:    "",
			files:  []string{"./testdata/example.csv"},
			expErr: ErrInvalidOperation,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var res strings.Builder
			err := run(tc.files, tc.op, tc.col, &res)
			if tc.expErr != nil {
				assert.Errorf(t, err, "expected error, got nil instead")
				assert.ErrorIsf(t, err, tc.expErr, "expected error %q matches %q", err, tc.expErr)
				return
			}
			assert.NoErrorf(t, err, "expected no error, got %q instead", err)
			assert.Equalf(t, tc.exp, res.String(), "expected %q, got %q instead", tc.exp, &res)
		})
	}
}

func BenchmarkRunAvg(b *testing.B) {
	filenames, err := filepath.Glob("./testdata/benchmark/*.csv")
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := run(filenames, "avg", 2, io.Discard); err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkRunSum(b *testing.B) {
	filenames, err := filepath.Glob("./testdata/benchmark/*.csv")
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := run(filenames, "sum", 2, io.Discard); err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkRunMin(b *testing.B) {
	filenames, err := filepath.Glob("./testdata/benchmark/*.csv")
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := run(filenames, "Min", 2, io.Discard); err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkRunMax(b *testing.B) {
	filenames, err := filepath.Glob("./testdata/benchmark/*.csv")
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := run(filenames, "Max", 2, io.Discard); err != nil {
			b.Error(err)
		}
	}
}
