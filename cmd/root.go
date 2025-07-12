package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sy",
	Short: "Symbol Usage - SCIP-based code navigation tool",
	Long: `Symbol Usage は SCIP を使ってコードベースを読むためのプログラムです。
	
Find symbol references and their callers/callees in your codebase.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
