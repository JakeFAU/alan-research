package main

import (
	"encoding/json"
	"fmt"
	m "github.com/JakeFAU/alan-research/research/authorityservice"
	c "github.com/JakeFAU/alan-research/research/delegationcorrespondence"
	"os"
)

func main() {
	if len(os.Args) != 1 {
		fmt.Fprintln(os.Stderr, "progress takes no arguments")
		os.Exit(2)
	}
	out := []c.ProgressReport{}
	for _, s := range []string{"LateGrant", "Revocable"} {
		r, e := c.QueryProgress(m.Config{Policy: "S2", Scenario: s})
		if e != nil {
			fmt.Fprintln(os.Stderr, e)
			os.Exit(2)
		}
		out = append(out, r)
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if e := enc.Encode(out); e != nil {
		fmt.Fprintln(os.Stderr, e)
		os.Exit(2)
	}
	for _, r := range out {
		if r.PredicateMismatches != 0 {
			os.Exit(1)
		}
	}
}
