package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/JakeFAU/alan-research/research/delegationcorrespondence"
)

func main() { os.Exit(run(os.Args[1:], os.Stdout, os.Stderr)) }

func run(args []string, output, diagnostics io.Writer) int {
	flags := flag.NewFlagSet("revocation", flag.ContinueOnError)
	flags.SetOutput(diagnostics)
	dir := flags.String("case", "", "directory of the declared binding case")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	if *dir == "" || flags.NArg() != 0 {
		fmt.Fprintln(diagnostics, "supply exactly one case directory with -case")
		return 2
	}
	report, err := delegationcorrespondence.CheckRevocationDirectory(*dir)
	if err != nil {
		fmt.Fprintln(diagnostics, err)
		return 2
	}
	encoder := json.NewEncoder(output)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(report); err != nil {
		fmt.Fprintln(diagnostics, err)
		return 2
	}
	if !report.Supported {
		return 1
	}
	return 0
}
