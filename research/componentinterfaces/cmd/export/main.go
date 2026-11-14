package main

import (
	"fmt"
	"github.com/JakeFAU/alan-research/research/componentinterfaces"
	"os"
)

func main() {
	if len(os.Args) != 1 {
		fmt.Fprintln(os.Stderr, "usage: export (run from repository root)")
		os.Exit(2)
	}
	if err := componentinterfaces.Export(os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
