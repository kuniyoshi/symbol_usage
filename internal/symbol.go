package internal

import (
	"fmt"
	"strings"

	"github.com/sourcegraph/scip/bindings/go/scip"
)

type SymbolUsage struct {
	Symbol      string
	File        string
	Line        int32
	IsReference bool
}

func convertUserSymbolToSCIP(userSymbol string) string {
	parts := strings.Split(userSymbol, ".")
	if len(parts) < 2 {
		return userSymbol
	}

	scipParts := []string{"scip", ""}

	for i, part := range parts {
		if i == len(parts)-1 {
			scipParts = append(scipParts, part+"().")
		} else {
			scipParts = append(scipParts, part+"/")
		}
	}

	return strings.Join(scipParts, " ")
}

func findSymbolReferences(index *scip.Index, targetSymbol string) []SymbolUsage {
	var results []SymbolUsage
	seenDefinitions := make(map[string]bool)

	for _, doc := range index.Documents {
		for _, occ := range doc.Occurrences {
			if matchesSymbol(occ.Symbol, targetSymbol) {
				usage := SymbolUsage{
					Symbol:      occ.Symbol,
					File:        doc.RelativePath,
					Line:        occ.Range[0],
					IsReference: occ.SymbolRoles&int32(scip.SymbolRole_Definition) == 0,
				}

				key := fmt.Sprintf("%s:%d", usage.File, usage.Line)
				if !usage.IsReference {
					if seenDefinitions[key] {
						continue
					}
					seenDefinitions[key] = true
				}

				results = append(results, usage)
			}
		}
	}

	return results
}

func matchesSymbol(scipSymbol, targetSymbol string) bool {
	if scipSymbol == targetSymbol {
		return true
	}

	if strings.Contains(scipSymbol, targetSymbol) {
		return true
	}

	convertedTarget := convertUserSymbolToSCIP(targetSymbol)
	return strings.Contains(scipSymbol, convertedTarget)
}

func findCallers(index *scip.Index, targetSymbol string) []SymbolUsage {
	var callers []SymbolUsage
	symbolRefs := make(map[string][]SymbolUsage)

	for _, doc := range index.Documents {
		currentSymbol := ""

		for _, symbol := range doc.Symbols {
			for _, rel := range symbol.Relationships {
				if rel.IsReference && matchesSymbol(rel.Symbol, targetSymbol) {
					currentSymbol = symbol.Symbol
					break
				}
			}
		}

		if currentSymbol != "" {
			for _, occ := range doc.Occurrences {
				if occ.Symbol == currentSymbol {
					usage := SymbolUsage{
						Symbol: currentSymbol,
						File:   doc.RelativePath,
						Line:   occ.Range[0],
					}
					symbolRefs[currentSymbol] = append(symbolRefs[currentSymbol], usage)
				}
			}
		}
	}

	for _, usages := range symbolRefs {
		if len(usages) > 0 {
			callers = append(callers, usages[0])
		}
	}

	return callers
}