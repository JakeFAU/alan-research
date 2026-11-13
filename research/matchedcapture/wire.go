// Package matchedcapture checks source-conditioned paired histories without
// executing producer code or opening evidence databases.
package matchedcapture

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	dc "github.com/JakeFAU/alan-research/research/delegationcorrespondence"
)

const Version = "alan/phase3b7-matched/v1"
const service = `["a",1,4,0]`

type Digest struct {
	Bytes  int64  `json:"bytes"`
	SHA256 string `json:"sha256"`
}

func digest(b []byte) Digest {
	s := sha256.Sum256(b)
	return Digest{int64(len(b)), hex.EncodeToString(s[:])}
}

type Verdict struct {
	Valid   bool
	Reasons []string
}
type checkFailure string

func require(ok bool, path, why string) {
	if !ok {
		panic(checkFailure(path + ": " + why))
	}
}
func stage(fn func()) (v Verdict) {
	v.Reasons = []string{}
	defer func() {
		if p := recover(); p != nil {
			if f, ok := p.(checkFailure); ok {
				v.Reasons = append(v.Reasons, string(f))
			} else {
				panic(p)
			}
		}
	}()
	fn()
	v.Valid = true
	return
}
func decode[T any](b []byte, path string) T {
	_, e := parseJSON(b)
	require(e == nil, path, fmt.Sprint(e))
	var v T
	d := json.NewDecoder(bytes.NewReader(b))
	d.DisallowUnknownFields()
	e = d.Decode(&v)
	require(e == nil, path, fmt.Sprint(e))
	return v
}
func jsonBytes(v any) []byte {
	b, e := json.Marshal(v)
	if e != nil {
		panic(e)
	}
	return b
}
func same(a, b any) bool {
	r, e := CompareJSON(jsonBytes(a), jsonBytes(b))
	return e == nil && r.StructuralEqual
}
func read[T any](f map[string][]byte, p string) T {
	b, ok := f[p]
	require(ok, p, "missing member")
	return decode[T](b, p)
}
func lines[T any](f map[string][]byte, p string) []T {
	b, ok := f[p]
	require(ok, p, "missing journal")
	out := []T{}
	if len(b) == 0 {
		return out
	}
	require(b[len(b)-1] == '\n', p, "incomplete final journal record")
	for i, s := range strings.Split(strings.TrimSuffix(string(b), "\n"), "\n") {
		out = append(out, decode[T]([]byte(s), fmt.Sprintf("%s[%d]", p, i)))
	}
	return out
}

type event struct {
	AuthorityAdmissionOrdinal                        int64 `json:",omitempty"`
	Scenario, Token, Store, Boundary, Phase, Subject string
	PID                                              int
	Generation, Epoch, Sequence                      int64
	Intent, Attempt                                  string
	Root                                             string `json:",omitempty"`
	InFlight                                         []string
}
type ready struct {
	URL, Token, Store string
	PID               int
}
type opened struct {
	PID                                      int
	Scenario, Token, Store, Namespace, Error string
	Runs                                     map[string]dc.OwnerRun
}
type record struct {
	Ordinal int             `json:"ordinal"`
	Action  string          `json:"action"`
	Role    string          `json:"role"`
	PID     int             `json:"pid"`
	Data    json.RawMessage `json:"data"`
}
type pair struct {
	Version, Status, Mode string
	ControllerPID         int      `json:"controller_pid"`
	Completed             []string `json:"completed_branches"`
}
type process struct {
	Config  event
	PID     int
	Command []string
}

func sortedKeys[V any](m map[string]V) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
