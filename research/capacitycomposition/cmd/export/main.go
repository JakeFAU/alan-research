// Command export writes the frozen Phase 3C.2 evidence JSON to stdout.
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	capacity "github.com/JakeFAU/alan-research/research/capacitycomposition"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

func sourceHashes() ([]capacity.SourceHash, error) {
	paths := []string{"go.mod", "go.sum", "docs/foundational/phase3c2/00-contract.md", "docs/foundational/phase3c2/01-argument.md", "docs/foundational/phase3c2/02-evidence-schema.md"}
	for _, dir := range []string{"research/capacitycomposition", "research/contractcomposition", "research/twoagent"} {
		err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() && strings.HasSuffix(path, ".go") {
				paths = append(paths, path)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	slices.Sort(paths)
	out := []capacity.SourceHash{}
	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		sum := sha256.Sum256(data)
		out = append(out, capacity.SourceHash{Path: filepath.ToSlash(path), SHA256: hex.EncodeToString(sum[:])})
	}
	return out, nil
}

// run expects the repository root as cwd and accepts no arguments.
func run(w io.Writer, args []string) error {
	if len(args) != 0 {
		return fmt.Errorf("usage: go run ./research/capacitycomposition/cmd/export (from repository root)")
	}
	hashes, err := sourceHashes()
	if err != nil {
		return err
	}
	e, err := capacity.BuildFiniteEvidence()
	if err != nil {
		return err
	}
	e.Sources = hashes
	e.Probes, err = capacity.ProbeRepresentatives()
	if err != nil {
		return err
	}
	return writeEvidence(w, e)
}
func writeEvidence(w io.Writer, e capacity.Evidence) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(e)
}
func main() {
	if err := run(os.Stdout, os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
