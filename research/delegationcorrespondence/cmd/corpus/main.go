package main

import (
	"encoding/json"
	"fmt"
	c "github.com/JakeFAU/alan-research/research/delegationcorrespondence"
	"os"
)

func main() {
	if len(os.Args) != 1 {
		fmt.Fprintln(os.Stderr, "corpus takes no arguments")
		os.Exit(2)
	}
	r, e := c.GenerateCorpus()
	if e != nil {
		fmt.Fprintln(os.Stderr, e)
		os.Exit(2)
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if e = enc.Encode(r); e != nil {
		fmt.Fprintln(os.Stderr, e)
		os.Exit(2)
	}
}
