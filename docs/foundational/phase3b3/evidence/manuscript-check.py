"""Compare only the edited manuscript slice; do not certify historical structure."""
from pathlib import Path
import hashlib
import json
import re
import subprocess
import sys

repo = Path(__file__).resolve().parents[4]
md = (repo / 'paper/alan-dissertation.md').read_text()
tex = (repo / 'paper/alan-dissertation.tex').read_text()

def between(text, start, end):
    assert text.count(start) == 1, start
    return text.split(start, 1)[1].split(end, 1)[0].strip()

def normalize(text):
    text = re.sub(r'\[Phase 3B\.([23]) results\]\(\.\./docs/foundational/phase3b\1/([^)]*)\)',
                  lambda m: 'docs/foundational/phase3b' + m[1] + '/' + m[2], text)
    text = re.sub(r'\\path\{([^}]*)\}', r'\1', text)
    text = re.sub(r'\\paragraph\{([^}]*)\}', r'\1', text)
    text = text.replace('**', '').replace('`', '')
    return re.sub(r'\s+', ' ', text).strip()

pairs = {
    'Abstract': (between(md, '**Abstract.**', '\n---'), between(tex, r'\begin{abstract}', r'\end{abstract}')),
    'Phase3B3': (between(md, '### 6.11 Request-Bound Delegation: the Phase 3B.3 Realization Boundary', '\n---'),
                 between(tex, r'\subsection{Request-Bound Delegation: the Phase 3B.3 Realization Boundary}', r'\section{Conclusion}')),
    'Conclusion': (between(md, "## 8. Conclusion: ALAN's Contributions and Open Mathematical Problems", '\n---'),
                   between(tex, r'\section{Conclusion}', r'\bibliographystyle')),
}
# The only changed Phase3B2 paragraph has an equivalent results link/path.
start = 'An independent retained-history ledger checks ancestry, full tuples, bindings, temporal decisions and custody rather than final capability sets alone.'
pairs['Phase3B2BoundaryParagraph'] = (between(md, start, '\n\n'), between(tex, start, '\n\n'))
results = {}
for name, (a, b) in pairs.items():
    a, b = normalize(a), normalize(b)
    assert a == b, name + ' differs'
    results[name] = {'Aligned': True, 'NormalizedSHA256': hashlib.sha256(a.encode()).hexdigest()}
bib = repo / 'paper/references.bib'
base = subprocess.check_output(['git', 'show', '8197683c23a3ca774750fc3c2c389c5a64dd38b8:paper/references.bib'], cwd=repo)
assert bib.read_bytes() == base
result = {'Checks': results, 'BibTeXUnchanged': True,
          'Files': {name: hashlib.sha256((repo / name).read_bytes()).hexdigest() for name in ['paper/alan-dissertation.md', 'paper/alan-dissertation.tex', 'paper/references.bib', 'paper/alan-dissertation.pdf']},
          'Scope': 'Edited paragraphs only; no universal manuscript consistency claim.',
          'HistoricalDifferences': ['Markdown new subsection6.11 is TeX/PDF6.6 because preceding structure differs.', 'Markdown has a separate historical falsification section absent from TeX; Conclusions are8 versus7.']}
output = json.dumps(result, indent=2) + '\n'
if len(sys.argv) == 2:
    Path(sys.argv[1]).write_text(output)
else:
    print(output, end='')
