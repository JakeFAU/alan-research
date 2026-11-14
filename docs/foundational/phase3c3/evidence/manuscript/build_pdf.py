"""Rebuild baseline/current manuscripts with existing cached Tectonic assets."""
from pathlib import Path
import json
import shutil
import subprocess
import time

ROOT = Path(__file__).resolve().parents[5]
SCRATCH = Path('/tmp/alan-phase3c3/task3-pdf')
BASE = 'ab01204a0b5121b81379043576ee7bbe46118a19'
RECORDS = []


def run(argv, cwd):
    started = time.monotonic()
    result = subprocess.run(argv, cwd=cwd, capture_output=True)
    RECORDS.append({'argv': argv, 'cwd': str(cwd), 'exit_code': result.returncode,
                    'elapsed_seconds': round(time.monotonic() - started, 3),
                    'stdout': result.stdout.decode(), 'stderr': result.stderr.decode()})
    (SCRATCH / 'build-commands.json').write_text(json.dumps(RECORDS, indent=2) + '\n')
    assert result.returncode == 0, RECORDS[-1]
    return result.stdout


for name in ['base', 'current']:
    folder = SCRATCH / name
    folder.mkdir(parents=True, exist_ok=True)
    for source in ['alan-dissertation.tex', 'references.bib']:
        raw = subprocess.check_output(['git', 'show', BASE + ':paper/' + source], cwd=ROOT) if name == 'base' else (ROOT / 'paper' / source).read_bytes()
        (folder / source).write_bytes(raw)
    run(['tectonic', '--only-cached', '--keep-logs', 'alan-dissertation.tex'], folder)
    (folder / 'pdfinfo.txt').write_bytes(run(['pdfinfo', str(folder / 'alan-dissertation.pdf')], ROOT))
    run(['pdftotext', '-layout', str(folder / 'alan-dissertation.pdf'), str(folder / 'text.txt')], ROOT)
shutil.copyfile(SCRATCH / 'current/alan-dissertation.pdf', ROOT / 'paper/alan-dissertation.pdf')
print('Baseline/current builds and extraction passed; commands retained in', SCRATCH)
