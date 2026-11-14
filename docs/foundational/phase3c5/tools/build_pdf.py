"""Build baseline/current PDFs into one fresh directory; never install artifacts."""
from pathlib import Path
import argparse
import hashlib
import json
import subprocess
import time

ROOT = Path(__file__).resolve().parents[4]
BASE = '347d93ba25a5f7a10f75144b1a44107be09a2f39'


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('output', type=Path, help='new directory; parent must exist')
    args = parser.parse_args()
    # mkdir is exclusive: reject files, directories and symlinks before any write.
    try:
        args.output.mkdir()
    except FileExistsError:
        raise SystemExit('REFUSE: output already exists: ' + str(args.output))
    output = args.output.resolve()
    commands, inputs = [], []

    def run(argv, cwd):
        start = time.monotonic()
        result = subprocess.run(argv, cwd=cwd, capture_output=True)
        commands.append(dict(argv=argv, cwd=str(cwd), exit_code=result.returncode,
                             elapsed_seconds=round(time.monotonic()-start, 3),
                             stdout=result.stdout.decode(), stderr=result.stderr.decode()))
        (output/'build-commands.json').write_text(json.dumps(commands, indent=2)+'\n')
        if result.returncode:
            raise SystemExit('Build command failed; records retained in ' + str(output))
        return result.stdout

    for name in ['base', 'current']:
        folder = output/name
        folder.mkdir()
        for source in ['alan-dissertation.tex', 'references.bib']:
            path = 'paper/' + source
            raw = run(['git', 'show', BASE+':'+path], ROOT) if name == 'base' else (ROOT/path).read_bytes()
            (folder/source).write_bytes(raw)
            inputs.append(dict(kind=name, path=path, bytes=len(raw), sha256=hashlib.sha256(raw).hexdigest()))
        (output/'build-inputs.json').write_text(json.dumps(inputs, indent=2)+'\n')
        run(['tectonic', '--only-cached', '--keep-logs', 'alan-dissertation.tex'], folder)
        (folder/'pdfinfo.txt').write_bytes(run(['pdfinfo', str(folder/'alan-dissertation.pdf')], ROOT))
        run(['pdftotext', '-layout', str(folder/'alan-dissertation.pdf'), str(folder/'text.txt')], ROOT)
    files = []
    for path in sorted(output.rglob('*')):
        if path.is_file():
            raw = path.read_bytes()
            files.append(dict(path=str(path.relative_to(output)), bytes=len(raw), sha256=hashlib.sha256(raw).hexdigest()))
    (output/'manifest.json').write_text(json.dumps(dict(status='PASS', files=files), indent=2)+'\n')
    print('PASS: built into ' + str(output) + '; no repository artifact installed')


if __name__ == '__main__':
    main()
