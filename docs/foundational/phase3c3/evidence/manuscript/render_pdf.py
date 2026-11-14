"""Compare extracted pages and render changed pages plus both continuations."""
from pathlib import Path
import hashlib
import json
import re
import subprocess

SCRATCH = Path('/tmp/alan-phase3c3/task3-pdf')
texts = {name: (SCRATCH / name / 'text.txt').read_text().split('\f') for name in ['base', 'current']}
for name in texts:
    if not texts[name][-1].strip():
        texts[name].pop()
changed = [i + 1 for i, text in enumerate(texts['current']) if i >= len(texts['base']) or text != texts['base'][i]]
warnings = {name: [line for line in (SCRATCH / name / 'alan-dissertation.log').read_text().splitlines() if re.search(r'(Underfull \\hbox|Overfull \\hbox|LaTeX.*Warning)', line)] for name in texts}
comparison = {'page_counts': {name: len(texts[name]) for name in texts}, 'changed_current_pages': changed,
              'warnings': warnings, 'exact_warning_lines_equal': warnings['base'] == warnings['current']}
(SCRATCH / 'comparison.json').write_text(json.dumps(comparison, indent=2) + '\n')
pages = sorted({p for i in changed for p in [i - 1, i, i + 1] if 1 <= p <= len(texts['current'])})
records = []
(SCRATCH / 'renders').mkdir(exist_ok=True)
prior = {p.name: hashlib.sha256(p.read_bytes()).hexdigest() for p in (SCRATCH / 'renders').glob('page-*.png')}
(SCRATCH / 'initial-render-hashes.json').write_text(json.dumps(prior, indent=2) + '\n')
for p in pages:
    argv = ['pdftoppm', '-f', str(p), '-l', str(p), '-r', '100', '-png', '-singlefile', str(SCRATCH / 'current/alan-dissertation.pdf'), str(SCRATCH / 'renders' / f'page-{p:02}')]
    result = subprocess.run(argv, capture_output=True, text=True)
    records.append({'argv': argv, 'exit_code': result.returncode, 'stdout': result.stdout, 'stderr': result.stderr})
    assert result.returncode == 0
(SCRATCH / 'render-commands.json').write_text(json.dumps(records, indent=2) + '\n')
print(json.dumps(comparison, indent=2))
print('Rendered', pages)
print('Changed since initial visual inspection:', [p for p in pages if prior.get(f'page-{p:02}.png') != hashlib.sha256((SCRATCH / 'renders' / f'page-{p:02}.png').read_bytes()).hexdigest()])
