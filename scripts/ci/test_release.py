"""Guard release identity and recovery using an in-memory GitHub fixture."""

import copy
import unittest

import release


SHA = "a" * 40
RUN = {
    "id": 42, "run_attempt": 1, "head_sha": SHA, "head_branch": "main",
    "event": "push", "status": "completed", "conclusion": "success",
    "html_url": "https://github.com/example/research/actions/runs/42",
}


class GitHubFixture:
    def __init__(self):
        self.runs = [copy.deepcopy(RUN)]
        self.ref = None
        self.release = None
        self.writes = []
        self.fail_publication = False
        self.racing_tag = False

    def request(self, method, path, data=None, *, missing_ok=False):
        if method == "GET":
            if path.startswith("/actions/workflows/main.yml/runs?"):
                return {"workflow_runs": self.runs}
            if path.startswith("/releases/tags/"):
                return self.release
            if path.startswith("/git/ref/tags/"):
                return self.ref
        self.writes.append((path, data))
        if method == "POST" and path == "/git/refs":
            if self.racing_tag:
                raise RuntimeError("HTTP 422: reference already exists")
            self.ref = {"object": {"type": "commit", "sha": data["sha"]}}
            return self.ref
        if method == "POST" and path == "/releases":
            if self.fail_publication:
                raise RuntimeError("HTTP 503: publication failed")
            self.release = {**data, "html_url": "https://github.com/example/research/releases/tag/" + data["tag_name"]}
            return self.release
        raise AssertionError(f"Unexpected request {method} {path}")


class ReleaseTests(unittest.TestCase):
    def test_semver(self):
        for version in ["v0.1.0", "v12.0.1", "v0.2.0-rc.1", "v0.2.0-0.a-1"]:
            release.validate_version(version)
        for version in ["1.0.0", "v01.0.0", "v1.0", "v1.0.0-01", "v1.0.0-", "v1.0.0+x", "v1.0.0\n"]:
            with self.subTest(version=version), self.assertRaises(ValueError):
                release.validate_version(version)

    def test_rejects_unvalidated_or_failed_commit_without_writes(self):
        for override in [
            {"head_sha": "b" * 40}, {"head_branch": "feature"}, {"event": "workflow_dispatch"},
            {"status": "in_progress", "conclusion": None}, {"conclusion": "failure"},
            {"conclusion": "cancelled"}, {"conclusion": "skipped"},
        ]:
            with self.subTest(override=override):
                api = GitHubFixture()
                api.runs[0].update(override)
                with self.assertRaises(ValueError):
                    release.publish(api, "v0.1.0", SHA)
                self.assertEqual(api.writes, [])

    def test_newer_failed_run_cannot_hide_behind_old_success(self):
        api = GitHubFixture()
        api.runs.append({**RUN, "id": 43, "conclusion": "failure"})
        with self.assertRaises(ValueError):
            release.publish(api, "v0.1.0", SHA)
        self.assertEqual(api.writes, [])

    def test_invalid_identity_is_rejected_before_writes(self):
        for version, sha in [("v01.0.0", SHA), ("v0.1.0", "main"), ("v0.1.0", "aaaa")]:
            with self.subTest(version=version, sha=sha):
                api = GitHubFixture()
                with self.assertRaises(ValueError):
                    release.publish(api, version, sha)
                self.assertEqual(api.writes, [])

    def test_publishes_exact_commit_and_validation_provenance(self):
        api = GitHubFixture()
        release.publish(api, "v0.1.0", SHA)
        self.assertEqual(api.ref["object"]["sha"], SHA)
        self.assertEqual(api.release["target_commitish"], SHA)
        self.assertIn(SHA, api.release["body"])
        self.assertIn(RUN["html_url"], api.release["body"])
        self.assertFalse(api.release["prerelease"])

    def test_prerelease(self):
        api = GitHubFixture()
        release.publish(api, "v0.2.0-rc.1", SHA)
        self.assertTrue(api.release["prerelease"])
        self.assertEqual(api.release["make_latest"], "false")

    def test_rejects_existing_release_and_conflicting_tag(self):
        for existing in ["release", "tag"]:
            with self.subTest(existing=existing):
                api = GitHubFixture()
                if existing == "release":
                    api.release = {"tag_name": "v0.1.0"}
                else:
                    api.ref = {"object": {"type": "commit", "sha": "b" * 40}}
                with self.assertRaises(ValueError):
                    release.publish(api, "v0.1.0", SHA)
                self.assertEqual(api.writes, [])

    def test_tag_creation_race_prevents_publication(self):
        api = GitHubFixture()
        api.racing_tag = True
        with self.assertRaises(RuntimeError):
            release.publish(api, "v0.1.0", SHA)
        self.assertIsNone(api.release)

    def test_retry_after_tag_creation_finishes_same_release(self):
        api = GitHubFixture()
        api.fail_publication = True
        with self.assertRaises(RuntimeError):
            release.publish(api, "v0.1.0", SHA)
        self.assertEqual(api.ref["object"]["sha"], SHA)
        api.fail_publication = False
        release.publish(api, "v0.1.0", SHA)
        self.assertEqual(sum(path == "/git/refs" for path, _ in api.writes), 1)
        self.assertEqual(api.release["target_commitish"], SHA)


if __name__ == "__main__":
    unittest.main()
