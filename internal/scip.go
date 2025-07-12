package internal

import (
	"fmt"
	"os"
	"strings"

	"github.com/sourcegraph/scip/bindings/go/scip"
	"google.golang.org/protobuf/proto"
)

func ReadSCIPIndex(filepath string) (*scip.Index, error) {
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

func DisplaySymbolUsage(index *scip.Index, targetSymbol string) error {
	references := findSymbolReferences(index, targetSymbol)
	callers := findCallers(index, targetSymbol)

	callerMap := make(map[string]bool)
	for _, caller := range callers {
		callerMap[FormatSymbolName(caller.Symbol)] = true
	}

	displayedSymbols := make(map[string]bool)

	for _, ref := range references {
		symbolName := FormatSymbolName(ref.Symbol)

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
				subName := FormatSymbolName(subRef.Symbol)
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

func FormatSymbolName(scipSymbol string) string {
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

func GetAllSymbols(index *scip.Index) []string {
	symbolMap := make(map[string]bool)
	var symbols []string

	for _, doc := range index.Documents {
		for _, symbol := range doc.Symbols {
			if !symbolMap[symbol.Symbol] {
				symbolMap[symbol.Symbol] = true
				symbols = append(symbols, FormatSymbolName(symbol.Symbol))
			}
		}
		
		for _, occ := range doc.Occurrences {
			if !symbolMap[occ.Symbol] {
				symbolMap[occ.Symbol] = true
				symbols = append(symbols, FormatSymbolName(occ.Symbol))
			}
		}
	}

	return symbols
}