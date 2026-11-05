package delegationcorrespondence

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type IndexEntry struct {
	Version, Name, Schedule, Directory, Classification string
	Manifest                                           json.RawMessage
	Interpretation                                     json.RawMessage
	Cut                                                json.RawMessage
}
type PhaseReport struct {
	Path, SHA256 string
	Verdict      RawVerdict
	Projection   Projection
}
type CaseReport struct {
	Name, ClaimedClassification, OriginalDirectory string
	CaseSHA256                                     string
	Final                                          RawVerdict
	Phases                                         []PhaseReport
	Files                                          map[string]string
	Outcomes                                       map[string]json.RawMessage
	Projection                                     Projection
	External                                       ExternalReport
	ContinuityFailure                              *Failure
}
type CorpusCheck struct {
	Version, Root, OriginalPrefix, IndexSHA256 string
	Cases                                      []CaseReport
	FilesVerified                              int
	SourceCommits                              []string
	Discrepancies                              []string
}
type ProjectionReport struct {
	Supported                       bool
	Reasons, Stutters, Abstractions []string
	Mismatch                        string
	Scope                           string
}
type ExternalReport struct {
	Valid, Present         bool
	Failure                *Failure
	Contacts, Applications int
	AcknowledgedLosses     int
	Limitations            []string
	Exposures              []Exposure
	References             []RowRef
	OwnerEpochs            map[string]map[string]int64
	Owners                 map[string]OwnerEvidence
	OwnerBindings          []OwnerBoundaryBinding
}

type OwnerRun struct {
	ID                        RunID
	ProgramID                 string
	Epoch, Sequence, Revision int64
	Projection                ProjectionState
}
type OwnerEvidence struct {
	OpenPath, OpenSHA256, PostOpenPath, PostOpenSHA256, PriorPath, PriorSHA256 string
	PID                                                                        int
	Runs                                                                       map[string]OwnerRun
	RenewalChecked                                                             bool
}
type OwnerBoundaryBinding struct {
	Path, Owner, Original string
	Event, PID            int
	Run                   RunID
	Epoch, GateOrdinal    int64
}
type Exposure struct {
	Original, Endpoint, WireSHA256               string
	Contacts, Applications                       int
	AllowedGates, RefusedGates, ReturnedOrdinals []int64
}

func resolve(root, prefix, original string) (string, error) {
	rel, e := filepath.Rel(prefix, original)
	if e != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) || filepath.IsAbs(rel) {
		return "", fmt.Errorf("path outside declared original prefix: %s", original)
	}
	return filepath.Join(root, rel), nil
}
func CheckCorpus(root, prefix string) (CorpusCheck, error) {
	out := CorpusCheck{Version: "alan/raw-correspondence-check/v1", Root: root, OriginalPrefix: prefix}
	index, e := os.ReadFile(filepath.Join(root, "case-index.json"))
	if e != nil {
		return out, e
	}
	out.IndexSHA256 = sha(index)
	var entries []IndexEntry
	if e = strictJSON(index, &entries); e != nil {
		return out, e
	}
	names := map[string]bool{}
	for _, entry := range entries {
		if names[entry.Name] {
			return out, fmt.Errorf("duplicate case %s", entry.Name)
		}
		names[entry.Name] = true
		dir, e := resolve(root, prefix, entry.Directory)
		if e != nil {
			return out, e
		}
		report := CaseReport{Name: entry.Name, ClaimedClassification: entry.Classification, OriginalDirectory: entry.Directory, Files: map[string]string{}, Outcomes: map[string]json.RawMessage{}}
		db, e := os.ReadFile(filepath.Join(dir, "digests.json"))
		if e != nil {
			return out, e
		}
		var digests map[string]string
		if e = strictJSON(db, &digests); e != nil {
			return out, e
		}
		paths := []string{}
		for path := range digests {
			paths = append(paths, path)
		}
		sort.Strings(paths)
		externalFiles := map[string][]byte{}
		for _, path := range paths {
			full, e := resolve(dir, dir, filepath.Join(dir, path))
			if e != nil {
				return out, e
			}
			data, e := os.ReadFile(full)
			if e != nil {
				return out, e
			}
			hash := sha(data)
			if hash != digests[path] {
				return out, fmt.Errorf("digest mismatch %s/%s", entry.Name, path)
			}
			out.FilesVerified++
			report.Files[path] = hash
			if strings.HasSuffix(path, ".jsonl") || strings.HasSuffix(path, ".json") {
				externalFiles[path] = data
			}
			if path == "case.json" {
				report.CaseSHA256 = hash
				var actual IndexEntry
				if e = strictJSON(data, &actual); e != nil {
					return out, e
				}
				if string(bytesJSON(entry)) != string(bytesJSON(actual)) {
					return out, fmt.Errorf("index/case mismatch %s", entry.Name)
				}
			}
			if path == "raw-tables.json" || strings.HasSuffix(path, "/inspect.json") {
				v := CheckRaw(data)
				report.Phases = append(report.Phases, PhaseReport{path, hash, v, projectChecked(data, v)})
				if path == "raw-tables.json" || path == "final/inspect.json" || path == "after-grant/inspect.json" {
					report.Final = v
				}
			}
			if strings.HasSuffix(path, "outcomes.json") || strings.HasSuffix(path, "open.json") || strings.HasSuffix(path, "run.json") {
				report.Outcomes[path] = json.RawMessage(data)
			}
		}
		if report.Final.InputSHA256 == "" {
			return out, fmt.Errorf("missing final raw input %s", entry.Name)
		}
		finalPath := "final/inspect.json"
		if _, ok := report.Files["raw-tables.json"]; ok {
			finalPath = "raw-tables.json"
		} else if _, ok := report.Files["after-grant/inspect.json"]; ok {
			finalPath = "after-grant/inspect.json"
		}
		rawData := externalFiles[finalPath]
		report.Projection = projectChecked(rawData, report.Final)
		report.External = CheckExternal(rawData, externalFiles)
		report.ContinuityFailure = CheckContinuity(report.Phases, externalFiles)
		for _, owner := range []string{"first-owner/open.json", "fresh-owner/open.json", "reopened/inspect.json"} {
			if _, exists := report.Files[owner]; exists {
				report.Projection.Supported = false
				report.Projection.Reasons = append(report.Projection.Reasons, "whole_store_crash_paired_expansion_not_established")
				report.Projection.Scope = "retained logical service subsequence only; chronological process recovery unsupported"
				break
			}
		}
		var source map[string]json.RawMessage
		if strictJSON(externalFiles["source.json"], &source) == nil {
			var head string
			json.Unmarshal(source["HEAD"], &head)
			head = strings.TrimSpace(head)
			found := false
			for _, x := range out.SourceCommits {
				found = found || x == head
			}
			if !found {
				out.SourceCommits = append(out.SourceCommits, head)
			}
		}
		for _, d := range report.Final.Discrepancies {
			out.Discrepancies = append(out.Discrepancies, entry.Name+":"+d)
		}
		out.Cases = append(out.Cases, report)
	}
	return out, nil
}
func projectSummary(v RawVerdict) ProjectionReport {
	p := ProjectionReport{Scope: "retained logical service subsequence; no process availability correspondence", Stutters: []string{"lifecycle outcome selection", "grant-ready transport and selection", "terminal protocol after both original results consumed"}, Abstractions: []string{"explicit original-ID alpha renaming", "whole-action many-to-one B tuples for nonrevoking reciprocal trace", "equality-preserving per-case result atoms"}}
	if !v.Valid {
		p.Reasons = append(p.Reasons, "invalid_raw_history")
		return p
	}
	if len(v.RejectedAuthorityOperations) > 0 {
		p.Reasons = append(p.Reasons, "rejected_authority_control_not_in_lategrant_issuance")
	}
	for _, e := range v.Effects {
		if e.Source.Subject == "b" && e.Allowed && !v.Custody.BSerialGuard {
			p.Reasons = append(p.Reasons, "program_service_order_outside_A_servicing_B_serial")
		}
	}
	if v.Revocations > 0 {
		p.Reasons = append(p.Reasons, "separate_grant_with_revocation_not_in_frozen_configuration")
	}
	if len(v.Custody.Originals) != 2 {
		p.Reasons = append(p.Reasons, "not_two_original_reciprocal_history")
	}
	for _, e := range v.Effects {
		if len(e.Gates) > 1 {
			p.Reasons = append(p.Reasons, "repeated_physical_exposure_not_atomic_model_apply")
		}
		if e.Allowed && !e.Result && len(e.Gates) > 0 {
			p.Reasons = append(p.Reasons, "unknown_or_uncalled_prefix_requires_external_evidence")
		}
		grantAlreadyPresent := false
		for _, grant := range v.Grants {
			for _, fact := range v.OrderedFacts {
				if fact.Kind == "lifecycle" && fact.OperationID == grant && fact.Order < e.Order {
					grantAlreadyPresent = true
				}
			}
		}
		if !e.Allowed && (e.Reason != "authority_grant_missing" || grantAlreadyPresent) {
			p.Reasons = append(p.Reasons, "explicit_selector_policy_narrows_implicit_model_choice")
		}
	}
	if len(v.Discrepancies) > 0 {
		p.Mismatch = v.Discrepancies[0]
		p.Reasons = append(p.Reasons, p.Mismatch)
	}
	p.Supported = len(p.Reasons) == 0
	return p
}
