package main

import (
	"encoding/json"
	"fmt"
	c "github.com/JakeFAU/alan-research/research/delegationcorrespondence"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "mutations RAW_SNAPSHOT")
		os.Exit(2)
	}
	data, e := os.ReadFile(os.Args[1])
	if e != nil {
		fmt.Fprintln(os.Stderr, e)
		os.Exit(2)
	}
	v := c.CheckRaw(data)
	if !v.Valid {
		fmt.Fprintln(os.Stderr, "requires a valid reciprocal source")
		os.Exit(2)
	}
	r := c.RawMutations(data)
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if e = enc.Encode(r); e != nil {
		fmt.Fprintln(os.Stderr, e)
		os.Exit(2)
	}
	for _, m := range r {
		if m.Verdict.Valid {
			os.Exit(1)
		}
	}
}
