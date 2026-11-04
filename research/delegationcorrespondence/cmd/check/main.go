package main

import (
	c "github.com/JakeFAU/alan-research/research/delegationcorrespondence"
	"os"
)

func main() { os.Exit(c.RunCheck(os.Args[1:], os.Stdout, os.Stderr)) }
