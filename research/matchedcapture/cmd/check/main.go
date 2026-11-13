package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/JakeFAU/alan-research/research/matchedcapture"
)

func main() { os.Exit(run(os.Args[1:], os.Stdout, os.Stderr)) }
func run(args []string, out, diagnostics io.Writer) int {
	flags := flag.NewFlagSet("matchedcapture", flag.ContinueOnError)
	flags.SetOutput(diagnostics)
	pair := flags.String("pair", "", "sealed matched pair directory")
	if e := flags.Parse(args); e != nil {
		return 2
	}
	if *pair == "" || flags.NArg() != 0 {
		fmt.Fprintln(diagnostics, "supply exactly one directory with -pair")
		return 2
	}
	report, e := matchedcapture.CheckDirectory(*pair)
	if e != nil {
		fmt.Fprintln(diagnostics, e)
		return 2
	}
	enc := json.NewEncoder(out)
	enc.SetIndent("", "  ")
	if e = enc.Encode(report); e != nil {
		fmt.Fprintln(diagnostics, e)
		return 2
	}
	if !report.Supported {
		return 1
	}
	return 0
}
