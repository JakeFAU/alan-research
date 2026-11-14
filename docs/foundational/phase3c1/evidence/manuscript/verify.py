"""Verify retained Phase 3C.1 manuscript/qualification bytes; no builds or writes."""
import gzip
import hashlib
import json
import re
from pathlib import Path

ROOT = Path(__file__).resolve().parents[5]
PACKET = ROOT / 'docs/foundational/phase3c1'


def digest(data):
    return hashlib.sha256(data).hexdigest()


manifest = json.loads((PACKET / 'evidence/manuscript/manifest.json').read_text())
for name, expected in manifest['sources']['current'].items():
    assert digest((ROOT / 'paper' / name).read_bytes()) == expected, name
for name, key in [('alan-dissertation.md', 'markdown_sha256'),
                  ('alan-dissertation.pdf', 'pdf_sha256')]:
    assert digest((ROOT / 'paper' / name).read_bytes()) == manifest[key], name
for artifact in manifest['artifacts']:
    packed = (ROOT / artifact['path']).read_bytes()
    raw = gzip.decompress(packed)
    assert digest(packed) == artifact['compressed_sha256']
    assert digest(raw) == artifact['raw_sha256']
    assert len(raw) == artifact['raw_bytes']
    assert packed[3] & 8 == 0 and packed[4:8] == b'\0' * 4
archive = json.loads((PACKET / 'evidence/qualification/archive-manifest.json').read_text())
for artifact in archive['files']:
    retained = (ROOT / artifact['path']).read_bytes()
    raw = gzip.decompress(retained) if artifact['path'].endswith('.gz') else retained
    assert digest(retained) == artifact['retained_sha256']
    assert digest(raw) == artifact['raw_sha256']
    assert len(raw) == artifact['raw_bytes']
md = (ROOT / 'paper/alan-dissertation.md').read_text()
tex = (ROOT / 'paper/alan-dissertation.tex').read_text()
m = md.split('### 6.16 Preserving Named Progress Contracts (Phase 3C.1)\n', 1)[1].split('\n---\n', 1)[0].strip()
t = tex.split(r'\subsection{Preserving Named Progress Contracts (Phase 3C.1)}', 1)[1].split(r'\section{Conclusion}', 1)[0].strip()
t = t.replace(r'\textbf{Conditional source argument.}', '**Conditional source argument.**')
t = t.replace(r"Abadi and Lamport's \emph{Conjoining Specifications}~\cite{abadi1995conjoining}", "Abadi and Lamport's *Conjoining Specifications* (1995)")
t = t.replace(r'\path{docs/foundational/phase3c1/04-results.md}', '[Phase 3C.1 results](../docs/foundational/phase3c1/04-results.md)')
assert m == t, 'new section content differs'
# The altered prose paragraphs are byte-equal once stripped of surrounding markup.
for start in ['Acceptance used a checker corrected after capture.',
              'The Phase 3B investigation is closed within those explicit fragments']:
    mpart = md[md.index(start):].split('\n', 1)[0]
    tpart = tex[tex.index(start):].split('\n', 1)[0]
    assert mpart == tpart, start
paths = [PACKET / name for name in ['04-results.md', '05-qualification.md', '06-reviews.md']]
paths.append(ROOT / 'docs/foundational/phase3/README.md')
links = 0
for path in paths:
    for target in re.findall(r'\[[^\]]+\]\(([^)]+)\)', path.read_text()):
        if '://' in target or target.startswith('#'):
            continue
        assert (path.parent / target.split('#', 1)[0]).exists(), (path, target)
        links += 1
assert (ROOT / 'docs/foundational/phase3c1/04-results.md').exists()
print(json.dumps({'qualification_archived_files': len(archive['files']),
                  'manuscript_compressed_artifacts': len(manifest['artifacts']),
                  'local_links_checked': links,
                  'new_section_exact_after_markup_normalization': True,
                  'changed_abstract_conclusion_prose_equal': True,
                  'current_paper_hashes_match': True,
                  'result': 'PASS'}, indent=2))
