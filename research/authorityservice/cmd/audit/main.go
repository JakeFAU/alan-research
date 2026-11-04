// Command audit exports independent finite-graph checks and retained hostile
// controls. It does not assert exhaustive enumeration of the history language.
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	model "github.com/JakeFAU/alan-research/research/authorityservice"
)

type report struct {
	Coverage string
	Graphs   []model.GraphAudit
	Controls model.OracleCorpusReport
}

func main() { os.Exit(run(os.Args[1:], os.Stdout, os.Stderr)) }
func run(args []string, out, errOut io.Writer) int {
	if len(args) != 0 {
		fmt.Fprintln(errOut, "audit takes no arguments")
		return 2
	}
	result := report{
		Coverage: "one independently checked shortest BFS prefix per finite state; every outgoing arc; equality of independent ledgers at merged targets; not enumeration of all histories",
		Graphs:   []model.GraphAudit{},
	}
	for _, config := range model.Configurations() {
		graph, err := model.AuditGraph(config)
		if err != nil {
			fmt.Fprintln(errOut, err)
			return 2
		}
		result.Graphs = append(result.Graphs, graph)
	}
	corpus, err := model.OracleCorpus()
	if err != nil {
		fmt.Fprintln(errOut, err)
		return 2
	}
	result.Controls = corpus
	encoder := json.NewEncoder(out)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(result); err != nil {
		fmt.Fprintln(errOut, err)
		return 2
	}
	return 0
}
