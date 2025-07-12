package cmd

import (
	"log"

	"github.com/kuniyoshi/symbol_usage/internal"
	"github.com/spf13/cobra"
)

var (
	findVerbose bool
)

var findCmd = &cobra.Command{
	Use:   "find <scip-file> <symbol>",
	Short: "Find symbol usage in SCIP index",
	Long:  `Find references and callers/callees of a symbol in the SCIP index.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		scipFile := args[0]
		targetSymbol := args[1]

		index, err := internal.ReadSCIPIndex(scipFile)
		if err != nil {
			log.Fatalf("Failed to read SCIP index: %v", err)
		}

		if findVerbose {
			if err := internal.DisplaySymbolUsageVerbose(index, targetSymbol); err != nil {
				log.Fatalf("Failed to find symbol usage: %v", err)
			}
		} else {
			if err := internal.DisplaySymbolUsage(index, targetSymbol); err != nil {
				log.Fatalf("Failed to find symbol usage: %v", err)
			}
		}
	},
	Example: `  sy find /tmp/index.scip Foo.bar
  sy find ./index.scip MyClass.method`,
}

func init() {
	findCmd.Flags().BoolVarP(&findVerbose, "verbose", "v", false, "Show verbose output with SCIP symbol names")
	rootCmd.AddCommand(findCmd)
}
