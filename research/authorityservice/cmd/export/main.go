// Command export writes deterministic complete-graph counts and shortest witnesses.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	model "github.com/JakeFAU/alan-research/research/authorityservice"
	"io"
	"os"
)

func main() { os.Exit(run(os.Args[1:], os.Stdout, os.Stderr)) }
func run(args []string, out, errOut io.Writer) int {
	flags := flag.NewFlagSet("export", flag.ContinueOnError)
	flags.SetOutput(errOut)
	selected := flags.String("config", "all", "all or a frozen Policy/Scenario name")
	safety := flags.Bool("assert-authority-safe", false, "exit 1 when a reachable authority violation falsifies safety")
	completion := flags.Bool("assert-universal-completion", false, "exit 1 when a reachable state cannot safely complete")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	if flags.NArg() != 0 {
		fmt.Fprintln(errOut, "unexpected positional arguments")
		return 2
	}
	reports := []model.Report{}
	counterexample := false
	for _, c := range model.Configurations() {
		if *selected != "all" && *selected != c.Name() {
			continue
		}
		r := model.Explore(c)
		if r.OperationalInvariantViolations+r.CustodyViolations != 0 {
			fmt.Fprintln(errOut, "apparatus invariant failure")
			return 2
		}
		counterexample = counterexample || (*safety && r.UnsafeStates > 0) || (*completion && r.CannotCompleteSafely > 0)
		reports = append(reports, r)
	}
	if len(reports) == 0 {
		fmt.Fprintln(errOut, "unknown configuration")
		return 2
	}
	encoder := json.NewEncoder(out)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(reports); err != nil {
		fmt.Fprintln(errOut, err)
		return 2
	}
	if counterexample {
		fmt.Fprintln(errOut, "candidate claim falsified; preserved counterexample in report")
		return 1
	}
	return 0
}
