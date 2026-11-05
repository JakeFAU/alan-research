package delegationcorrespondence

import (
	"fmt"
	"sort"
	"strings"
)

type tableShape struct{ Columns, Key string }

var rawShapes = map[string]tableShape{
	"adapter_attempts":                {"descriptor,intent_id", "intent_id"},
	"adapter_results":                 {"descriptor,intent_id", "intent_id"},
	"admissions":                      {"descriptor,domain,id,observation_id,result,transport_id", "id"},
	"attempts":                        {"descriptor,id,intent_id,result", "id"},
	"composition_admission_decisions": {"", ""},
	"composition_bindings":            {"", ""},
	"composition_charges":             {"decision_order,descriptor,kind,operation_id,root_id", "operation_id,kind"},
	"composition_closures":            {"decision_order,descriptor,root_id,scope_id", "scope_id"},
	"composition_deployment":          {"descriptor,singleton", "singleton"},
	"composition_facts":               {"decision_order,descriptor,id,kind", "id"},
	"composition_freezes":             {"descriptor,root_id,scope_id", "scope_id"},
	"composition_invocations":         {"", ""},
	"composition_members":             {"", ""},
	"composition_outcome_deliveries":  {"", ""},
	"composition_permissions":         {"descriptor,singleton", "singleton"},
	"composition_profile":             {"descriptor,singleton,version", "singleton"},
	"composition_publications":        {"", ""},
	"composition_replies":             {"", ""},
	"composition_requests":            {"", ""},
	"composition_reservations":        {"descriptor,operation_id,root_id", "operation_id"},
	"composition_results":             {"", ""},
	"composition_root_publications":   {"decision_order,descriptor,root_id,scope_id", "scope_id"},
	"composition_roots":               {"descriptor,id", "id"},
	"composition_static_members":      {"descriptor,generation,root_id,subject", "subject,generation"},
	"composition_templates":           {"", ""},
	"corpus_assets":                   {"", ""},
	"external_attempts":               {"descriptor,intent_id", "intent_id"},
	"external_receipts":               {"descriptor,intent_id", "intent_id"},
	"history":                         {"epoch,generation,intents,selected_observation,sequence,source_revision,subject,successor", "subject,generation,sequence"},
	"intents":                         {"descriptor,generation,id,ordinal,sequence,source_revision,subject", "id"},
	"lifecycle_results":               {"decision_sequence,descriptor,id,result", "id"},
	"message_acknowledgments":         {"admission_id,intent_id", "intent_id"},
	"metadata":                        {"key_version,namespace,schema_version,semantics,singleton", "singleton"},
	"observations":                    {"consumed_sequence,descriptor,generation,id,subject", "id"},
	"permissions":                     {"kind,obligation,operation,revision,source,target", "kind,source,target,operation,revision,obligation"},
	"programs":                        {"bytes,id,semantics", "id"},
	"run_origins":                     {"descriptor,generation,kind,operation_id,predecessor_generation,subject", "subject,generation"},
	"runs":                            {"epoch,generation,initial,program_id,projection,sequence,subject", "subject,generation"},
	"settlements":                     {"descriptor,id,intent_id,observation_id", "id"},
	"subjects":                        {"generation,id,owner,retired,retry_cancelled,revision", "id"},
	"terminal_obligations":            {"delivery_status,destinations,generation,operation_id,outcome,sequence,source_revision,subject", "subject,generation"},
}

func rowKey(t string, r Row) string {
	parts := []any{}
	for _, k := range strings.Split(rawShapes[t].Key, ",") {
		parts = append(parts, r[k])
	}
	return string(bytesJSON(parts))
}
func checkSchema(raw RawSnapshot) {
	sparse := raw.Namespace == "" && raw.KeyVersion == 0 && len(raw.Settings) == 0 && len(raw.Roots) == 0
	if !sparse {
		need(len(raw.Tables) == len(rawShapes), "raw_table_set", "", fmt.Sprint(len(raw.Tables)))
	}
	for table := range raw.Tables {
		_, known := rawShapes[table]
		need(known, "unknown_table", table, "")
	}
	names := []string{}
	for t := range rawShapes {
		names = append(names, t)
	}
	sort.Strings(names)
	for _, t := range names {
		shape := rawShapes[t]
		rows, ok := raw.Tables[t]
		need(ok || sparse, "missing_table", t, "")
		if shape.Columns == "" {
			need(len(rows) == 0, "outside_static_domain", t, "")
			continue
		}
		want := strings.Split(shape.Columns, ",")
		seen := map[string]bool{}
		for _, r := range rows {
			need(len(r) == len(want), "row_columns", t, "")
			for _, k := range want {
				_, ok := r[k]
				need(ok, "row_columns", t, k)
			}
			key := rowKey(t, r)
			need(!seen[key], "duplicate_primary_key", t, key)
			seen[key] = true
		}
	}
}
func (c *checker) checkTerminals() {
	seen := map[RunID]bool{}
	for _, row := range c.raw.Tables["terminal_obligations"] {
		run := RunID{textCell(row, "subject"), intCell(row, "generation")}
		seq := intCell(row, "sequence")
		h, ok := c.history[historyKey(run, seq)]
		p := decodeCell[ProjectionState](row, "outcome")
		if op := nullableText(row, "operation_id"); op != "" {
			a := c.life[op]
			need(a.Accepted && a.After != nil && a.After.Run == run && a.After.Sequence == seq, "terminal_lifecycle", op, "")
			h.Next = a.After.Projection
		}
		need(ok && h.Next.Status != "open" && equal(h.Next, p) && intCell(row, "source_revision") == h.Revision, "terminal_original", run.Subject, "")
		seen[run] = true
		var destinations []struct{ Target, Obligation string }
		need(canonical(rawCell(row, "destinations"), &destinations) == nil, "terminal_destinations", run.Subject, "")
		var status map[string]bool
		need(canonical(rawCell(row, "delivery_status"), &status) == nil, "terminal_delivery", run.Subject, "")
		expectedRoutes := map[string]bool{}
		for _, p := range c.permissions {
			if p.Kind == "observation" && p.Source == run.Subject && p.Operation == "terminal" && p.Revision == h.Revision {
				expectedRoutes[identity("route", p.Target, p.Obligation)] = true
			}
		}
		need(len(destinations) == len(expectedRoutes), "terminal_route_completeness", run.Subject, "")
		seenRoutes := map[string]bool{}
		for key := range status {
			need(expectedRoutes[key], "terminal_delivery_route", key, "")
		}
		for _, d := range destinations {
			key := identity("route", d.Target, d.Obligation)
			need(expectedRoutes[key] && !seenRoutes[key], "terminal_route_completeness", key, "")
			seenRoutes[key] = true
			allowed := false
			for _, p := range c.permissions {
				allowed = allowed || p == (Permission{"observation", run.Subject, d.Target, "terminal", h.Revision, d.Obligation})
			}
			need(allowed, "terminal_permission", run.Subject, d.Target)
			if status[identity("route", d.Target, d.Obligation)] {
				id := identity("admission", "terminal", identity("terminal-delivery", run.Subject, run.Generation, seq, d.Target, d.Obligation))
				need(c.admitResults[id].Accepted, "terminal_delivery", id, "")
			}
		}
	}
	for run, r := range c.runs {
		need((r.Projection.Status != "open") == seen[run], "terminal_obligation_presence", run.Subject, "")
	}
}
