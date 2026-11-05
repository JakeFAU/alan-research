package delegationcorrespondence

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
)

type Expectation struct {
	RawValid bool
	Mismatch string
}

// RunCheck returns 0 for declared matched outcomes, 1 for a semantic discrepancy,
// and 2 for usage/input apparatus errors. Expectations are explicit input.
func RunCheck(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("check", flag.ContinueOnError)
	fs.SetOutput(stderr)
	root := fs.String("root", "", "extracted raw corpus root")
	prefix := fs.String("original-prefix", "", "original absolute corpus prefix; default root")
	input := fs.String("input", "", "single raw JSON snapshot")
	binding := fs.String("binding-case", "", "single explicit binding case directory")
	expect := fs.String("expect", "", "explicit JSON map of case name to {RawValid,Mismatch}")
	if e := fs.Parse(args); e != nil {
		return 2
	}
	modes := 0
	for _, value := range []string{*root, *input, *binding} {
		if value != "" {
			modes++
		}
	}
	if fs.NArg() != 0 || modes != 1 {
		fmt.Fprintln(stderr, "supply exactly one of -root, -input or -binding-case")
		return 2
	}
	enc := json.NewEncoder(stdout)
	enc.SetIndent("", "  ")
	if *binding != "" {
		report, err := CheckBindingDirectory(*binding)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 2
		}
		if err = enc.Encode(report); err != nil {
			fmt.Fprintln(stderr, err)
			return 2
		}
		if report.ContinuityFailure != nil {
			return 1
		}
		for _, phase := range report.Phases {
			if !phase.Verdict.Valid {
				return 1
			}
		}
		if len(report.Phases) > 0 && !report.External.Valid {
			return 1
		}
		return 0
	}
	if *input != "" {
		data, e := os.ReadFile(*input)
		if e != nil {
			fmt.Fprintln(stderr, e)
			return 2
		}
		var outer RawSnapshot
		if e = strictJSON(data, &outer); e != nil {
			fmt.Fprintln(stderr, e)
			return 2
		}
		r := CheckRaw(data)
		if e = enc.Encode(r); e != nil {
			fmt.Fprintln(stderr, e)
			return 2
		}
		if !r.Valid {
			return 1
		}
		return 0
	}
	if *prefix == "" {
		*prefix = *root
	}
	r, e := CheckCorpus(*root, *prefix)
	if e != nil {
		fmt.Fprintln(stderr, e)
		return 2
	}
	expected := map[string]Expectation{}
	if *expect != "" {
		data, e := os.ReadFile(*expect)
		if e != nil {
			fmt.Fprintln(stderr, e)
			return 2
		}
		if e = strictJSON(data, &expected); e != nil {
			fmt.Fprintln(stderr, e)
			return 2
		}
		if len(expected) != len(r.Cases) {
			fmt.Fprintln(stderr, "expectations must cover every case exactly")
			return 2
		}
	}
	exit := 0
	for _, c := range r.Cases {
		want := Expectation{RawValid: true}
		if *expect != "" {
			var ok bool
			want, ok = expected[c.Name]
			if !ok {
				fmt.Fprintln(stderr, "missing expectation", c.Name)
				return 2
			}
		}
		if c.Final.Valid != want.RawValid || c.Projection.Mismatch != want.Mismatch || c.ContinuityFailure != nil || c.Final.Valid && !c.External.Valid || c.Projection.Error != "" {
			exit = 1
		}
		for _, p := range c.Phases {
			if want.RawValid && !p.Verdict.Valid {
				exit = 1
			}
		}
	}
	if e = enc.Encode(r); e != nil {
		fmt.Fprintln(stderr, e)
		return 2
	}
	return exit
}
