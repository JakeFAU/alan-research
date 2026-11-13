"""Read-only checks of Phase 3B.6 packaging; does not execute tests or SQLite."""
import gzip
import hashlib
import json
from pathlib import Path
import re
import subprocess
import tarfile
from urllib.parse import unquote, urlsplit

ROOT = Path.cwd()
PACKET = ROOT / 'docs/foundational/phase3b6'
EVIDENCE = PACKET / 'evidence'

def sha(data):
    return hashlib.sha256(data).hexdigest()

def load(path):
    return json.loads(path.read_text())

head = subprocess.check_output(['git', 'rev-parse', 'HEAD'], text=True).strip()
tracked = set(subprocess.check_output(['git', 'ls-files', '-z']).decode().split('\0'))
files = [p for p in PACKET.rglob('*') if p.is_file()]
assert all(not p.is_symlink() for p in files)
untracked = [str(p.relative_to(ROOT)) for p in files if str(p.relative_to(ROOT)) not in tracked]
assert not untracked, untracked
sources = load(EVIDENCE / 'task4/qualification-final/source.json')
for name, value in sources.items():
    assert sha((ROOT / name).read_bytes()) == value, name
freeze = load(EVIDENCE / 'task4/freeze.json')
executable_delta = subprocess.check_output([
    'git', 'diff', '--name-only', freeze['head'], '--', '*.go', 'go.mod', 'go.sum', 'research',
], text=True)
assert not executable_delta, executable_delta
untracked_go = subprocess.check_output(['git', 'ls-files', '--others', '--exclude-standard', '*.go'], text=True)
assert not untracked_go, untracked_go
checksums = {}
for task in ('task1', 'task2', 'task3'):
    rows = (EVIDENCE / task / 'SHA256SUMS').read_text().splitlines()
    for row in rows:
        expected, name = row.split(None, 1)
        path = EVIDENCE / task / name.lstrip('*')
        assert sha(path.read_bytes()) == expected, path
    checksums[task] = len(rows)
decoded = 0
for path in files:
    if path.suffix == '.json':
        load(path)
    if path.suffix == '.gz':
        data = gzip.decompress(path.read_bytes())
        decoded += 1
        if path.name.endswith('.json.gz'):
            json.loads(data)
pdfdir = EVIDENCE / 'task4/pdf'
pdf = load(pdfdir / 'manifest.json')
archive = pdfdir / pdf['archive']
assert sha(archive.read_bytes()) == pdf['archive_sha256']
with tarfile.open(archive) as tar:
    members = tar.getmembers()
    assert len(members) == len(pdf['members'])
    assert {m.name for m in members} == set(pdf['members'])
    for member in members:
        assert member.isfile()
        data = tar.extractfile(member).read()
        assert {'bytes': len(data), 'sha256': sha(data)} == pdf['members'][member.name]
    for name in ('alan-dissertation.md', 'alan-dissertation.tex', 'references.bib', 'alan-dissertation.pdf'):
        assert (ROOT / 'paper' / name).read_bytes() == tar.extractfile(name).read(), name
assert sha((ROOT / 'paper/alan-dissertation.pdf').read_bytes()) == pdf['installed_pdf_sha256']
workflowdir = EVIDENCE / 'workflow'
workflow = load(workflowdir / 'manifest.json')['files']
assert {p.name for p in workflowdir.iterdir() if p.name != 'manifest.json'} == set(workflow)
for name, expected in workflow.items():
    data = (workflowdir / name).read_bytes()
    assert {'bytes': len(data), 'sha256': sha(data)} == expected, name
links = []
for path in sorted(PACKET.glob('*.md')):
    for target in re.findall(r'\[[^\]]+\]\(([^)]+)\)', path.read_text()):
        parsed = urlsplit(target)
        if parsed.scheme or not parsed.path:
            continue
        resolved = (path.parent / unquote(parsed.path)).resolve()
        assert resolved.exists(), (path, target)
        if resolved.is_file():
            assert str(resolved.relative_to(ROOT)) in tracked, (path, target)
        else:
            assert any(p.is_file() and str(p.relative_to(ROOT)) in tracked for p in resolved.rglob('*')), (path, target)
        links.append({'source': str(path.relative_to(ROOT)), 'target': target})
assert subprocess.run(['git', 'diff', '--check'], stdout=subprocess.PIPE).returncode == 0
assert subprocess.run(['git', 'diff', '--cached', '--check'], stdout=subprocess.PIPE).returncode == 0
print(json.dumps({
    'head_at_check': head, 'script_sha256': sha(Path(__file__).read_bytes()),
    'frozen_inputs_checked': len(sources), 'packet_files_tracked': len(files),
    'no_go_module_or_research_changes_since_freeze': True,
    'task_checksum_members': checksums, 'gzip_files_decoded': decoded,
    'pdf_archive_members': len(pdf['members']), 'pdf_sources_equal': True,
    'workflow_files_checked': len(workflow), 'relative_links_checked': len(links),
    'relative_links': links, 'unstaged_and_staged_diff_checks_passed': True,
}, indent=2))
