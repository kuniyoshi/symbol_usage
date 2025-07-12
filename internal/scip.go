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

func DisplaySymbolUsageVerbose(index *scip.Index, targetSymbol string) error {
	fmt.Printf("Searching for symbol: %s\n", targetSymbol)
	fmt.Printf("(Will also match SCIP patterns containing this symbol)\n\n")

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
			fmt.Printf("%s    %-50s [%s]\n", prefix, symbolName, ref.Symbol)

			subRefs := findSymbolReferences(index, ref.Symbol)
			displayedSubs := make(map[string]bool)
			for _, subRef := range subRefs {
				subName := FormatSymbolName(subRef.Symbol)
				if !displayedSubs[subName] && subName != symbolName {
					displayedSubs[subName] = true
					fmt.Printf("         %-50s [%s]\n", subName, subRef.Symbol)
				}
			}
		} else {
			fmt.Printf("%s    %-50s [%s]\n", prefix, symbolName, ref.Symbol)
		}
	}

	return nil
}

func FormatSymbolName(scipSymbol string) string {
	parts := strings.Fields(scipSymbol)
	if len(parts) < 3 {
		return scipSymbol
	}

	// Extract the symbol part (last part after the space-separated components)
	symbolPart := parts[len(parts)-1]

	// Clean up the symbol part
	symbolPart = strings.TrimSuffix(symbolPart, "/")
	symbolPart = strings.TrimSuffix(symbolPart, "().")
	symbolPart = strings.TrimSuffix(symbolPart, ".")

	// Extract package path from backticks if present
	if strings.Contains(symbolPart, "`") {
		if start := strings.Index(symbolPart, "`"); start != -1 {
			if end := strings.LastIndex(symbolPart, "`"); end != -1 && end > start {
				packagePath := symbolPart[start+1 : end]
				suffix := symbolPart[end+1:]

				// Simplify package path
				simplified := simplifyPackagePath(packagePath)
				return simplified + suffix
			}
		}
	}

	return symbolPart
}

func simplifyPackagePath(path string) string {
	// Remove common prefixes and just keep the meaningful parts
	parts := strings.Split(path, "/")

	// Find the start of the meaningful part
	startIdx := 0
	for i, part := range parts {
		if part == "github.com" || part == "golang.org" || part == "gopkg.in" {
			// Skip the domain and the next part (usually username/org)
			if i+2 < len(parts) {
				startIdx = i + 2
			}
			break
		}
	}

	if startIdx < len(parts) {
		return strings.Join(parts[startIdx:], "/")
	}

	return path
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

type SymbolInfo struct {
	SCIP      string
	Formatted string
}

func GetAllSymbolsVerbose(index *scip.Index) []SymbolInfo {
	symbolMap := make(map[string]bool)
	var symbols []SymbolInfo

	for _, doc := range index.Documents {
		for _, symbol := range doc.Symbols {
			if !symbolMap[symbol.Symbol] {
				symbolMap[symbol.Symbol] = true
				symbols = append(symbols, SymbolInfo{
					SCIP:      symbol.Symbol,
					Formatted: FormatSymbolName(symbol.Symbol),
				})
			}
		}

		for _, occ := range doc.Occurrences {
			if !symbolMap[occ.Symbol] {
				symbolMap[occ.Symbol] = true
				symbols = append(symbols, SymbolInfo{
					SCIP:      occ.Symbol,
					Formatted: FormatSymbolName(occ.Symbol),
				})
			}
		}
	}

	return symbols
}
