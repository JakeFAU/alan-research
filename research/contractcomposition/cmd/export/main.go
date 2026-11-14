// Command export emits the versioned deterministic Phase 3C.1 evidence packet.
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/JakeFAU/alan-research/research/contractcomposition"
)

func writeExport(w io.Writer) error {
	evidence, err := contractcomposition.BuildEvidence()
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(evidence)
}
func main() {
	if err := writeExport(os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
