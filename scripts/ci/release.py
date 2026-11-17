"""Publish a semantic version of the exact main commit validated by CI."""

import argparse
import json
import os
import re
from urllib.error import HTTPError
from urllib.parse import quote, urlencode
from urllib.request import Request, urlopen


def validate_version(version):
    number = r"(?:0|[1-9][0-9]*)"
    identifier = r"(?:0|[1-9][0-9]*|[0-9]*[A-Za-z-][0-9A-Za-z-]*)"
    if not re.fullmatch(rf"v{number}\.{number}\.{number}(?:-{identifier}(?:\.{identifier})*)?", version):
        raise ValueError("Expected vMAJOR.MINOR.PATCH, optionally followed by a SemVer prerelease")


class GitHub:
    def __init__(self, repository, token):
        self.base = f"https://api.github.com/repos/{repository}"
        self.token = token

    def request(self, method, path, data=None, *, missing_ok=False):
        payload = None if data is None else json.dumps(data).encode()
        request = Request(self.base + path, data=payload, method=method, headers={
            "Authorization": f"Bearer {self.token}",
            "Accept": "application/vnd.github+json",
            "Content-Type": "application/json",
            "X-GitHub-Api-Version": "2022-11-28",
        })
        try:
            with urlopen(request, timeout=30) as response:
                return json.load(response)
        except HTTPError as error:
            if error.code == 404 and missing_ok:
                return None
            raise RuntimeError(f"GitHub {method} {path} returned HTTP {error.code}") from error


def validated_run(api, sha):
    query = urlencode({"branch": "main", "event": "push", "head_sha": sha, "per_page": 100})
    response = api.request("GET", f"/actions/workflows/main.yml/runs?{query}")
    runs = [run for run in response["workflow_runs"] if
            run["head_sha"] == sha and run["head_branch"] == "main" and run["event"] == "push"]
    if not runs:
        raise ValueError(f"No post-merge validation exists for {sha}")
    # A newer failed or running attempt must not be hidden by an older green run.
    latest = max(runs, key=lambda run: (run["id"], run.get("run_attempt", 1)))
    if latest["status"] != "completed" or latest["conclusion"] != "success":
        raise ValueError(f"Post-merge validation is {latest['status']}/{latest['conclusion']}")
    return latest


def publish(api, version, sha):
    validate_version(version)
    if not re.fullmatch(r"[0-9a-f]{40}", sha):
        raise ValueError("A full commit SHA is required")
    run = validated_run(api, sha)
    tag = quote(version, safe="")
    if api.request("GET", f"/releases/tags/{tag}", missing_ok=True) is not None:
        raise ValueError(f"Release {version} already exists")
    ref_path = f"/git/ref/tags/{tag}"
    ref = api.request("GET", ref_path, missing_ok=True)
    if ref is None:
        # Creating the ref explicitly makes a racing tag creation fail atomically.
        api.request("POST", "/git/refs", {"ref": f"refs/tags/{version}", "sha": sha})
    elif ref["object"]["type"] != "commit" or ref["object"]["sha"] != sha:
        raise ValueError(f"Tag {version} already identifies a different object")
    # A retry after tag creation can finish publication using the same SHA.
    ref = api.request("GET", ref_path)
    if ref["object"]["type"] != "commit" or ref["object"]["sha"] != sha:
        raise ValueError(f"Tag {version} changed before publication")
    provenance = (
        f"Research commit `{sha}`.\n\n"
        f"Post-merge validation [run {run['id']}, attempt {run.get('run_attempt', 1)}]({run['html_url']}).\n\n"
        "Validation covers the CI scopes documented in `docs/ci.md`.\n"
    )
    release = api.request("POST", "/releases", {
        "tag_name": version,
        "target_commitish": sha,
        "name": version,
        "body": provenance,
        "generate_release_notes": True,
        "prerelease": "-" in version,
        "make_latest": "false" if "-" in version else "true",
    })
    return release["html_url"]


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("version")
    version = parser.parse_args().version
    if os.environ.get("GITHUB_REF") != "refs/heads/main" or os.environ.get("GITHUB_EVENT_NAME") != "workflow_dispatch":
        raise ValueError("Release requires a manual dispatch on main")
    api = GitHub(os.environ["GITHUB_REPOSITORY"], os.environ["GH_TOKEN"])
    sha = os.environ["RELEASE_SHA"]
    url = publish(api, version, sha)
    print(f"Published {version}: {url}")
    summary = os.environ.get("GITHUB_STEP_SUMMARY")
    if summary:
        with open(summary, "a") as stream:
            stream.write(f"Published [{version}]({url}) from `{sha}`.\n")


if __name__ == "__main__":
    main()
