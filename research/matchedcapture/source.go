package matchedcapture

import (
	"bytes"
	"compress/gzip"
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime/debug"
	"sort"
	"strconv"
	"strings"
)

//go:embed testdata/prepared-source.json.gz
var preparedGzip []byte

const PreparedSHA256 = "557e832fc84e7cbdb47d1d4767938f80ddd1ca1d77c8c647cb665a3436546994"

func preparedBytes() []byte {
	z, e := gzip.NewReader(bytes.NewReader(preparedGzip))
	if e != nil {
		panic(e)
	}
	b, e := io.ReadAll(z)
	if e != nil {
		panic(e)
	}
	if e = z.Close(); e != nil {
		panic(e)
	}
	require(digest(b).SHA256 == PreparedSHA256, "builtin prepared source", "anchor corrupt")
	return b
}
func loadDirectory(root string) (map[string][]byte, map[string]Digest, error) {
	files := map[string][]byte{}
	digests := map[string]Digest{}
	st, e := os.Lstat(root)
	if e != nil {
		return nil, nil, e
	}
	if !st.IsDir() {
		return nil, nil, fmt.Errorf("pair must be a real directory")
	}
	e = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !d.Type().IsRegular() {
			return fmt.Errorf("nonregular member: %s", path)
		}
		name, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		name = filepath.ToSlash(name)
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		digests[name] = digest(b)
		if strings.HasSuffix(name, ".json") || strings.HasSuffix(name, ".jsonl") || strings.HasSuffix(name, ".log") {
			files[name] = b
		}
		return nil
	})
	return files, digests, e
}
func checkSeal(f map[string][]byte, d map[string]Digest) {
	m := read[struct {
		Version string
		Files   map[string]Digest
	}](f, "members.json")
	require(m.Version == Version, "members.json", "version")
	require(len(m.Files) == len(d)-1, "members.json", "member coverage")
	for _, n := range sortedKeys(m.Files) {
		h := m.Files[n]
		require(filepath.IsLocal(n) && filepath.ToSlash(filepath.Clean(n)) == n && n != "members.json", "members.json", "unsafe member "+n)
		actual, exists := d[n]
		require(exists && h.Bytes >= 0 && validSHA(h.SHA256) && actual == h, n, "missing member or length/hash mismatch")
	}
	for _, n := range sortedKeys(f) {
		b := f[n]
		if strings.HasSuffix(n, ".json") {
			_, e := parseJSON(b)
			require(e == nil, n, fmt.Sprint(e))
		} else if strings.HasSuffix(n, ".jsonl") {
			_ = lines[json.RawMessage](f, n)
		}
	}
}

type source struct {
	Version        string            `json:"version"`
	Prepared       json.RawMessage   `json:"prepared"`
	GoVersion      string            `json:"go_version"`
	GOOS           string            `json:"goos"`
	GOARCH         string            `json:"goarch"`
	Race           bool              `json:"race"`
	Environment    map[string]string `json:"environment"`
	TestArgv       []string          `json:"test_argv"`
	BuildInfo      debug.BuildInfo   `json:"build_info"`
	OwnerBinary    Digest            `json:"owner_binary"`
	EndpointBinary Digest            `json:"endpoint_binary"`
}

func checkSource(f map[string][]byte, p pair) {
	s := read[source](f, "source.json")
	require(s.Version == Version, "source.json", "version")
	c, e := CompareJSON(s.Prepared, preparedBytes())
	require(e == nil && c.StructuralEqual, "source.json.prepared", "does not match externally frozen complete manifest")
	require(s.GoVersion != "" && s.GoVersion == s.BuildInfo.GoVersion && s.GOOS == "darwin" && s.GOARCH == "arm64", "source.json", "unsupported build platform/version")
	require((p.Mode == "race") == s.Race && (p.Mode == "normal" || p.Mode == "race"), "source.json", "race mode join")
	require(s.OwnerBinary.Bytes > 0 && validSHA(s.OwnerBinary.SHA256) && s.EndpointBinary.Bytes > 0 && validSHA(s.EndpointBinary.SHA256), "source.json", "missing binary identities")
	require(len(s.TestArgv) > 1 && s.TestArgv[0] != "", "source.json", "test executable argv")
	checkTestSelection(s.TestArgv[1:])
	settings := map[string]string{}
	for _, v := range s.BuildInfo.Settings {
		_, exists := settings[v.Key]
		require(!exists, "source.json", "duplicate build setting")
		settings[v.Key] = v.Value
	}
	require(settings["GOOS"] == s.GOOS && settings["GOARCH"] == s.GOARCH && settings["CGO_ENABLED"] == "1" && settings["-tags"] == "integration" && (settings["-race"] == "true") == s.Race, "source.json", "build settings")
	require(s.BuildInfo.Path == "github.com/JakeFAU/alan/internal/durabletest.test", "source.json", "owner build path")
	dep := false
	for _, v := range s.BuildInfo.Deps {
		require(v != nil, "source.json", "null build dependency")
		if v.Path == "github.com/mattn/go-sqlite3" {
			dep = v.Version == "v1.14.52" && v.Replace == nil
		}
	}
	require(dep, "source.json", "pinned sqlite dependency")
	for _, k := range []string{"GOTOOLCHAIN", "GOWORK", "GOPROXY", "GOSUMDB", "GOMODCACHE", "GOCACHE", "GOFLAGS", "GOENV", "CGO_ENABLED"} {
		_, ok := s.Environment[k]
		require(ok, "source.json.environment", "missing "+k)
	}
	require(len(s.Environment) == 9 && s.Environment["GOTOOLCHAIN"] == "local" && s.Environment["GOWORK"] == "off" && s.Environment["GOPROXY"] == "off" && s.Environment["GOSUMDB"] == "off", "source.json.environment", "unfrozen execution environment")
}
func validSHA(s string) bool {
	if len(s) != 64 {
		return false
	}
	for _, c := range s {
		if !(c >= '0' && c <= '9' || c >= 'a' && c <= 'f') {
			return false
		}
	}
	return true
}

func checkArchiveShape(d map[string]Digest) {
	allowed := map[string]bool{}
	for _, name := range []string{"members.json", "pair.json", "source.json", "seed/provision.json", "seed/manifest.json", "seed/copies.json"} {
		allowed[name] = true
	}
	capture := func(dir string) {
		for _, name := range []string{"inspect.json", "typed-inspection.json", "outcomes.json", "replay.json", "cost-sidecar.json", "history-000.json", "runtime.sqlite", "runtime.sqlite-wal", "runtime.sqlite-shm"} {
			allowed[dir+"/"+name] = true
		}
	}
	capture("seed/initial")
	for _, branch := range []string{"zero", "applied"} {
		for _, name := range []string{"protocol.jsonl", "recovery-open.json", "controlled-revocation.json", "run.json", "first-owner/open.json", "first-owner/pre-call.json", "first-owner/worker-release.json", "first-owner/boundaries.jsonl", "first-owner/abort.json", "first-owner/process.log"} {
			allowed[branch+"/"+name] = true
		}
		for _, cut := range []string{"first-owner/post-open", "stopped", "reopened", "final"} {
			capture(branch + "/" + cut)
		}
		for _, endpoint := range []string{"endpoint-a", "endpoint-b"} {
			for _, name := range []string{"process.log", "ready.json", "journal.jsonl", "stopped.json", "final.json", "resource.sqlite", "resource.sqlite-wal", "resource.sqlite-shm"} {
				allowed[branch+"/"+endpoint+"/"+name] = true
			}
		}
	}
	names := []string{}
	for name := range d {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		require(allowed[name], name, "outside complete frozen archive schema")
	}
}

// Use the owning Go flag grammar (including --, double dash, and separated
// values), while admitting only the frozen non-selection-changing driver flags.
func checkTestSelection(args []string) {
	flags := flag.NewFlagSet("capture", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	var selected string
	selections := 0
	flags.Func("test.run", "", func(v string) error { selections++; selected = v; return nil })
	count := flags.Int("test.count", 1, "")
	flags.Duration("test.timeout", 0, "")
	flags.String("test.testlogfile", "", "")
	flags.Bool("test.paniconexit0", false, "")
	flags.BoolFunc("test.v", "", func(v string) error {
		if v == "test2json" {
			return nil
		}
		_, e := strconv.ParseBool(v)
		return e
	})
	err := flags.Parse(args)
	require(err == nil && flags.NArg() == 0 && selections == 1 &&
		selected == "^TestPhase3B7MatchedCapture$" && *count == 1,
		"source.json.test_argv", "unsupported or nonunique effective capture selection")
}
