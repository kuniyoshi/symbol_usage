package cmd

import (
	"fmt"
	"log"
	"sort"

	"github.com/spf13/cobra"
	"github.com/kuniyoshi/symbol_usage/internal"
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

		symbols := internal.GetAllSymbols(index)
		sort.Strings(symbols)

		for _, symbol := range symbols {
			fmt.Println(symbol)
		}
	},
	Example: `  sy list /tmp/index.scip
  sy list ./index.scip`,
}

func init() {
	rootCmd.AddCommand(listCmd)
}