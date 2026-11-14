#!/usr/bin/env python3
"""Evaluate deliberately false hypotheses. Expected process exit: 1."""
import gzip
import json
from pathlib import Path
import sys

payload = json.loads(gzip.decompress((Path(__file__).parent / "export.json.gz").read_bytes()))
products = {g["ID"]: g for g in payload["Products"]}
if sys.argv[1] == "union-sufficiency":
    witness = payload["Negative"]
    comparison = payload["Comparisons"][witness["Comparison"]]
    left, right = products[witness["LeftProduct"]], products[witness["RightProduct"]]
    antecedent = comparison["InitialLanguage"]["Equal"]
    consequent = left["Nonblocking"] == right["Nonblocking"]
    print(json.dumps({"Hypothesis": "initial-language/derivative-union equality suffices for contextual nonblocking",
        "EqualInitialLanguages": antecedent, "LeftNonblocking": left["Nonblocking"],
        "RightNonblocking": right["Nonblocking"], "ProfilePrefix": comparison["Word"],
        "LeftPath": witness["LeftPath"], "RightPath": witness["RightPath"]}, sort_keys=True))
    sys.exit(0 if not antecedent or consequent else 1)
if sys.argv[1] == "inevitable-completion":
    witness = payload["Loop"]
    graph = products[witness["Product"]]
    nodes = {n["ID"]: n for n in graph["Nodes"]}
    cycle = witness["Cycle"]
    closed_private_unmarked = bool(cycle["Steps"]) and cycle["Start"] == cycle["End"] and all(
        s["Label"] == "tau" and not nodes[s["From"]]["Marked"] and
        {"Label": s["Label"], "To": s["To"]} in nodes[s["From"]]["Arcs"] for s in cycle["Steps"])
    print(json.dumps({"Hypothesis": "nonblocking implies inevitable completion in the abstract grammar",
        "Nonblocking": graph["Nonblocking"], "ClosedPrivateUnmarkedCycle": closed_private_unmarked,
        "Witness": witness, "RuntimeRealization": "unestablished"}, sort_keys=True))
    sys.exit(1 if graph["Nonblocking"] and closed_private_unmarked else 0)
raise SystemExit("unknown hypothesis")
