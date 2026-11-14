"""Render a captured build into a fresh directory without publishing evidence."""
from pathlib import Path
import argparse
import gzip
import hashlib
import json
import shutil
import subprocess
import re


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('build', type=Path, help='output of tools/build_pdf.py')
    parser.add_argument('output', type=Path, help='new directory; parent must exist')
    args = parser.parse_args()
    try:
        args.output.mkdir()
    except FileExistsError:
        raise SystemExit('REFUSE: output already exists: ' + str(args.output))
    output, build = args.output.resolve(), args.build.resolve()
    manifest = json.loads((build/'manifest.json').read_text())
    if manifest['status'] != 'PASS':
        raise SystemExit('REFUSE: input build did not pass')
    for row in manifest['files']:
        raw = (build/row['path']).read_bytes()
        if len(raw) != row['bytes'] or hashlib.sha256(raw).hexdigest() != row['sha256']:
            raise SystemExit('REFUSE: input build identity changed: '+row['path'])
    shutil.copyfile(build/'manifest.json', output/'input-build-manifest.json')
    shutil.copyfile(build/'build-commands.json', output/'build-commands.json')
    pages, warnings = {}, {}
    for name in ['base', 'current']:
        for source in ['alan-dissertation.pdf', 'alan-dissertation.log', 'text.txt', 'pdfinfo.txt']:
            (output/(name+'-'+source+'.gz')).write_bytes(gzip.compress((build/name/source).read_bytes(), mtime=0))
        pages[name] = (build/name/'text.txt').read_text().split('\f')
        if not pages[name][-1].strip():
            pages[name].pop()
        warnings[name] = [line for line in (build/name/'alan-dissertation.log').read_text().splitlines()
                          if 'Warning:' in line or line.startswith(('Underfull', 'Overfull'))]
    changed = [i+1 for i, text in enumerate(pages['current']) if i >= len(pages['base']) or text != pages['base'][i]]
    coverage = sorted({j for i in changed for j in [i-1, i, i+1] if 1 <= j <= len(pages['current'])})
    normalize = lambda lines: [re.sub(r'(lines? )\d+(?:--\d+)?', r'\1N', line) for line in lines]
    comparison = dict(page_counts={k:len(v) for k,v in pages.items()}, warnings=warnings,
                      normalized_warnings_equal=normalize(warnings['base']) == normalize(warnings['current']),
                      changed_pages=changed, rendered_pages=coverage,
                      visual_review='NOT PERFORMED: rendering alone is not visual inspection')
    (output/'comparison.json').write_text(json.dumps(comparison, indent=2)+'\n')
    renders = output/'renders'
    renders.mkdir()
    commands = []
    for page in coverage:
        argv = ['pdftoppm', '-f', str(page), '-l', str(page), '-r', '100', '-png', '-singlefile',
                str(build/'current/alan-dissertation.pdf'), str(renders/f'page-{page:02}')]
        result = subprocess.run(argv, capture_output=True, text=True)
        commands.append(dict(argv=argv, exit_code=result.returncode, stdout=result.stdout, stderr=result.stderr))
        (output/'render-commands.json').write_text(json.dumps(commands, indent=2)+'\n')
        if result.returncode:
            raise SystemExit('Render failed; partial output retained in '+str(output))
    files = []
    for path in sorted(output.rglob('*')):
        if path.is_file():
            raw = path.read_bytes()
            files.append(dict(path=str(path.relative_to(output)), bytes=len(raw), sha256=hashlib.sha256(raw).hexdigest()))
    (output/'manifest.json').write_text(json.dumps(dict(status='PASS', files=files), indent=2)+'\n')
    print('PASS: rendered into '+str(output)+'; no repository artifact installed; visual review not performed')


if __name__ == '__main__':
    main()
