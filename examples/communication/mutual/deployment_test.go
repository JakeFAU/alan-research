package mutual_test

import (
	"crypto/sha256"
	"encoding/hex"
	"reflect"
	"testing"

	"github.com/JakeFAU/alan-research/examples/communication/mutual"
)

func TestNew(t *testing.T) {
	for _, profile := range []mutual.Profile{mutual.SerialSerial, mutual.ServiceA, mutual.ServiceB, mutual.ServiceBoth} {
		t.Run(string(profile), func(t *testing.T) {
			p, err := mutual.New(profile)
			if err != nil {
				t.Fatal(err)
			}
			if len(p.Runs) != 2 || len(p.Subjects) != 2 || len(p.Programs) != 2 || p.Composition != nil {
				t.Fatal("fixture must contain exactly two ordinary roots")
			}
			messages, inputs := 0, 0
			for _, permission := range p.Permissions {
				switch permission.Kind {
				case "message":
					messages++
				case "input":
					inputs++
				default:
					t.Fatal("fixture acquired an unexpected capability", permission)
				}
			}
			if messages != 4 || inputs != 2 {
				t.Fatal("wrong reciprocal/hold boundaries", p.Permissions)
			}
			for _, program := range p.Programs {
				digest := sha256.Sum256(program.Bytes)
				if program.ID != hex.EncodeToString(digest[:]) {
					t.Fatal("artifact identity does not bind exact public program bytes")
				}
			}
			again, err := mutual.New(profile)
			if err != nil || !reflect.DeepEqual(p, again) {
				t.Fatal("construction is not deterministic", err)
			}
		})
	}
	if p, err := mutual.New("unknown"); err == nil || len(p.Runs) != 0 {
		t.Fatal("unknown service policy yielded a runnable provision")
	}
}
