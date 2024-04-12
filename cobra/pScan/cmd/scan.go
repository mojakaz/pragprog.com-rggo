/*
Copyright © 2024 Kazuki Takemoto

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"os"
	"pragprog.com/rggo/cobra/pScan/scan"
	"strconv"
	"strings"
	"time"
)

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Run a port scan on the hosts",
	RunE: func(cmd *cobra.Command, args []string) error {
		hostsFile := viper.GetString("hosts-file")
		ports, err := cmd.Flags().GetIntSlice("ports")
		if err != nil {
			return err
		}
		pr, err := cmd.Flags().GetString("port-range")
		if err != nil {
			return err
		}
		useUDP, err := cmd.Flags().GetBool("use-udp")
		if err != nil {
			return err
		}
		timeout, err := cmd.Flags().GetDuration("timeout")
		if err != nil {
			return err
		}
		return scanAction(os.Stdout, hostsFile, ports, pr, useUDP, timeout)
	},
}

func scanAction(out io.Writer, hostsFile string, ports []int, pr string, useUDP bool, timeout time.Duration) error {
	hl := &scan.HostsList{}
	if err := hl.Load(hostsFile); err != nil {
		return err
	}
	if pr != "" {
		r := strings.Split(pr, "-")
		var (
			start, end int
		)
		start, err := strconv.Atoi(r[0])
		if err != nil {
			return err
		}
		end, err = strconv.Atoi(r[1])
		if err != nil {
			return err
		}
		if start < 0 || end > 65543 || start > end {
			return fmt.Errorf("invaild port range: start %d, end %d", start, end)
		}
		for i := start; i < end+1; i++ {
			ports = append(ports, i)
		}
	}
	results := scan.Run(hl, ports, useUDP, timeout)
	return printResults(out, results)
}

func printResults(out io.Writer, results []scan.Results) error {
	message := ""
	for _, r := range results {
		message += fmt.Sprintf("%s:", r.Host)
		if r.NotFound {
			message += fmt.Sprintf(" Host not found\n")
			continue
		}
		message += fmt.Sprintln()
		for _, p := range r.PortStates {
			message += fmt.Sprintf("\t%d: %s\n", p.Port, p.Open)
		}
		message += fmt.Sprintln()
	}
	_, err := fmt.Fprintln(out, message)
	return err
}

func init() {
	rootCmd.AddCommand(scanCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// scanCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// scanCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	scanCmd.PersistentFlags().IntSliceP("ports", "p", []int{22, 80, 443}, "ports to scan")
	scanCmd.PersistentFlags().StringP("port-range", "r", "", "port range to scan: ex) 1-3000")
	scanCmd.PersistentFlags().BoolP("use-udp", "u", false, "use UDP when true, default is false (use TCP).")
	scanCmd.PersistentFlags().DurationP("timeout", "t", 1*time.Second, "custom timeout duration")
}
