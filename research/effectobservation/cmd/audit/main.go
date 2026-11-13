// Command audit emits the independent graph audit and semantic controls.
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/JakeFAU/alan-research/research/effectobservation"
)

func main() { os.Exit(run(os.Args[1:], os.Stdout, os.Stderr)) }

func run(args []string, out, errOut io.Writer) int {
	if len(args) != 0 {
		fmt.Fprintln(errOut, "usage: audit")
		return 2
	}
	report, err := effectobservation.AuditGraph(effectobservation.Explore())
	if err != nil {
		fmt.Fprintln(errOut, err)
		return 1
	}
	controls, err := effectobservation.AuditControls()
	if err != nil {
		fmt.Fprintln(errOut, err)
		return 1
	}
	result := struct {
		Report   effectobservation.AuditReport     `json:"report"`
		Controls []effectobservation.ControlResult `json:"controls"`
	}{report, controls}
	encoder := json.NewEncoder(out)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(result); err != nil {
		fmt.Fprintln(errOut, err)
		return 1
	}
	return 0
}
