package main

import (
	"fmt"
	"github.com/JakeFAU/alan-research/research/residualantichains"
	"os"
)

func main() {
	if err := residualantichains.Export(os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
