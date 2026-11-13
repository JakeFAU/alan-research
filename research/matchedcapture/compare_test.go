package matchedcapture

import "testing"

// A projection-only comparator, float64 coercion, or permissive JSON decoder
// would miss these differences in the declared complete observation.
func TestCompleteJSON(t *testing.T) {
	for _, tc := range []struct {
		name, left, right, path string
		equal, malformed        bool
	}{
		{"object-order", `{"a":1,"b":[]}`, ` { "b":[], "a":1 }`, "", true, false},
		{"hidden-state", `{"StoreView":{},"Tables":{"runs":[{"projection":"PC=4"}]}}`, `{"StoreView":{},"Tables":{"runs":[{"projection":"PC=5"}]}}`, `$.Tables.runs[0].projection`, false, false},
		{"duplicate", `{"x":1,"x":1}`, `{"x":1}`, "", false, true},
		{"null-empty", `{"x":null}`, `{"x":[]}`, `$.x`, false, false},
		{"array-order", `{"x":[1,2]}`, `{"x":[2,1]}`, `$.x[0]`, false, false},
		{"large-number", `{"x":9007199254740992}`, `{"x":9007199254740993}`, `$.x`, false, false},
		{"descriptor-exact", `{"x":"{\"a\":1}"}`, `{"x":"{ \"a\":1}"}`, `$.x`, false, false},
		{"receipt-evidence", `{"ReceiptEvidence":null}`, `{"ReceiptEvidence":[]}`, `$.ReceiptEvidence`, false, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r, err := CompareJSON([]byte(tc.left), []byte(tc.right))
			if (err != nil) != tc.malformed {
				t.Fatalf("malformed=%v: %v", tc.malformed, err)
			}
			if err != nil {
				return
			}
			if r.StructuralEqual != tc.equal {
				t.Fatalf("equality=%v, want %v", r.StructuralEqual, tc.equal)
			}
			if tc.path != "" && (len(r.Paths) == 0 || r.Paths[0] != tc.path) {
				t.Fatalf("paths %v want first %s", r.Paths, tc.path)
			}
		})
	}
}

func TestUnambiguousDifferencePath(t *testing.T) {
	r, e := CompareJSON([]byte(`{"a.b":1}`), []byte(`{"a.b":2}`))
	if e != nil || len(r.Paths) != 1 || r.Paths[0] != `$["a.b"]` {
		t.Fatalf("ambiguous field path: %v %v", e, r.Paths)
	}
}

func TestMalformedBytesStillReportByteEquality(t *testing.T) {
	data := []byte(`{"x":1,"x":1}`)
	r, e := CompareJSON(data, data)
	if e == nil || !r.ByteEqual || r.StructuralEqual {
		t.Fatalf("byte equality lost on invalid JSON: %+v %v", r, e)
	}
}
