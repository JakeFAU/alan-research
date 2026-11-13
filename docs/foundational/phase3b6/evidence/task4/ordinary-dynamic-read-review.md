# Unchanged ordinary group: dynamic input check

The58 selected packages have no Go test dependency on effectobservation/delegationcorrespondence. A focused scan of their test sources found external-module builds in author and examples/evidence/investigation; their fixture imports and commands use unchanged public/runtime packages. Durable-test child builds target cmd/alan, cmd/alan-resource and internal/executor, also unchanged. main_integration_test.go is integration-tagged and outside this ordinary invocation.

internal/durabletest/phase2_process_test.go:68-136 walks Go/module source files to emit provenance metadata. It may include excluded in-progress research source hashes. The metadata is logged/archived by phaseEmit, not used to execute or assert semantics. This ordinary suite is test qualification under the528-input unchanged-source snapshot, not a new frozen-source runtime evidence campaign; do not promote its transient provenance records into new Phase3B6 capture evidence. The new relation uses only the sealed earlier producer records. The scan is a focused source check, not a proof of arbitrary filesystem noninterference.

Final qualification still must compare all528 snapshot inputs against final files and verify disjoint ordinary package groups cover the final go list ./... inventory. Any input drift requires impacted requalification.
