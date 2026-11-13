package matchedcapture

import (
	"encoding/json"

	dc "github.com/JakeFAU/alan-research/research/delegationcorrespondence"
)

type WorldFacts struct {
	Contacts, Applications int
	SelectedKey            string
	Receipt                *dc.Receipt
}
type endpointSnapshot struct {
	ForeignKeys, BusyTimeoutMS                 int
	SQLiteVersion, SQLiteSourceID, JournalMode string
	Synchronous, QueryOnly                     int
	Requests                                   []dc.WireRequest
	Applications                               []dc.Receipt
	Total, Contacts                            int64
}
type endpointExport struct {
	Ready    ready
	Snapshot endpointSnapshot
	Error    string
}
type journalRecord struct {
	Event              string
	PID                int
	Process            string
	Contact            int
	Method, Path, Body string
	Receipt            *dc.Receipt
}

func checkWorld(f map[string][]byte, d map[string]Digest, branch string) WorldFacts {
	q := selectedRequest(f, branch)
	result := WorldFacts{SelectedKey: q.Key}
	for _, label := range []string{"a", "b"} {
		prefix := branch + "/endpoint-" + label
		r := read[ready](f, prefix+"/ready.json")
		require(r.PID > 0 && r.Token != "" && r.Store != "" && r.URL != "", prefix, "complete ready identity")
		require(d[prefix+"/resource.sqlite"].Bytes > 0, prefix, "missing closed endpoint bytes")
		journals := lines[journalRecord](f, prefix+"/journal.jsonl")
		want := 0
		if branch == "applied" && label == "b" {
			want = 1
		}
		require(len(journals) == 2*want, prefix, "complete journal must have exact ingress/commit cardinality")
		var receipt dc.Receipt
		if want == 1 {
			in, commit := journals[0], journals[1]
			require(in.Event == "ingress" && in.PID == r.PID && in.Process == r.Token && in.Contact == 1 && in.Method == "POST" && in.Path == "/apply" && in.Receipt == nil, prefix, "ingress identity/coverage")
			wire := decode[dc.WireRequest]([]byte(in.Body), prefix+"/journal.jsonl.Body")
			require(wire == q, prefix, "exact selected wire")
			require(commit.Event == "committed" && commit.PID == r.PID && commit.Process == r.Token && commit.Contact == 1 && commit.Method == "" && commit.Path == "" && commit.Body == "" && commit.Receipt != nil, prefix, "commit identity/coverage")
			receipt = *commit.Receipt
			require(receipt == (dc.Receipt{Key: q.Key, Profile: q.Profile, Delta: 1, Mutation: 1, Total: 1, Process: r.Token, PID: r.PID}), prefix, "authentic selected receipt")
			result.Receipt = &receipt
		}
		var first endpointSnapshot
		for _, cut := range []string{"stopped", "final"} {
			p := prefix + "/" + cut + ".json"
			outer := read[map[string]json.RawMessage](f, p)
			require(len(outer) == 3, p, "complete endpoint export")
			e := read[endpointExport](f, p)
			s := e.Snapshot
			require(e.Error == "<nil>" && same(e.Ready, r), p, "snapshot ready identity")
			require(s.ForeignKeys == 1 && s.BusyTimeoutMS == 5000 && s.SQLiteVersion != "" && s.SQLiteSourceID != "" && s.JournalMode == "wal" && s.Synchronous == 2 && s.QueryOnly == 1, p, "snapshot configuration")
			require(s.Contacts == int64(want) && s.Total == int64(want) && len(s.Requests) == want && len(s.Applications) == want, p, "snapshot contact/application coverage")
			if want == 1 {
				require(s.Requests[0] == q && s.Applications[0] == receipt, p, "snapshot journal binding")
			}
			if cut == "stopped" {
				first = s
			} else {
				require(same(first, s), p, "endpoint changed after owner death")
			}
		}
		result.Contacts += want
		result.Applications += want
	}
	return result
}
