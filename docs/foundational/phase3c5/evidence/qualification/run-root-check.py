from pathlib import Path
import datetime, hashlib, json, os, subprocess, sys, time

COMMANDS = {
    "build": ["go", "build", "./..."],
    "module-verify": ["go", "mod", "verify"],
    "module-tidy-diff": ["go", "mod", "tidy", "-diff"],
    "package-list": ["go", "list", "./..."],
}
name = sys.argv[1]
command = COMMANDS[name]
root = Path.cwd()
out = Path("/tmp/alan-phase3c5/qualification") / name
out.mkdir(parents=True, exist_ok=False)
env = os.environ.copy()
env["GOCACHE"] = "/tmp/alan-phase3c5/go-cache"
env["GOMODCACHE"] = "/tmp/alan-phase3c5/mod-cache"
env["STATICCHECK_CACHE"] = "/tmp/alan-phase3c5/staticcheck-cache"
head = subprocess.check_output(["git", "rev-parse", "HEAD"], text=True).strip()
job = {"argv": command, "cwd": str(root), "head": head,
       "environment": {k: env[k] for k in ("GOCACHE", "GOMODCACHE", "STATICCHECK_CACHE")},
       "started_utc": datetime.datetime.now(datetime.timezone.utc).isoformat()}
(out / "job.json").write_text(json.dumps(job, indent=2) + "\n")
start = time.monotonic()
r = subprocess.run(command, env=env, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
(out / "stdout.log").write_bytes(r.stdout)
(out / "stderr.log").write_bytes(r.stderr)
result = {"exit_code": r.returncode, "elapsed_seconds": round(time.monotonic()-start,3),
          "stdout_bytes": len(r.stdout), "stderr_bytes": len(r.stderr),
          "stdout_sha256": hashlib.sha256(r.stdout).hexdigest(),
          "stderr_sha256": hashlib.sha256(r.stderr).hexdigest()}
(out / "result.json").write_text(json.dumps(result, indent=2) + "\n")
print(json.dumps({"check":name, **result}))
sys.exit(r.returncode)
