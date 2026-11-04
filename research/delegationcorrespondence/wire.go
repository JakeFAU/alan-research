// Package delegationcorrespondence checks a bounded retained-history contract.
// Its authority checker imports neither runtime code nor model predicates.
package delegationcorrespondence

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"unicode/utf8"
)

const authorityVersion = "alan/request-authority/v1"

type RunID struct {
	Subject    string
	Generation int64
}
type Value struct {
	Kind   string
	Number int64
	Text   string
	Flag   bool
}
type Permission struct {
	Kind, Source, Target, Operation string
	Revision                        int64
	Obligation                      string
}
type AuthorityRoot struct {
	Run        RunID
	RootRun    RunID
	Permission Permission
}
type Claim struct {
	Version                               string
	Root                                  AuthorityRoot
	Recipient, RecipientRoot              RunID
	RecipientRevision                     int64
	OriginalRequestID, ServiceOperationID string
}
type Revoke struct {
	Version string
	Root    AuthorityRoot
}
type Use struct {
	Version, Path string
	Root          AuthorityRoot
	GrantID       string
}
type Decision struct {
	Version, Namespace      string
	Root                    AuthorityRoot
	RootID, RecipientRootID string
}
type Witness struct {
	Version, Namespace                                          string
	Use                                                         *Use
	RootID, ServingRootID                                       string
	Allowed                                                     bool
	Reason                                                      string
	OriginalRequestID, ServiceOperationID, RequestObservationID string
	RequestConsumedSequence, GrantDecisionSequence              int64
}
type Intent struct {
	Authority    *Use `json:",omitempty"`
	ID           string
	Run          RunID
	Sequence     int64
	Ordinal      int
	Kind, Target string
	Generation   int64
	Obligation   string
	Payload      map[string]Value
}
type ProjectionState struct {
	PC     int
	Data   map[string]Value
	Status string
}
type Source struct {
	Run                       RunID
	Epoch, Sequence, Revision int64
}
type Authorization struct {
	ID, IntentID     string
	Allowed          bool
	Reason           string
	Source           Source
	ClosureOperation string
}
type AuthorizationRequest struct {
	Intent   Intent
	Revision int64
}
type LifecycleRequest struct {
	Intent           Intent
	SourceRevision   int64
	Operation        string
	ExpectedRevision int64
}
type TargetState struct {
	Run                       RunID
	ProgramID, Owner          string
	Revision, Epoch, Sequence int64
	Retired, RetryCancelled   bool
	Projection                ProjectionState
}
type LifecycleResult struct {
	Authority        *Decision `json:",omitempty"`
	Request          LifecycleRequest
	DecisionSequence int64
	Accepted         bool
	Reason           string
	Successor        int64
	Before, After    *TargetState
}
type WireRequest struct {
	Key, Profile, Subject, Authorization string
	Generation, Epoch, Sequence, Delta   int64
}
type Receipt struct {
	Key, Profile           string
	Delta, Mutation, Total int64
	Process                string
	PID                    int
}
type ExternalAttempt struct {
	Authority       *Witness `json:",omitempty"`
	Intent          Intent
	AuthorizationID string
	Request         WireRequest
	Allowed         bool
	Reason          string
}
type AdapterAttempt struct {
	Authority       *Witness `json:",omitempty"`
	Intent          Intent
	AuthorizationID string
	Prepared        Prepared
	Allowed         bool
	Reason          string
}
type Prepared struct {
	Request                           AdapterRequest
	Digest, Protocol, Endpoint, Model string
	MaxOutputTokens, MaxResponseBytes int64
	Wire                              string
}
type AdapterRequest struct {
	ID, Key, Capability, Effect string
	Input, Output               json.RawMessage
	Payload                     map[string]Value
	Assets                      []json.RawMessage
}
type Candidate struct {
	ID, Key, Disposition                           string
	Payload                                        map[string]Value
	HTTPStatus                                     int
	RequestID, ResponseID, ReturnedModel, Response string
	Usage                                          json.RawMessage
}
type AdapterResult struct {
	Digest      string
	Candidate   Candidate
	Values      map[string]Value
	Disposition string
}
type Observation struct {
	ID                       string
	Run                      RunID
	Kind, Provenance, Source string
	SourceGeneration         int64
	Obligation, Operation    string
	Payload                  map[string]Value
}
type Admission struct {
	TransportID, Kind, LogicalID string
	Source                       RunID
	Target                       string
	Generation                   int64
	Obligation, Operation        string
	Payload                      map[string]Value
	IntentID                     string
	TerminalSequence             int64
}
type AdmissionResult struct {
	TransportID   string
	Accepted      bool
	Reason        string
	Bound         *RunID
	ObservationID string
}
type Reservation struct {
	OperationID, RootID, Kind string
	Run                       RunID
	Sequence                  int64
	Allowed                   bool
	Reason, AttemptDigest     string
}
type Charge struct {
	RootID, OperationID, Kind                        string
	Invocations, CreatedRuns, Effects, DecisionOrder int64
}
type Gate struct {
	Version, Namespace string
	Ordinal            int64
	Source             Source
	AttemptDigest      string
	Allowed            bool
	Reason             string
}
type Fact struct {
	AuthorityAdmission        *Gate `json:",omitempty"`
	Version, ID, Kind         string
	Order                     int64
	RootIDs                   []string
	Run                       RunID
	Sequence                  int64
	OperationID, RecordDigest string
}
type Root struct {
	ID, Version             string
	Run                     RunID
	Members                 []RunID
	Templates, Capabilities []string
	Limits                  Limits
}
type Limits struct{ MaxInvocations, MaxCreatedRuns, MaxDepth, MaxEffects, MaxInFlight int64 }
type Profile struct {
	RequestAuthority string
	Version          string
	Roots            []Root
	Templates        []json.RawMessage
}
type Subject struct {
	ID                      string
	Generation              int64
	Owner                   string
	Revision                int64
	Retired, RetryCancelled bool
}
type InitialRun struct {
	ID         RunID
	ProgramID  string
	Projection ProjectionState
}
type Deployment struct {
	Subjects []Subject
	Runs     []InitialRun
}
type Origin struct {
	Subject   Subject
	ProgramID string
	Initial   ProjectionState
}
type StaticMember struct {
	RootID string
	Run    RunID
}
type Expr struct {
	Ref     string
	Literal *Value
	Offset  *int64
}
type Filter struct {
	Kind, Provenance                          string
	Source, Generation, Obligation, Operation Expr
}
type Effect struct {
	Authority                      *Use `json:",omitempty"`
	Kind                           string
	Target, Generation, Obligation Expr
	Payload                        map[string]Expr
}
type Instruction struct {
	Op, Dst           string
	A, B              Expr
	Target, Otherwise int
	Match             *Filter
	Bind              map[string]string
	Effect            *Effect
	Updates           map[string]Expr
	Outcome           string
}
type Program struct {
	Semantics string
	Code      []Instruction
	Author    json.RawMessage `json:",omitempty"`
}
type Freeze struct {
	InvocationID, Trigger                            string
	Run                                              RunID
	Sequence, DecisionOrder                          int64
	Status                                           string
	ReplyIntents, InvocationIntents, ExternalIntents []string
}
type Outcome struct {
	ID, InvocationID, OperationID, TemplateID string
	Caller                                    RunID
	Status, Reason                            string
	ResponsePresent                           bool
	Response                                  map[string]Value
	DecisionOrder                             int64
}
type Row map[string]json.RawMessage
type RawSnapshot struct {
	Namespace  string
	KeyVersion int
	Settings   json.RawMessage
	Tables     map[string][]Row
	Roots      json.RawMessage
}

func bytesJSON(v any) []byte {
	b, e := json.Marshal(v)
	if e != nil {
		panic(e)
	}
	return b
}
func sha(b []byte) string                { s := sha256.Sum256(b); return hex.EncodeToString(s[:]) }
func identity(v ...any) string           { return string(bytesJSON(v)) }
func digest(domain string, v any) string { return sha(bytesJSON([]any{domain, v})) }
func rootID(namespace string, r RunID) string {
	return digest("alan/resource-profile/v1", []any{namespace, r})
}
func equal(a, b any) bool { return reflect.DeepEqual(a, b) }

// strictJSON rejects duplicate keys recursively before typed decoding. Numbers
// never traverse float64; the second decoder rejects unknown typed fields.
func strictJSON(b []byte, out any) error {
	if !utf8.Valid(b) {
		return fmt.Errorf("invalid UTF-8 JSON")
	}
	d := json.NewDecoder(bytes.NewReader(b))
	d.UseNumber()
	var walk func() error
	walk = func() error {
		t, e := d.Token()
		if e != nil {
			return e
		}
		switch x := t.(type) {
		case json.Delim:
			switch x {
			case '{':
				seen := map[string]bool{}
				for d.More() {
					k, e := d.Token()
					if e != nil {
						return e
					}
					s, ok := k.(string)
					if !ok || seen[s] {
						return fmt.Errorf("duplicate or nonstring key %v", k)
					}
					seen[s] = true
					if e = walk(); e != nil {
						return e
					}
				}
			case '[':
				for d.More() {
					if e = walk(); e != nil {
						return e
					}
				}
			default:
				return fmt.Errorf("unexpected delimiter")
			}
			_, e = d.Token()
			return e
		}
		return nil
	}
	if e := walk(); e != nil {
		return e
	}
	if _, e := d.Token(); e != io.EOF {
		return fmt.Errorf("trailing JSON")
	}
	d = json.NewDecoder(bytes.NewReader(b))
	d.DisallowUnknownFields()
	d.UseNumber()
	return d.Decode(out)
}
func canonical(b []byte, out any) error {
	if e := strictJSON(b, out); e != nil {
		return e
	}
	if !bytes.Equal(b, bytesJSON(out)) {
		return fmt.Errorf("noncanonical descriptor")
	}
	return nil
}
func textCell(r Row, k string) string {
	var s string
	if e := strictJSON(r[k], &s); e != nil {
		fail("cell_type", k, e.Error())
	}
	return s
}
func intCell(r Row, k string) int64 {
	var n int64
	if e := strictJSON(r[k], &n); e != nil {
		fail("cell_type", k, e.Error())
	}
	return n
}
func nullableText(r Row, k string) string {
	if bytes.Equal(r[k], []byte("null")) {
		return ""
	}
	return textCell(r, k)
}
func rawCell(r Row, k string) json.RawMessage { return json.RawMessage(textCell(r, k)) }
func decodeCell[T any](r Row, k string) T {
	var out T
	if e := canonical(rawCell(r, k), &out); e != nil {
		fail("descriptor_encoding", k, e.Error())
	}
	return out
}
func need(ok bool, edge, id, detail string) {
	if !ok {
		fail(edge, id, detail)
	}
}
func fail(edge, id, detail string) { panic(Failure{Edge: edge, ID: id, Detail: detail}) }
func locator(id string) (RunID, int64, int) {
	var a []json.RawMessage
	need(strictJSON([]byte(id), &a) == nil && len(a) == 4, "canonical_intent_id", id, "")
	var s string
	var g, q int64
	var o int
	need(strictJSON(a[0], &s) == nil && strictJSON(a[1], &g) == nil && strictJSON(a[2], &q) == nil && strictJSON(a[3], &o) == nil, "canonical_intent_id", id, "")
	need(s != "" && g > 0 && q > 0 && o >= 0 && identity(s, g, q, o) == id, "canonical_intent_id", id, "")
	return RunID{s, g}, q, o
}
