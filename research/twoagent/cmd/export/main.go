// Command export writes deterministic evidence for the frozen two-agent model.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/JakeFAU/alan-research/research/twoagent"
)

func main() {
	assertNaive := flag.Bool("assert-naive-recovery", false, "exit 1 when local safety and deduplication fail to imply useful recovery (expected falsification)")
	flag.Parse()
	reports := []twoagent.Report{twoagent.Explore(twoagent.DeliveryACK), twoagent.Explore(twoagent.DurableACK)}
	for _, r := range reports {
		for _, w := range r.LossWitnesses {
			if err := twoagent.CheckTrace(r.Profile, w.Trace); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(2)
			}
		}
	}
	output := struct {
		Contract string
		Reports  []twoagent.Report
	}{Contract: "phase3a-two-agent-custody-v1", Reports: reports}
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(output); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	locallySafe := reports[0].LocalSafetyViolations == 0 && reports[0].LocalTransitionViolations == 0
	if *assertNaive && locallySafe && len(reports[0].LossWitnesses) > 0 {
		fmt.Fprintln(os.Stderr, "FALSIFIED: local safety plus original-request deduplication does not imply useful recovery; see independently checked recovered-loss witnesses")
		os.Exit(1)
	}
}
