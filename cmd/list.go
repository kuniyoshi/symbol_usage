package cmd

import (
	"fmt"
	"log"
	"sort"

	"github.com/kuniyoshi/symbol_usage/internal"
	"github.com/spf13/cobra"
)

var (
	listVerbose bool
)

var listCmd = &cobra.Command{
	Use:   "list <scip-file>",
	Short: "List all symbols in SCIP index",
	Long:  `List all symbols found in the SCIP index file.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		scipFile := args[0]

		index, err := internal.ReadSCIPIndex(scipFile)
		if err != nil {
			log.Fatalf("Failed to read SCIP index: %v", err)
		}

		if listVerbose {
			symbolsWithSCIP := internal.GetAllSymbolsVerbose(index)
			sort.Slice(symbolsWithSCIP, func(i, j int) bool {
				return symbolsWithSCIP[i].Formatted < symbolsWithSCIP[j].Formatted
			})

			for _, sym := range symbolsWithSCIP {
				fmt.Printf("%-50s => %s\n", sym.Formatted, sym.SCIP)
			}
		} else {
			symbols := internal.GetAllSymbols(index)
			sort.Strings(symbols)

			for _, symbol := range symbols {
				fmt.Println(symbol)
			}
		}
	},
	Example: `  sy list /tmp/index.scip
  sy list ./index.scip`,
}

func init() {
	listCmd.Flags().BoolVarP(&listVerbose, "verbose", "v", false, "Show verbose output with SCIP symbol names")
	rootCmd.AddCommand(listCmd)
}
