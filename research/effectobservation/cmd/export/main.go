// Command export emits the complete split effect-observation graph.
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/JakeFAU/alan-research/research/effectobservation"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, out, errOut io.Writer) int {
	if len(args) != 0 {
		fmt.Fprintln(errOut, "usage: export")
		return 2
	}
	encoder := json.NewEncoder(out)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(effectobservation.Explore()); err != nil {
		fmt.Fprintln(errOut, err)
		return 1
	}
	return 0
}
