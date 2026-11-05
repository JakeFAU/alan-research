package delegationcorrespondence

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// CheckBindingDirectory checks only declared ordered captures and all hashed
// JSON evidence. SQLite is never opened by the independent checker.
func CheckBindingDirectory(dir string) (report CaseReport, err error) {
	report.Files = map[string]string{}
	report.Outcomes = map[string]json.RawMessage{}
	caseBytes, err := os.ReadFile(filepath.Join(dir, "case.json"))
	if err != nil {
		return report, err
	}
	var entry BindingCase
	if err = json.Unmarshal(caseBytes, &entry); err != nil {
		return report, err
	}
	if entry.Version != "alan/phase3b4-case/v1" {
		return report, fmt.Errorf("unknown binding case version %q", entry.Version)
	}
	report.Name, report.ClaimedClassification, report.OriginalDirectory, report.CaseSHA256 = entry.Name, entry.Classification, dir, sha(caseBytes)
	hashes, err := os.ReadFile(filepath.Join(dir, "digests.json"))
	if err != nil {
		return report, err
	}
	var digests map[string]string
	if err = strictJSON(hashes, &digests); err != nil {
		return report, err
	}
	files := map[string][]byte{}
	paths := []string{}
	for p := range digests {
		paths = append(paths, p)
	}
	sort.Strings(paths)
	for _, p := range paths {
		full, e := resolve(dir, dir, filepath.Join(dir, p))
		if e != nil {
			return report, e
		}
		data, e := os.ReadFile(full)
		if e != nil {
			return report, e
		}
		if sha(data) != digests[p] {
			return report, fmt.Errorf("digest mismatch %s", p)
		}
		report.Files[p] = sha(data)
		if strings.HasSuffix(p, ".json") || strings.HasSuffix(p, ".jsonl") {
			files[p] = data
		}
	}
	files["case.json"] = caseBytes
	for _, capture := range entry.Captures {
		p := capture.Directory + "/inspect.json"
		data, ok := files[p]
		if !ok {
			return report, fmt.Errorf("missing declared capture %s", p)
		}
		v := CheckRaw(data)
		if entry.PrincipalRequests[0] != "" {
			v = CheckBindingRaw(data, entry.Spec())
		}
		projection := projectChecked(data, v)
		if len(v.Custody.Auxiliary) > 0 {
			projection.Supported = false
			projection.Reasons = append(projection.Reasons, "auxiliary_originals_outside_restricted_correspondence")
		}
		for _, c := range entry.Captures {
			if c.AcknowledgedSIGKILL {
				projection.Supported = false
				projection.Reasons = append(projection.Reasons, "whole_store_crash_paired_expansion_not_established")
				break
			}
		}
		report.Phases = append(report.Phases, PhaseReport{Path: p, SHA256: sha(data), Verdict: v, Projection: projection})
		report.Final = v
		report.Projection = projection
	}
	if len(entry.Captures) == 0 {
		report.Projection.Reasons = []string{"no_store_created_input_validation_obstruction"}
		return report, nil
	}
	last := entry.Captures[len(entry.Captures)-1].Directory + "/inspect.json"
	report.External = CheckExternal(files[last], files)
	report.ContinuityFailure = CheckContinuity(report.Phases, files)
	return report, nil
}

func checkBindingLosses(c BindingCase, files map[string][]byte, out *ExternalReport) {
	for _, capture := range c.Captures {
		if !capture.AcknowledgedSIGKILL {
			continue
		}
		matched := false
		for name, b := range files {
			if name != "kill.json" && !(strings.HasPrefix(name, "loss-") && strings.HasSuffix(name, ".json")) {
				continue
			}
			var loss struct {
				PID            int
				Wait           string
				Acknowledgment boundaryEvent
				Command        []string
			}
			need(strictJSON(b, &loss) == nil, "loss_encoding", name, "")
			if loss.PID != capture.OwnerPID {
				continue
			}
			need(loss.Wait == "signal: killed", "loss_acknowledgment", name, loss.Wait)
			ack := loss.Acknowledgment
			if ack.PID == 0 {
				var acknowledgments map[string]boundaryEvent
				need(strictJSON(files["acknowledgments.json"], &acknowledgments) == nil, "loss_boundary", name, "")
				ack = acknowledgments["Runtime"]
			}
			need(ack.PID == capture.OwnerPID && ack.Boundary != "" && ack.Token != "", "loss_boundary", name, "")
			if ack.Subject != "" {
				var raw RawSnapshot
				need(strictJSON(files[capture.Directory+"/inspect.json"], &raw) == nil, "loss_snapshot", name, "")
				found := false
				wantSequence := ack.Sequence
				if ack.Boundary == "execution.staged" && ack.Phase == "pre-commit" {
					wantSequence--
				}
				for _, r := range raw.Tables["runs"] {
					if textCell(r, "subject") == ack.Subject {
						found = intCell(r, "generation") == ack.Generation && intCell(r, "epoch") == ack.Epoch && intCell(r, "sequence") == wantSequence
					}
				}
				need(found, "loss_boundary_snapshot", name, ack.Subject)
			}
			matched = true
		}
		need(matched, "missing_acknowledged_loss", capture.Directory, "")
		out.AcknowledgedLosses++
	}
}
