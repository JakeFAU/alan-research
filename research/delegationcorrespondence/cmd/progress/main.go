package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	m "github.com/JakeFAU/alan-research/research/authorityservice"
	c "github.com/JakeFAU/alan-research/research/delegationcorrespondence"
)

func main() { os.Exit(run(os.Args[1:], os.Stdout, os.Stderr)) }
func run(args []string, output, errors io.Writer) int {
	flags := flag.NewFlagSet("progress", flag.ContinueOnError)
	flags.SetOutput(errors)
	selected := flags.String("config", "all", "all historical progress queries or S2/SeparateGrantRevocation")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	if flags.NArg() != 0 {
		fmt.Fprintln(errors, "unexpected positional arguments")
		return 2
	}
	scenarios := []string{"LateGrant", "Revocable"}
	if *selected == m.GrantRevocationConfig().Name() {
		scenarios = []string{"SeparateGrantRevocation"}
	} else if *selected != "all" {
		fmt.Fprintln(errors, "unknown progress configuration")
		return 2
	}
	out := []c.ProgressReport{}
	for _, s := range scenarios {
		r, e := c.QueryProgress(m.Config{Policy: "S2", Scenario: s})
		if e != nil {
			fmt.Fprintln(errors, e)
			return 2
		}
		out = append(out, r)
	}
	enc := json.NewEncoder(output)
	enc.SetIndent("", "  ")
	if e := enc.Encode(out); e != nil {
		fmt.Fprintln(errors, e)
		return 2
	}
	for _, r := range out {
		if r.PredicateMismatches != 0 {
			return 1
		}
	}
	return 0
}
