package cmd

// whris is forked from https://github.com/harakeishi/whris
// (MIT License)
//
// The MIT License (MIT)
//
// Copyright 2022 harakeishi. (original author)
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/likexian/whois"
	"github.com/nao1215/morrigan/internal/print"
	"github.com/spf13/cobra"
)

var whrisCmd = &cobra.Command{
	Use:   "whris DOMAIN_NAMEs",
	Short: "Displays management information for IPs associated with the domain",
	Long: `"whris" outputs the target domain and IP from the input domain.
	
It is as well as the administrator information for that IP: administrator name,
network name, range of IPs to be managed, and country name.`,
	Example: `  morrigan whris example.com`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := whris(cmd, args); err != nil {
			print.Err(err)
			os.Exit(1)
		}
		return nil
	},
}

func init() {
	whrisCmd.Flags().BoolP("verbose", "v", false, "verbose output")
	rootCmd.AddCommand(whrisCmd)
}

func whris(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("whris subcommand need one or more argument (you specify domain name)")
	}

	v, err := cmd.Flags().GetBool("verbose")
	if err != nil {
		return err
	}

	for n, domain := range args {
		if n > 0 {
			fmt.Println()
		}
		if _, err := resolve(domain, v); err != nil {
			return err
		}
	}
	return nil
}

type networkAdmin struct {
	ipRange string
	netName string
	country string
	admin   string
}

type domainSummary struct {
	targetDomain        string
	targetIP            string
	whoisResponseServer string
	whoisResponse       string
	parseResult         []networkAdmin
}

func resolve(domain string, verbose bool) (domainSummary, error) {
	var summary domainSummary
	summary.whoisResponseServer = "whois.apnic.net"
	summary.targetDomain = domain
	addr, err := net.ResolveIPAddr("ip", domain)
	if err != nil {
		return domainSummary{}, err
	}

	summary.targetIP = addr.String()
	summary.whoisResponse, err = whois.Whois(summary.targetIP, "whois.iana.org")
	if err != nil {
		return domainSummary{}, err
	}

	summary.setWhoisResponseServerFromWhoisResponse()
	summary.ParseWhoisResponse()
	if !summary.parseCheck() {
		summary.parseResult = summary.parseResult[1:]
		summary.whoisResponse, err = whois.Whois(summary.targetIP, summary.whoisResponseServer)
		if err != nil {
			return domainSummary{}, err
		}
		summary.ParseWhoisResponse()
	}
	summary.printResult(verbose)
	return summary, nil
}

func (s *domainSummary) ParseWhoisResponse() {
	paragraph := s.breakDownWhoisResponseIntoParagraphs()
	for _, v := range paragraph {
		tmp := networkAdmin{}
		row := strings.Split(v, "\n")
		for _, val := range row {
			col := strings.Split(val, ":")
			switch col[0] {
			case "inetnum", "NetRange":
				tmp.ipRange = strings.TrimSpace(col[1])
			case "netname", "NetName":
				tmp.netName = strings.TrimSpace(col[1])
			case "country", "Country":
				tmp.country = strings.TrimSpace(col[1])
			case "descr", "Organization", "organization", "owner":
				if tmp.admin == "" {
					tmp.admin = strings.TrimSpace(col[1])
				}
			}
		}
		if tmp.ipRange != "" {
			s.parseResult = append(s.parseResult, tmp)
		}
	}
}

func (s *domainSummary) printResult(verbose bool) {
	fmt.Println("[Basic Information]")
	fmt.Printf(" Target domain: %s\n", s.targetDomain)
	fmt.Printf(" Target ip    : %s\n", s.targetIP)
	fmt.Printf(" Network Admin: %s\n", s.parseResult[len(s.parseResult)-1].admin)
	fmt.Printf(" Network name : %s\n", s.parseResult[len(s.parseResult)-1].netName)
	fmt.Printf(" ip range     : %s\n", s.parseResult[len(s.parseResult)-1].ipRange)
	fmt.Printf(" country      : %s\n", s.parseResult[len(s.parseResult)-1].country)
	if verbose {
		fmt.Println("[Detailed Information]")
		for i, v := range s.parseResult {
			fmt.Printf(" %d:\n", i)
			fmt.Printf(" Network Admin: %s\n", v.admin)
			fmt.Printf(" Network name : %s\n", v.netName)
			fmt.Printf(" ip range     : %s\n", v.ipRange)
			fmt.Printf(" country      : %s\n", v.country)
		}
	}
}

func (s *domainSummary) breakDownWhoisResponseIntoParagraphs() []string {
	switch s.whoisResponseServer {
	case "whois.arin.net":
		return strings.Split(s.whoisResponse, "#")
	case "whois.apnic.net", "whois.ripe.net", "whois.lacnic.net":
		return strings.Split(s.whoisResponse, "%")
	default:
		return strings.Split(s.whoisResponse, "\n\n")
	}
}

func (s *domainSummary) setWhoisResponseServerFromWhoisResponse() {
	tmp := strings.Split(s.whoisResponse, "\n")
	for _, v := range tmp {
		col := strings.Split(v, ":")
		if col[0] == "refer" {
			s.whoisResponseServer = strings.TrimSpace(col[1])
			break
		}
	}
}

func (s *domainSummary) parseCheck() bool {
	list := []string{"apnic", "arin", "ripe", "lacnic"}
	for _, v := range list {
		if strings.Contains(strings.ToLower(s.parseResult[1].netName), v) {
			s.whoisResponseServer = fmt.Sprintf("whois.%s.net", v)
			return false
		}
	}
	return true
}
