package nominal

import (
	"os"
	"strings"
	"testing"
)

// These exact artifacts were compiled with research main 3e64c837, before
// ordinary author validation rejected the attacks. They retain the independent
// nominal judgment on already-encoded programs; no new compiler admission is
// implied. program checks the canonical bytes and content identity.
func TestHistoricalNominalDataFlowAttacks(t *testing.T) {
	for _, tc := range []struct{ name, id, reason string }{
		{"literal-forgery", "4a18e6f400d3b4de49830957ddcdd1b06fab130b782955c162e59db17007dd64", `operand "copy" requires handle[helper], got string`},
		{"overwrite-with-data", "17cd472a0ac5de029f7b8a83dd1541127441c9a0acba19d374f5133c7917a03f", `register "copy" changes sort from handle[helper] to string`},
		{"loop-launder", "71707a7891d154daeb4b5c1949e6096f4f29067ee73f3a6b9c1710e1ef9f1e27", `pc 6 (set) steps/2:foreach/body/1:set:carry: register "carry" changes sort from string to handle[helper]`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			raw, err := os.ReadFile("testdata/" + tc.name + "-v1.json")
			if err != nil {
				t.Fatal(err)
			}
			p, err := program(tc.id, raw)
			if err != nil {
				t.Fatalf("historical artifact: %v", err)
			}
			err = checkOperands(p)
			if err == nil || !strings.Contains(err.Error(), tc.reason) {
				t.Fatalf("nominal attack not refused at expected boundary: %v", err)
			}
		})
	}
}
