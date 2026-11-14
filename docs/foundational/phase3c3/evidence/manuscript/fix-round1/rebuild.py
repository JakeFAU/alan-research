"""Capture the current-only build and renders after the P2 wording correction."""
from pathlib import Path
import gzip
import hashlib
import json
import re
import shutil
import subprocess
import time

OUT = Path(__file__).resolve().parent
ROOT = OUT.parents[5]
SCRATCH = Path('/tmp/alan-phase3c3/task3-fix1')
SCRATCH.mkdir(parents=True, exist_ok=True)
prior = json.loads((OUT / 'prior-manifest.json').read_text())
records = []
sha = lambda raw: hashlib.sha256(raw).hexdigest()


def run(argv, cwd=ROOT):
    start = time.monotonic()
    result = subprocess.run(argv, cwd=cwd, capture_output=True)
    records.append({'argv': argv, 'cwd': str(cwd), 'exit_code': result.returncode,
                    'elapsed_seconds': round(time.monotonic() - start, 3),
                    'stdout': result.stdout.decode(), 'stderr': result.stderr.decode()})
    (OUT / 'commands.json').write_text(json.dumps(records, indent=2) + '\n')
    assert result.returncode == 0, records[-1]
    return result.stdout


for name in ['alan-dissertation.tex', 'references.bib']:
    shutil.copyfile(ROOT / 'paper' / name, SCRATCH / name)
run(['tectonic', '--only-cached', '--keep-logs', 'alan-dissertation.tex'], SCRATCH)
info = run(['pdfinfo', str(SCRATCH / 'alan-dissertation.pdf')])
(SCRATCH / 'pdfinfo.txt').write_bytes(info)
run(['pdftotext', '-layout', str(SCRATCH / 'alan-dissertation.pdf'), str(SCRATCH / 'text.txt')])
for name in ['alan-dissertation.log', 'pdfinfo.txt', 'text.txt']:
    (OUT / ('current-' + name + '.gz')).write_bytes(gzip.compress((SCRATCH / name).read_bytes(), mtime=0))
text = (SCRATCH / 'text.txt').read_text().split('\f')
if not text[-1].strip():
    text.pop()
old_text = gzip.decompress((OUT.parent / 'current-text.txt.gz').read_bytes()).decode().split('\f')
if not old_text[-1].strip():
    old_text.pop()
base_text = gzip.decompress((OUT.parent / 'base-text.txt.gz').read_bytes()).decode().split('\f')
if not base_text[-1].strip():
    base_text.pop()
changed = [i + 1 for i, page in enumerate(text) if i >= len(old_text) or page != old_text[i]]
comp = {'prior_current_pages': len(old_text), 'current_pages': len(text),
        'changed_from_reviewed_pages': changed,
        'changed_from_baseline_pages': [i + 1 for i, page in enumerate(text) if i >= len(base_text) or page != base_text[i]],
        'warnings': [line for line in (SCRATCH / 'alan-dissertation.log').read_text().splitlines() if re.search(r'(Underfull \\hbox|Overfull \\hbox|LaTeX.*Warning)', line)],
        'prior_render_sha256': prior['render_sha256'], 'current_render_sha256': {}}
assert comp['warnings'] == prior['warning_comparison']['current']
pages = sorted(set(prior['rendered_and_visually_inspected_pages']) | {p for i in changed for p in [i - 1, i, i + 1] if 1 <= p <= len(text)})
(SCRATCH / 'renders').mkdir(exist_ok=True)
for page in pages:
    output = SCRATCH / 'renders' / f'page-{page:02}'
    run(['pdftoppm', '-f', str(page), '-l', str(page), '-r', '100', '-png', '-singlefile', str(SCRATCH / 'alan-dissertation.pdf'), str(output)])
    comp['current_render_sha256'][str(page)] = sha(output.with_suffix('.png').read_bytes())
comp['changed_render_pages'] = [p for p in pages if comp['current_render_sha256'][str(p)] != prior['render_sha256'].get(str(p))]
(OUT / 'comparison.json').write_text(json.dumps(comp, indent=2) + '\n')
shutil.copyfile(SCRATCH / 'alan-dissertation.pdf', ROOT / 'paper/alan-dissertation.pdf')
print(json.dumps(comp, indent=2))
