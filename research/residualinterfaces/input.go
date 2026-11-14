package residualinterfaces

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

const predecessorPath = "docs/foundational/phase3c3/evidence/task1/export.json.gz"
const predecessorCompressed = "86a9d22f9808a67dd0066b6b77dd2b3c01a2125d10863c6fbe1369097f6624c8"
const predecessorRaw = "a073a22fe8ac62e05f2d15b0f5146b19db9606e38ce4619261e7e4108e27269d"

type source struct {
	Path   string
	SHA256 string
}
type sourceBinding struct {
	Path   string
	SHA256 string
	Use    string
}
type inputBinding struct {
	Path             string
	CompressedBytes  int
	CompressedSHA256 string
	RawBytes         int
	RawSHA256        string
	Schema           string
	Sources          []source
}
type predecessor struct {
	Schema      string
	Alphabet    []symbol
	Monitors    []automaton
	Plants      []automaton
	Controllers []automaton
	Sources     []source
}

func hash(b []byte) string { return fmt.Sprintf("%x", sha256.Sum256(b)) }
func readInput(root string) (predecessor, inputBinding, error) {
	var p predecessor
	binding := inputBinding{Path: predecessorPath}
	compressed, err := os.ReadFile(filepath.Join(root, predecessorPath))
	if err != nil {
		return p, binding, err
	}
	if hash(compressed) != predecessorCompressed {
		return p, binding, fmt.Errorf("predecessor compressed identity mismatch")
	}
	z, err := gzip.NewReader(bytes.NewReader(compressed))
	if err != nil {
		return p, binding, err
	}
	raw, err := io.ReadAll(z)
	closeErr := z.Close()
	if err != nil {
		return p, binding, err
	}
	if closeErr != nil {
		return p, binding, closeErr
	}
	if len(raw) != 25384964 || hash(raw) != predecessorRaw {
		return p, binding, fmt.Errorf("predecessor raw identity mismatch")
	}
	// This intentionally selects descriptor and provenance fields only. All old
	// products, counts, histories and outcome comparisons are discarded by decoding.
	if err = json.Unmarshal(raw, &p); err != nil {
		return p, binding, err
	}
	if p.Schema != "alan/phase3c3/v1" || len(p.Sources) != 27 || len(p.Plants) != 72 || len(p.Monitors) != 5 || len(p.Controllers) != 324 {
		return p, binding, fmt.Errorf("predecessor roster mismatch")
	}
	if len(p.Alphabet) != len(labels) {
		return p, binding, fmt.Errorf("alphabet size mismatch")
	}
	for i, s := range p.Alphabet {
		if s.ID != labels[i] {
			return p, binding, fmt.Errorf("alphabet mismatch")
		}
	}
	for _, s := range p.Sources {
		b, e := os.ReadFile(filepath.Join(root, s.Path))
		if e != nil {
			return p, binding, e
		}
		if hash(b) != s.SHA256 {
			return p, binding, fmt.Errorf("predecessor source mismatch: %s", s.Path)
		}
	}
	for _, group := range [][]automaton{p.Plants, p.Monitors, p.Controllers} {
		for _, a := range group {
			if err := validate(a); err != nil {
				return p, binding, err
			}
		}
	}
	binding = inputBinding{Path: predecessorPath, CompressedBytes: len(compressed), CompressedSHA256: hash(compressed), RawBytes: len(raw), RawSHA256: hash(raw), Schema: p.Schema, Sources: p.Sources}
	return p, binding, nil
}
func sourceManifest(root string, old []source) ([]sourceBinding, error) {
	uses := map[string]string{}
	for _, s := range old {
		uses[s.Path] = "hash-checked predecessor provenance"
	}
	for _, path := range []string{"go.mod", "go.sum"} {
		uses[path] = "module"
	}
	for _, path := range []string{"00-contract.md", "01-argument.md", "02-evidence-schema.md"} {
		uses["docs/foundational/phase3c4/"+path] = "hash-checked specification"
	}
	uses[predecessorPath] = "decoded pinned automata"
	for _, dir := range []string{"research/residualinterfaces", "research/residualinterfaces/cmd/export"} {
		entries, err := os.ReadDir(filepath.Join(root, dir))
		if err != nil {
			return nil, err
		}
		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".go") && !strings.HasSuffix(entry.Name(), "_test.go") {
				uses[dir+"/"+entry.Name()] = "compiled producer"
			}
		}
	}
	paths := []string{}
	for path := range uses {
		paths = append(paths, path)
	}
	slices.Sort(paths)
	out := []sourceBinding{}
	for _, path := range paths {
		b, err := os.ReadFile(filepath.Join(root, path))
		if err != nil {
			return nil, err
		}
		out = append(out, sourceBinding{Path: path, SHA256: hash(b), Use: uses[path]})
	}
	return out, nil
}
