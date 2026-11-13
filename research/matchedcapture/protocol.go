package matchedcapture

import (
	"encoding/json"
	"net/url"
	"path/filepath"
	"strings"

	dc "github.com/JakeFAU/alan-research/research/delegationcorrespondence"
)

type branchProtocol struct {
	Owner, Recovery event
	Endpoints       map[string]ready
	Processes       map[string]process
}
type transcript struct {
	Rows  []record
	Index int
	Path  string
}

func (t *transcript) next(action, role string, pid int) json.RawMessage {
	require(t.Index < len(t.Rows), t.Path, "missing "+action)
	r := t.Rows[t.Index]
	t.Index++
	require(r.Ordinal == t.Index && r.Action == action && r.Role == role && (pid < 0 || r.PID == pid), t.Path, "unexpected action/ordinal/role/PID at "+action)
	return r.Data
}
func checkProtocol(f map[string][]byte, p pair, branch string) branchProtocol {
	path := branch + "/protocol.jsonl"
	t := transcript{Rows: lines[record](f, path), Path: path}
	out := branchProtocol{Endpoints: map[string]ready{}, Processes: map[string]process{}}
	require(p.ControllerPID > 0, path, "controller PID")
	source := read[source](f, "source.json")
	require(len(source.TestArgv) > 0, path, "missing source executable")
	q := selectedRequest(f, branch)
	start := func(role string) process {
		d := decode[struct {
			Config  event
			Command []string
		}](t.next("start-request", role, 0), path)
		c := d.Config
		require(c.PID == 0 && c.Scenario == "phase3b7/"+branch+"/"+role && c.Token != "" && filepath.IsAbs(c.Store) && filepath.Base(filepath.Dir(c.Store)) == branch, path, "start identity")
		require(c.Generation == 0 && c.Epoch == 0 && c.Intent == "" && c.Attempt == "" && c.AuthorityAdmissionOrdinal == 0 && c.Root == "" && len(c.InFlight) == 0, path, "unexpected start selectors")
		switch role {
		case "owner":
			require(c.Boundary == "runner.waiting" && c.Phase == "in-flight" && c.Subject == "" && c.Sequence == 0 && filepath.Base(c.Store) == "runtime.sqlite" && same(d.Command, []string{source.TestArgv[0], "-test.run=^TestPhase3B7Owner$", "-test.v"}), path, "owner configuration/command")
		case "endpoint-a":
			require(c.Boundary == "" && c.Phase == "" && c.Subject == "" && c.Sequence == 0, path, "endpoint A selectors")
		case "endpoint-b":
			require(c.Boundary == "resource.applied" && c.Phase == "post-commit" && c.Subject == "a" && c.Sequence == 4, path, "endpoint B selectors")
		}
		if role != "owner" {
			require(filepath.Base(c.Store) == role+".sqlite" && len(d.Command) == 9 && filepath.Base(d.Command[0]) == "alan-resource" && d.Command[1] == "--db" && d.Command[2] == c.Store && d.Command[3] == "--token" && d.Command[4] == c.Token && d.Command[5] == "--listen" && d.Command[7] == "--journal" && filepath.IsAbs(d.Command[8]) && strings.HasSuffix(d.Command[8], "/"+branch+"/"+role+"/journal.jsonl"), path, "endpoint command/store/journal")
		}
		require(t.Index < len(t.Rows), path, "missing start result")
		pid := t.Rows[t.Index].PID
		require(pid > 0 && pid != p.ControllerPID, path, "child PID")
		result := decode[struct{ Error string }](t.next("start-result", role, pid), path)
		require(result.Error == "<nil>", path, "start failed")
		for _, other := range out.Processes {
			require(other.PID != pid && other.Config.Token != c.Token && other.Config.Store != c.Store, path, "aliased live process")
		}
		proc := process{c, pid, d.Command}
		out.Processes[role] = proc
		_, ok := f[branch+"/"+role+"/process.log"]
		if role == "owner" {
			_, ok = f[branch+"/first-owner/process.log"]
		}
		require(ok, path, "missing process log")
		return proc
	}
	for _, role := range []string{"endpoint-a", "endpoint-b"} {
		proc := start(role)
		d := decode[struct {
			Record ready
			Error  string
		}](t.next("ready", role, proc.PID), path)
		r := read[ready](f, branch+"/"+role+"/ready.json")
		require(d.Error == "<nil>" && same(r, d.Record) && r.PID == proc.PID && r.Token == proc.Config.Token && r.Store == proc.Config.Store, path, "ready/config identity")
		u, e := url.Parse(r.URL)
		require(e == nil && u.Scheme == "http" && u.Hostname() == "127.0.0.1" && u.Port() != "" && u.Port() != "0" && u.Path == "/apply" && u.RawQuery == "" && u.Fragment == "" && u.User == nil, path, "logical endpoint URL")
		listen := "127.0.0.1:0"
		if branch == "applied" {
			listen = u.Host
		}
		require(proc.Command[6] == listen, path, "listener/rebinding")
		out.Endpoints[role] = r
	}
	require(out.Endpoints["endpoint-a"].URL != out.Endpoints["endpoint-b"].URL, path, "aliased URLs")
	if branch == "zero" {
		d := decode[struct{ File string }](t.next("seed-copies", "controller", p.ControllerPID), path)
		require(d.File == "seed/copies.json", path, "seed copies marker")
	}
	owner := start("owner")
	out.Owner = owner.Config
	out.Owner.PID = owner.PID
	ack := func(action, role string, pid int) event {
		d := decode[struct {
			Event event
			Error string
		}](t.next(action, role, pid), path)
		require(d.Error == "<nil>", path, action+" failed")
		proc := out.Processes[role]
		identity(d.Event, proc.Config, proc.PID, path)
		return d.Event
	}
	worker := ack("pre-call-ack", "owner", owner.PID)
	require(worker.Boundary == "dispatch.admitted" && worker.Phase == "pre-call" && worker.Intent == service && worker.Subject == q.Subject && worker.Generation == q.Generation && worker.Epoch == q.Epoch && worker.Sequence == q.Sequence && worker.Attempt == q.Authorization && worker.AuthorityAdmissionOrdinal == 1 && worker.Root == "" && len(worker.InFlight) == 0, path, "pre-call selected binding")
	require(same(worker, read[event](f, branch+"/first-owner/pre-call.json")), path, "pre-call sidecar join")
	scheduler := ack("scheduler-ack", "owner", owner.PID)
	require(scheduler.Boundary == "runner.waiting" && scheduler.Phase == "in-flight" && scheduler.Subject == "" && scheduler.Generation == 0 && scheduler.Epoch == 0 && scheduler.Sequence == 0 && scheduler.Intent == "" && scheduler.Attempt == "" && scheduler.Root == "" && scheduler.AuthorityAdmissionOrdinal == 0 && same(scheduler.InFlight, []string{q.Authorization}), path, "scheduler sole selected authorization")
	release := func(role, target string, pid int) {
		a := decode[struct{ Target, Command string }](t.next("release-request", role, pid), path)
		b := decode[struct {
			Target string
			Bytes  int
			Error  string
		}](t.next("release-result", role, pid), path)
		require(a.Target == target && a.Command == "c" && b.Target == target && b.Bytes == 1 && b.Error == "<nil>", path, "release transaction")
	}
	if branch == "applied" {
		release("owner", "worker", owner.PID)
		require(same(worker, read[event](f, branch+"/first-owner/worker-release.json")), path, "worker release sidecar")
		a := ack("application-ack", "endpoint-b", out.Processes["endpoint-b"].PID)
		require(a.Boundary == "resource.applied" && a.Phase == "post-commit" && a.Subject == q.Subject && a.Generation == q.Generation && a.Epoch == q.Epoch && a.Sequence == q.Sequence && a.Intent == q.Key && a.Attempt == q.Authorization && a.Root == "" && a.AuthorityAdmissionOrdinal == 0 && len(a.InFlight) == 0, path, "application acknowledgment wire")
	} else {
		_, ok := f[branch+"/first-owner/worker-release.json"]
		require(!ok, path, "zero released worker")
	}
	stop := func(role, purpose string, pid int) {
		a := decode[struct{ Purpose, Signal string }](t.next("kill-request", role, pid), path)
		b := decode[struct {
			Purpose       string
			KillError     string `json:"kill_error"`
			Error, Signal string
			ExitCode      int `json:"exit_code"`
			Status        string
		}](t.next("wait", role, pid), path)
		require(a.Purpose == purpose && a.Signal == "killed" && b.Purpose == purpose && b.KillError == "<nil>" && b.Error == "signal: killed" && b.Signal == "killed" && b.ExitCode == -1 && b.Status == "signal: killed", path, "confirmed SIGKILL/reap")
	}
	capture := func(cut string) {
		d := decode[struct{ Cut string }](t.next("capture", "controller", p.ControllerPID), path)
		require(d.Cut == cut, path, "capture cut")
	}
	snap := func(role, cut string) {
		d := decode[struct{ Cut, Error string }](t.next("endpoint-snapshot", role, out.Processes[role].PID), path)
		require(d.Cut == cut && d.Error == "<nil>", path, "snapshot cut")
	}
	stop("owner", "owner-loss", owner.PID)
	capture("stopped")
	snap("endpoint-a", "stopped")
	snap("endpoint-b", "stopped")
	if branch == "applied" {
		release("endpoint-b", "endpoint-response", out.Processes["endpoint-b"].PID)
	}
	reopen := decode[struct {
		Identity event
		Error    string
	}](t.next("recovery-open", "controller", p.ControllerPID), path)
	r := reopen.Identity
	require(reopen.Error == "<nil>" && r.PID == p.ControllerPID && r.Scenario == "phase3b7/"+branch+"/recovery" && r.Token != "" && r.Token != owner.Config.Token && r.Store == owner.Config.Store && r.Boundary == "" && r.Phase == "" && r.Subject == "" && r.Generation == 0 && r.Epoch == 0 && r.Sequence == 0 && r.Intent == "" && r.Attempt == "" && r.AuthorityAdmissionOrdinal == 0 && r.Root == "" && len(r.InFlight) == 0, path, "recovery identity")
	out.Recovery = r
	capture("reopened")
	rev := decode[struct{ File string }](t.next("revocation", "controller", p.ControllerPID), path)
	require(rev.File == "controlled-revocation.json", path, "revocation sidecar")
	close := decode[struct {
		RunError   string `json:"run_error"`
		CloseError string `json:"close_error"`
	}](t.next("recovery-close", "controller", p.ControllerPID), path)
	require(close.RunError == "<nil>" && close.CloseError == "<nil>", path, "recovery close")
	for _, role := range []string{"endpoint-a", "endpoint-b"} {
		stop(role, "endpoint-stop", out.Processes[role].PID)
		snap(role, "final")
	}
	capture("final")
	require(string(t.next("branch-complete", "controller", p.ControllerPID)) == "null", path, "complete payload")
	require(t.Index == len(t.Rows), path, "unexpected trailing actions")
	_, abort := f[branch+"/first-owner/abort.json"]
	require(!abort, path, "owner barrier returned")
	checkBoundaries(f, branch, owner, worker, q)
	checkOpen(f, branch+"/first-owner/open.json", branch+"/first-owner/post-open", out.Owner)
	checkOpen(f, branch+"/recovery-open.json", branch+"/reopened", out.Recovery)
	return out
}
func identity(e, config event, pid int, path string) {
	require(e.PID == pid && e.Scenario == config.Scenario && e.Token == config.Token && e.Store == config.Store, path, "process event identity")
}
func checkBoundaries(f map[string][]byte, branch string, owner process, worker event, q dc.WireRequest) {
	p := branch + "/first-owner/boundaries.jsonl"
	events := lines[json.RawMessage](f, p)
	pre, wait, authorized := 0, 0, 0
	for _, rawEvent := range events {
		e := decode[event](rawEvent, p)
		identity(e, owner.Config, owner.PID, p)
		switch e.Boundary {
		case "dispatch.admitted":
			require(e.Phase == "pre-call" && same(e, worker) && authorized == 1, p, "unexpected worker or pre-call before authorization")
			pre++
		case "external.authorized":
			_ = decodeComplete[event](rawEvent, p)
			// The passive native event is logged before the inherited observer
			// enriches its own scheduler-ack copy with InFlight.
			require(e.Intent == service && e.Attempt == q.Authorization &&
				e.Subject == q.Subject && e.Generation == q.Generation &&
				e.Epoch == q.Epoch && e.Sequence == q.Sequence &&
				e.Phase == "pre-transport" && e.Root == "" &&
				e.AuthorityAdmissionOrdinal == 0 && e.InFlight == nil,
				p, "selected native authorization coordinates/selectors")
			authorized++
		case "runner.waiting":
			if e.Phase == "in-flight" {
				wait++
			}
		case "external.finished":
			require(e.Intent != service, p, "selected worker returned before death")
		}
	}
	require(pre == 1 && wait == 1 && authorized == 1, p, "complete selected boundary coverage")
}
func checkIsolation(a, b branchProtocol) {
	for _, role := range []string{"endpoint-a", "endpoint-b", "owner"} {
		first := a.Processes[role]
		second := b.Processes[role]
		require(first.Config.Token != second.Config.Token && first.Config.Store != second.Config.Store, "protocol", "cross-branch process/store identity")
	}
	require(a.Recovery.Token != b.Recovery.Token, "protocol", "cross-branch recovery token")
	for _, role := range []string{"endpoint-a", "endpoint-b"} {
		x := a.Endpoints[role]
		require(x.URL == b.Endpoints[role].URL, "protocol", "logical URL drift")
	}
}
