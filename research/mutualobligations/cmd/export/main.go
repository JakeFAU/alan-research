// Command export emits deterministic finite-graph summaries and shortest traces.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	model "github.com/JakeFAU/alan-research/research/mutualobligations"
)

func main() { os.Exit(run(os.Args[1:], os.Stdout, os.Stderr)) }
func run(args []string, out, errOut io.Writer) int {
	flags := flag.NewFlagSet("export", flag.ContinueOnError)
	flags.SetOutput(errOut)
	profile := flags.String(
		"profile", "all",
		"all or serial/serial, servicing/serial, serial/servicing, servicing/servicing",
	)
	naive := flags.Bool("assert-naive-completion", false, "intentionally exit 1 if any reachable state cannot complete")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	if flags.NArg() != 0 {
		fmt.Fprintln(errOut, "unexpected positional arguments")
		return 2
	}
	selected := []model.Profile{}
	for _, p := range model.Profiles() {
		if *profile == "all" || *profile == string(p) {
			selected = append(selected, p)
		}
	}
	if len(selected) == 0 {
		fmt.Fprintln(errOut, "unknown profile")
		return 2
	}
	reports := make([]model.Report, 0, len(selected))
	counterexample := false
	for _, p := range selected {
		r := model.Explore(p)
		if r.LocalSafetyViolations+r.LocalTransitionViolations+r.CustodyViolations+r.TrapMismatch != 0 {
			fmt.Fprintln(errOut, "apparatus invariant failure")
			return 2
		}
		witnesses := []*model.Witness{
			r.ShortestTrap, r.AdmittedWait, r.RecoveredWait, r.ShortestComplete,
			r.PostApplicationCustody, r.ReleasedCustody,
		}
		for _, w := range witnesses {
			if w != nil {
				if err := model.CheckTrace(p, w.Trace); err != nil {
					fmt.Fprintln(errOut, err)
					return 2
				}
			}
		}
		reports = append(reports, r)
		counterexample = counterexample || r.CannotComplete > 0
	}
	encoder := json.NewEncoder(out)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(reports); err != nil {
		fmt.Fprintln(errOut, err)
		return 2
	}
	if *naive && counterexample {
		fmt.Fprintln(errOut, "naive completion claim falsified: reachable serial dependency trap")
		return 1
	}
	return 0
}
