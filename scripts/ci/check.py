"""Explicit CI scopes. Full suites require a separate command."""

import argparse
import json
import os
from pathlib import Path
import re
import shlex
import subprocess


ROOT = Path(__file__).resolve().parents[2]
SMOKE_PACKAGES = ["./schema", "./effects", "./compose", "./author", "./research/twoagent"]
STATICCHECK_PACKAGES = ["./schema", "./effects", "./compose", "./research/twoagent"]
RUNTIME_SMOKE = [
    "TestAdapterRunReopenReplay",
    "TestUnknownAdapterNeverRepeats",
    "TestRuntimeAuthoredControlPaths",
]


def run(*args, cwd=ROOT, capture=False):
    print(f"[{Path(cwd).relative_to(ROOT)}] {shlex.join(args)}", flush=True)
    result = subprocess.run(args, cwd=cwd, text=True, capture_output=capture)
    if capture and result.returncode:
        print(result.stdout, end="")
        print(result.stderr, end="")
    result.check_returncode()
    return result.stdout


def tests(packages, timeout, *, names=(), subtests=(), integration=False, race=False, cwd=ROOT):
    flags = ["-count=1", f"-timeout={timeout}"]
    if integration:
        flags += ["-tags", "integration"]
    if race:
        flags += ["-race"]
    if names:
        pattern = "^(" + "|".join(names) + ")$"
        # Go accepts a selector that matches nothing. Make removed/renamed gates fail.
        listed = run("go", "test", *flags, "-list", pattern, *packages, cwd=cwd, capture=True)
        missing = set(names) - set(listed.splitlines())
        if missing:
            raise RuntimeError(f"CI test selector lost tests: {sorted(missing)}")
        pattern += "".join(f"/^{re.escape(name)}$" for name in subtests)
        flags += ["-run", pattern]
    output = run("go", "test", *flags, "-v", *packages, cwd=cwd, capture=bool(subtests))
    if subtests:
        print(output, end="")
        full_name = "/".join([names[0], *subtests])
        if f"--- PASS: {full_name} (" not in output:
            raise RuntimeError(f"CI subtest did not pass: {full_name}")


def pr():
    files = run("git", "ls-files", "-z", "--", "*.go", ":(exclude)docs/**", capture=True)
    unformatted = run("gofmt", "-l", *files.rstrip("\0").split("\0"), capture=True)
    if unformatted:
        raise RuntimeError(f"Run gofmt on:\n{unformatted}")
    run("python3", "-m", "unittest", "discover", "-s", "scripts/ci", "-p", "test_*.py", "-v")
    run("go", "mod", "verify")
    run("go", "build", "./...")
    run("go", "vet", "./...")
    run("go", "vet", "-tags", "integration", "./...")
    run("staticcheck", *STATICCHECK_PACKAGES)
    run("staticcheck", "-tags", "integration", *STATICCHECK_PACKAGES)
    tests(SMOKE_PACKAGES, "120s")
    tests(["./runtime"], "120s", names=RUNTIME_SMOKE, integration=True)


def regression():
    runtime_names = RUNTIME_SMOKE + [
        "TestKnownAdapterResultResumesWithoutRegistry",
        "TestHistoricalReplayAndArtifactEncodingCompatibility",
        "TestRuntimeEnforcesAuthoredSchemasBeforeCommit",
        "TestRuntimeEnforcesUnambiguousSendContracts",
    ]
    for race in (False, True):
        tests(["./runtime"], "5m", names=runtime_names, integration=True, race=race)
        tests(["./internal/executor"], "5m", names=["TestReceiptConflictCrash"], integration=True, race=race)
        tests(
            ["./internal/durabletest"], "5m", integration=True, race=race,
            names=[
                "TestAutomaticCLINestedCommitCrashBoundaries",
                "TestAutomaticNestedAuthorizationDecisionTerminalAndRepeatedCrash",
            ],
        )


def consumer(name):
    cwd = ROOT / "examples" / name / "consumer"
    run("go", "mod", "verify", cwd=cwd)
    run("go", "vet", "./...", cwd=cwd)
    tests(["./..."], "240s" if name == "worktree" else "15m", cwd=cwd)
    if name == "worktree":
        tests(["./..."], "240s", race=True, cwd=cwd)
    else:
        tests(
            ["./..."], "10m", race=True, cwd=cwd,
            names=["TestAdaptivePublicConsumerDirectNested"],
            subtests=["nested=true", "capacity=2", "paired=true"],
        )


def security():
    run("govulncheck", "./...")
    run("govulncheck", "-tags", "integration", "./...")
    for name in ("worktree", "adaptive"):
        run("govulncheck", "./...", cwd=ROOT / "examples" / name / "consumer")


def lint_report():
    for tags in ([], ["-tags", "integration"]):
        command = ["staticcheck", "-f", "json", *tags, "./..."]
        print(shlex.join(command), flush=True)
        result = subprocess.run(command, cwd=ROOT, text=True, capture_output=True)
        print(result.stdout, end="")
        print(result.stderr, end="")
        diagnostics = [json.loads(line) for line in result.stdout.splitlines() if line.strip()]
        # Findings are advisory here. Tool/build failures still fail the job.
        if result.returncode not in (0, 1) or result.stderr or (result.returncode == 1 and not diagnostics):
            raise RuntimeError("Staticcheck report failed to execute")
        if any(item.get("code") == "compile" or "location" not in item for item in diagnostics):
            raise RuntimeError("Staticcheck could not load/compile the complete scope")
        label = "integration" if tags else "ordinary"
        log_dir = os.environ.get("CI_LOG_DIR")
        if log_dir:
            Path(log_dir, f"staticcheck-{label}.jsonl").write_text(result.stdout)
        message = f"Global {label} Staticcheck: {len(diagnostics)} diagnostics (advisory scope)."
        print(message)
        if diagnostics:
            print(f"::warning::{message}")
        summary = os.environ.get("GITHUB_STEP_SUMMARY")
        if summary:
            with open(summary, "a") as stream:
                stream.write(message + "\n\n")


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("suite", choices=[
        "pr", "regression", "worktree", "adaptive", "security", "lint-report",
        "full-ordinary", "full-integration", "full-adaptive-race",
    ])
    suite = parser.parse_args().suite
    run("go", "version")
    if suite == "pr":
        pr()
    elif suite == "regression":
        regression()
    elif suite in ("worktree", "adaptive"):
        consumer(suite)
    elif suite == "security":
        security()
    elif suite == "lint-report":
        lint_report()
    elif suite == "full-adaptive-race":
        tests(["./..."], "60m", race=True, cwd=ROOT / "examples/adaptive/consumer")
    else:
        tests(["./..."], "60m", integration=suite == "full-integration")


if __name__ == "__main__":
    main()
