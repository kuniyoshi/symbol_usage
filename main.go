package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/sourcegraph/scip/bindings/go/scip"
	"google.golang.org/protobuf/proto"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <scip-file> <symbol>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Example: %s /tmp/index.scip Foo.bar\n", os.Args[0])
		os.Exit(1)
	}

	scipFile := os.Args[1]
	targetSymbol := os.Args[2]

	index, err := readSCIPIndex(scipFile)
	if err != nil {
		log.Fatalf("Failed to read SCIP index: %v", err)
	}

	if err := displaySymbolUsage(index, targetSymbol); err != nil {
		log.Fatalf("Failed to find symbol usage: %v", err)
	}
}

func readSCIPIndex(filepath string) (*scip.Index, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var index scip.Index
	if err := proto.Unmarshal(data, &index); err != nil {
		return nil, fmt.Errorf("failed to unmarshal SCIP index: %w", err)
	}

	return &index, nil
}

func displaySymbolUsage(index *scip.Index, targetSymbol string) error {
	references := findSymbolReferences(index, targetSymbol)
	callers := findCallers(index, targetSymbol)

	callerMap := make(map[string]bool)
	for _, caller := range callers {
		callerMap[formatSymbolName(caller.Symbol)] = true
	}

	displayedSymbols := make(map[string]bool)

	for _, ref := range references {
		symbolName := formatSymbolName(ref.Symbol)

		if displayedSymbols[symbolName] {
			continue
		}
		displayedSymbols[symbolName] = true

		prefix := " "
		if strings.Contains(ref.Symbol, targetSymbol) || strings.Contains(symbolName, targetSymbol) {
			prefix = "*"
		}

		if callerMap[symbolName] {
			fmt.Printf("%s    %s\n", prefix, symbolName)

			subRefs := findSymbolReferences(index, ref.Symbol)
			displayedSubs := make(map[string]bool)
			for _, subRef := range subRefs {
				subName := formatSymbolName(subRef.Symbol)
				if !displayedSubs[subName] && subName != symbolName {
					displayedSubs[subName] = true
					fmt.Printf("         %s\n", subName)
				}
			}
		} else {
			fmt.Printf("%s    %s\n", prefix, symbolName)
		}
	}

	return nil
}

func formatSymbolName(scipSymbol string) string {
	parts := strings.Fields(scipSymbol)
	if len(parts) < 3 {
		return scipSymbol
	}

	var result []string
	for i := 2; i < len(parts); i++ {
		part := parts[i]
		part = strings.TrimSuffix(part, "/")
		part = strings.TrimSuffix(part, "().")
		part = strings.TrimSuffix(part, ".")
		if part != "" {
			result = append(result, part)
		}
	}

	return strings.Join(result, ".")
}
