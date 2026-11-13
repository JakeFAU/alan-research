package matchedcapture

import (
	"bytes"
	"compress/gzip"
	_ "embed"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"testing"

	dc "github.com/JakeFAU/alan-research/research/delegationcorrespondence"
)

//go:embed testdata/historical-json.json.gz
var historicalGzip []byte

func historical(t *testing.T) map[string]json.RawMessage {
	t.Helper()
	z, e := gzip.NewReader(bytes.NewReader(historicalGzip))
	if e != nil {
		t.Fatal(e)
	}
	b, e := io.ReadAll(z)
	if e != nil {
		t.Fatal(e)
	}
	z.Close()
	var m map[string]json.RawMessage
	if e = json.Unmarshal(b, &m); e != nil {
		t.Fatal(e)
	}
	return m
}
func fixture(t *testing.T) map[string][]byte {
	t.Helper()
	h := historical(t)
	f := map[string][]byte{}
	put := func(p string, v any) { f[p] = append(jsonBytes(v), '\n') }
	put("pair.json", pair{Version, "complete", "normal", 100, []string{"zero", "applied"}})
	env := map[string]string{"GOTOOLCHAIN": "local", "GOWORK": "off", "GOPROXY": "off", "GOSUMDB": "off", "GOMODCACHE": "synthetic", "GOCACHE": "synthetic", "GOFLAGS": "", "GOENV": "", "CGO_ENABLED": ""}
	put("source.json", source{Version: Version, Prepared: preparedBytes(), GoVersion: "go1.27.4", GOOS: "darwin", GOARCH: "arm64", Environment: env, TestArgv: []string{"/synthetic/owner.test", "-test.run=^TestPhase3B7MatchedCapture$"}, BuildInfo: debug.BuildInfo{GoVersion: "go1.27.4", Path: "github.com/JakeFAU/alan/internal/durabletest.test", Deps: []*debug.Module{{Path: "github.com/mattn/go-sqlite3", Version: "v1.14.52"}}, Settings: []debug.BuildSetting{{Key: "GOOS", Value: "darwin"}, {Key: "GOARCH", Value: "arm64"}, {Key: "CGO_ENABLED", Value: "1"}, {Key: "-tags", Value: "integration"}}}, OwnerBinary: digest([]byte("synthetic owner")), EndpointBinary: digest([]byte("synthetic endpoint"))})
	var oldManifest struct{ Reciprocal json.RawMessage }
	json.Unmarshal(h["manifest.json"], &oldManifest)
	f["seed/manifest.json"] = oldManifest.Reciprocal
	f["seed/provision.json"] = h["provision.json"]
	capture := func(prefix, old string) {
		for name, b := range h {
			if strings.HasPrefix(name, old+"/") {
				f[prefix+strings.TrimPrefix(name, old)] = append([]byte{}, b...)
			}
		}
		f[prefix+"/runtime.sqlite"] = []byte("SYNTHETIC OPAQUE BYTES, NOT SQLITE")
	}
	capture("seed/initial", "initial")
	seed := map[string]Digest{"runtime.sqlite": digest(f["seed/initial/runtime.sqlite"])}
	put("seed/copies.json", map[string]any{"seed": seed, "zero": seed, "applied": seed})
	var raw dc.RawSnapshot
	json.Unmarshal(h["stopped/inspect.json"], &raw)
	var attempt dc.ExternalAttempt
	json.Unmarshal([]byte(decode[string](raw.Tables["external_attempts"][0]["descriptor"], "fixture")), &attempt)
	var endpoints map[string]json.RawMessage
	json.Unmarshal(h["endpoints.json"], &endpoints)
	for bi, branch := range []string{"zero", "applied"} {
		base := branch + "/"
		for _, cut := range []string{"first-owner/post-open", "stopped", "reopened", "final"} {
			capture(base+cut, cut)
		}
		f[base+"controlled-revocation.json"] = h["controlled-revocation.json"]
		f[base+"run.json"] = h["run.json"]
		protocols := []record{}
		add := func(action, role string, pid int, data any) {
			protocols = append(protocols, record{len(protocols) + 1, action, role, pid, jsonBytes(data)})
		}
		procs := map[string]event{}
		readies := map[string]ready{}
		for ei, label := range []string{"a", "b"} {
			role := "endpoint-" + label
			c := event{Scenario: "phase3b7/" + branch + "/" + role, Token: branch + "-" + role, Store: "/synthetic/" + branch + "/" + role + ".sqlite"}
			if label == "b" {
				c.Boundary = "resource.applied"
				c.Phase = "post-commit"
				c.Subject = "a"
				c.Sequence = 4
			}
			pid := 200 + bi*10 + ei
			r := decode[ready](h[role+"-ready.json"], "fixture")
			r.PID = pid
			r.Token = c.Token
			r.Store = c.Store
			readies[role] = r
			procs[role] = c
			add("start-request", role, 0, map[string]any{"config": c, "command": []string{"/synthetic/alan-resource", "--db", c.Store, "--token", c.Token, "--listen", "127.0.0.1:0", "--journal", "/synthetic/raw/" + branch + "/" + role + "/journal.jsonl"}})
			if bi == 1 {
				var d map[string]any
				json.Unmarshal(protocols[len(protocols)-1].Data, &d)
				d["command"].([]any)[6] = strings.TrimSuffix(strings.TrimPrefix(r.URL, "http://"), "/apply")
				protocols[len(protocols)-1].Data = jsonBytes(d)
			}
			add("start-result", role, pid, map[string]any{"error": "<nil>"})
			add("ready", role, pid, map[string]any{"record": r, "error": "<nil>"})
			put(base+role+"/ready.json", r)
			f[base+role+"/process.log"] = []byte{}
			f[base+role+"/resource.sqlite"] = []byte("SYNTHETIC ENDPOINT " + branch + label)
			var snapshot map[string]any
			json.Unmarshal(endpoints["a"], &snapshot)
			f[base+role+"/journal.jsonl"] = []byte{}
			if branch == "applied" && label == "b" {
				json.Unmarshal(endpoints["b"], &snapshot)
				receipt := snapshot["Applications"].([]any)[0].(map[string]any)
				receipt["PID"] = pid
				receipt["Process"] = c.Token
				records := []map[string]any{{"Event": "ingress", "PID": pid, "Process": c.Token, "Contact": 1, "Method": "POST", "Path": "/apply", "Body": string(jsonBytes(attempt.Request))}, {"Event": "committed", "PID": pid, "Process": c.Token, "Contact": 1, "Receipt": receipt}}
				for _, v := range records {
					f[base+role+"/journal.jsonl"] = append(f[base+role+"/journal.jsonl"], append(jsonBytes(v), '\n')...)
				}
			}
			for _, cut := range []string{"stopped", "final"} {
				put(base+role+"/"+cut+".json", map[string]any{"ready": r, "snapshot": snapshot, "error": "<nil>"})
			}
		}
		if branch == "zero" {
			add("seed-copies", "controller", 100, map[string]any{"file": "seed/copies.json"})
		}
		owner := event{Scenario: "phase3b7/" + branch + "/owner", Token: branch + "-owner", Store: "/synthetic/" + branch + "/runtime.sqlite", Boundary: "runner.waiting", Phase: "in-flight"}
		pid := 202 + bi*10
		add("start-request", "owner", 0, map[string]any{"config": owner, "command": []string{"/synthetic/owner.test", "-test.run=^TestPhase3B7Owner$", "-test.v"}})
		add("start-result", "owner", pid, map[string]any{"error": "<nil>"})
		owner.PID = pid
		worker := owner
		worker.Boundary = "dispatch.admitted"
		worker.Phase = "pre-call"
		worker.Subject = "a"
		worker.Generation = 1
		worker.Epoch = 2
		worker.Sequence = 4
		worker.AuthorityAdmissionOrdinal = 1
		worker.Intent = service
		worker.Attempt = attempt.AuthorizationID
		scheduler := owner
		scheduler.InFlight = []string{attempt.AuthorizationID}
		add("pre-call-ack", "owner", pid, map[string]any{"event": worker, "error": "<nil>"})
		add("scheduler-ack", "owner", pid, map[string]any{"event": scheduler, "error": "<nil>"})
		put(base+"first-owner/pre-call.json", worker)
		authorized := worker
		authorized.Boundary = "external.authorized"
		authorized.Phase = "pre-transport"
		authorized.AuthorityAdmissionOrdinal = worker.AuthorityAdmissionOrdinal
		passive := owner
		f[base+"first-owner/boundaries.jsonl"] = append(append(append(jsonBytes(authorized), '\n'), append(jsonBytes(worker), '\n')...), append(jsonBytes(passive), '\n')...)
		f[base+"first-owner/process.log"] = []byte{}
		release := func(role, target string, p int) {
			add("release-request", role, p, map[string]any{"target": target, "command": "c"})
			add("release-result", role, p, map[string]any{"target": target, "bytes": 1, "error": "<nil>"})
		}
		if branch == "applied" {
			release("owner", "worker", pid)
			put(base+"first-owner/worker-release.json", worker)
			ack := procs["endpoint-b"]
			ack.PID = readies["endpoint-b"].PID
			ack.Generation = 1
			ack.Epoch = 2
			ack.Intent = attempt.Request.Key
			ack.Attempt = attempt.AuthorizationID
			add("application-ack", "endpoint-b", ack.PID, map[string]any{"event": ack, "error": "<nil>"})
		}
		kill := func(role, purpose string, p int) {
			add("kill-request", role, p, map[string]any{"purpose": purpose, "signal": "killed"})
			add("wait", role, p, map[string]any{"purpose": purpose, "kill_error": "<nil>", "error": "signal: killed", "signal": "killed", "exit_code": -1, "status": "signal: killed"})
		}
		kill("owner", "owner-loss", pid)
		add("capture", "controller", 100, map[string]any{"cut": "stopped"})
		for _, role := range []string{"endpoint-a", "endpoint-b"} {
			add("endpoint-snapshot", role, readies[role].PID, map[string]any{"cut": "stopped", "error": "<nil>"})
		}
		if branch == "applied" {
			release("endpoint-b", "endpoint-response", readies["endpoint-b"].PID)
		}
		recovery := event{PID: 100, Scenario: "phase3b7/" + branch + "/recovery", Token: branch + "-recovery", Store: owner.Store}
		add("recovery-open", "controller", 100, map[string]any{"identity": recovery, "error": "<nil>"})
		add("capture", "controller", 100, map[string]any{"cut": "reopened"})
		add("revocation", "controller", 100, map[string]any{"file": "controlled-revocation.json"})
		add("recovery-close", "controller", 100, map[string]any{"run_error": "<nil>", "close_error": "<nil>"})
		for _, role := range []string{"endpoint-a", "endpoint-b"} {
			kill(role, "endpoint-stop", readies[role].PID)
			add("endpoint-snapshot", role, readies[role].PID, map[string]any{"cut": "final", "error": "<nil>"})
		}
		add("capture", "controller", 100, map[string]any{"cut": "final"})
		add("branch-complete", "controller", 100, nil)
		for _, v := range protocols {
			f[base+"protocol.jsonl"] = append(f[base+"protocol.jsonl"], append(jsonBytes(v), '\n')...)
		}
		for _, o := range []struct {
			path, cut string
			identity  event
		}{{"first-owner/open.json", "first-owner/post-open", owner}, {"recovery-open.json", "reopened", recovery}} {
			var raw dc.RawSnapshot
			json.Unmarshal(f[base+o.cut+"/inspect.json"], &raw)
			runs := map[string]dc.OwnerRun{}
			for _, row := range raw.Tables["runs"] {
				subject := decode[string](row["subject"], "fixture")
				var projection dc.ProjectionState
				json.Unmarshal([]byte(decode[string](row["projection"], "fixture")), &projection)
				runs[subject] = dc.OwnerRun{ID: dc.RunID{Subject: subject, Generation: 1}, ProgramID: decode[string](row["program_id"], "fixture"), Epoch: decode[int64](row["epoch"], "fixture"), Sequence: decode[int64](row["sequence"], "fixture"), Revision: 1, Projection: projection}
			}
			put(base+o.path, opened{o.identity.PID, o.identity.Scenario, o.identity.Token, o.identity.Store, raw.Namespace, "<nil>", runs})
		}
	}
	return f
}
func sealFixture(t *testing.T, f map[string][]byte) string {
	t.Helper()
	d := t.TempDir()
	members := map[string]Digest{}
	for p, b := range f {
		if p == "members.json" {
			continue
		}
		members[p] = digest(b)
		if e := os.MkdirAll(filepath.Dir(filepath.Join(d, p)), 0700); e != nil {
			t.Fatal(e)
		}
		if e := os.WriteFile(filepath.Join(d, p), b, 0600); e != nil {
			t.Fatal(e)
		}
	}
	b := jsonBytes(map[string]any{"version": Version, "files": members})
	if e := os.WriteFile(filepath.Join(d, "members.json"), b, 0600); e != nil {
		t.Fatal(e)
	}
	return d
}
