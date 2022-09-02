package cmd

import (
	"fmt"
	"os"

	"github.com/nao1215/morrigan/internal/embedded"
	"github.com/nao1215/morrigan/internal/print"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var licenseCmd = &cobra.Command{
	Use:   "license [PACKAGEs]",
	Short: "Print the license of the package (or subcommand) used in the morrigan",
	Long: `Print the license of the package (or subcommand) used in the morrigan.
`,
	Example: `  morrigan license zip`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := licenseRun(cmd, args); err != nil {
			print.Err(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(licenseCmd)
}

type license struct {
	pkg       string
	orgAuthor string
	orgURL    string
	license   string
}

func licenseRun(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		fmt.Println("The morrigan command user many codes that forked from other projects.")
		fmt.Println("All packages and subcommands not listed in the table are MIT LICENSE.")
		fmt.Println("We pay our respects to the author of the original code.")
		fmt.Println("")
		printPkgList()
		fmt.Println("")
		fmt.Println("If you want to read the detailed license terms, pass the package")
		fmt.Println("name given in the table as arguments.")
		fmt.Println("[e.g.]")
		fmt.Println(" $ morrigan license morrigan")
		return nil
	}

	for i, v := range args {
		if i != 0 {
			fmt.Println("---------------------------------------------------------------------------")
		}
		fmt.Printf("[%s]\n", v)

		license, err := embedded.License(v)
		if err != nil {
			print.Err("no such package or subcommand: " + v)
		}

		for _, l := range license {
			fmt.Println(l)
		}
	}
	return nil
}

func printPkgList() {
	l := []license{
		{
			pkg:       "gocrypt",
			orgAuthor: "Akshay Moghe",
			orgURL:    "https://github.com/amoghe/go-crypt",
			license:   "MIT",
		},
		{
			pkg:       "netcat",
			orgAuthor: "Vadym Fedorov",
			orgURL:    "https://github.com/vfedoroff/go-netcat",
			license:   "MIT",
		},
		{
			pkg:       "morrigan",
			orgAuthor: "Naohiro CHIKAMATSU",
			orgURL:    "https://github.com/nao1215/morrigan",
			license:   "MIT",
		},
		{
			pkg:       "ping",
			orgAuthor: "u-root project",
			orgURL:    "https://github.com/u-root/u-root",
			license:   "BSD 3-Clause",
		},
		{
			pkg:       "whris",
			orgAuthor: "harakeishi",
			orgURL:    "https://github.com/harakeishi/whris",
			license:   "MIT",
		},
		{
			pkg:       "zip",
			orgAuthor: "Alex Mullins",
			orgURL:    "https://github.com/yeka/zip",
			license:   "MIT",
		},
	}

	tableData := [][]string{}
	for _, v := range l {
		tableData = append(tableData, []string{v.pkg, v.orgAuthor, v.orgURL, v.license})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Package (CMD)", "Original Author", "Original URL", "LICENSE"})
	table.SetAutoWrapText(false)

	for _, v := range tableData {
		table.Append(v)
	}
	table.Render()
}
