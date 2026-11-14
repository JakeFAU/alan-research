package main

import (
	"fmt"
	"github.com/JakeFAU/alan-research/research/residualinterfaces"
	"os"
)

func main() {
	if err := residualinterfaces.Export(os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
