"""Focused CLI qualification: fresh success, refusal, and artifact preservation."""
from pathlib import Path
import hashlib
import json
import subprocess
import sys

ROOT = Path(__file__).resolve().parents[4]
PACKET = ROOT/'docs/foundational/phase3c5'
HISTORICAL = 'fbcedaee554582ec9b2ca5b6ea8f873458516894'


def digest(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()


def tree(folder):
    return {str(p.relative_to(folder)):digest(p) for p in sorted(folder.rglob('*')) if p.is_file()}


def main():
    if len(sys.argv) != 2:
        raise SystemExit('usage: check_pdf_capture.py FRESH_OUTPUT_DIRECTORY')
    output = Path(sys.argv[1])
    try:
        output.mkdir()
    except FileExistsError:
        raise SystemExit('REFUSE: output already exists: '+str(output))
    output = output.resolve()
    names = subprocess.check_output(['git','ls-tree','-r','--name-only',HISTORICAL,'--',
        'docs/foundational/phase3c5/evidence','paper/alan-dissertation.md',
        'paper/alan-dissertation.tex','paper/references.bib','paper/alan-dissertation.pdf'], cwd=ROOT).decode().splitlines()
    retired = {f'docs/foundational/phase3c5/evidence/manuscript/{n}.py' for n in ['build_pdf','render_pdf']}
    protected = {p:digest(ROOT/p) for p in names if p not in retired}
    for path, sha in protected.items():
        original = subprocess.check_output(['git','show',HISTORICAL+':'+path],cwd=ROOT)
        if hashlib.sha256(original).hexdigest() != sha:
            raise SystemExit('Historical protected input already changed: '+path)
    (output/'protected-before.json').write_text(json.dumps(protected,indent=2)+'\n')
    commands = []

    def run(argv, expect, label):
        result = subprocess.run(argv,cwd=ROOT,capture_output=True,text=True)
        commands.append(dict(label=label,argv=argv,exit_code=result.returncode,
                             stdout=result.stdout,stderr=result.stderr))
        (output/'commands.json').write_text(json.dumps(commands,indent=2)+'\n')
        if result.returncode != expect:
            raise SystemExit('Unexpected exit: '+label)
        if expect and 'REFUSE:' not in result.stderr:
            raise SystemExit('Missing refusal diagnostic: '+label)

    build = [sys.executable,str(PACKET/'tools/build_pdf.py'),str(output/'build')]
    render = [sys.executable,str(PACKET/'tools/render_pdf.py'),str(output/'build'),str(output/'render')]
    run(build,0,'fresh build')
    run(render,0,'fresh render')
    before = {n:tree(output/n) for n in ['build','render']}
    run(build,1,'existing build refusal')
    run(render,1,'existing render refusal')
    for name in ['build_pdf','render_pdf']:
        run([sys.executable,str(PACKET/f'evidence/manuscript/{name}.py')],1,'retired '+name+' refusal')
    # Existing files and dangling symlinks also reserve a path.
    occupied = output/'occupied-file'
    occupied.write_text('retain this sentinel\n')
    link = output/'occupied-link'
    link.symlink_to(output/'absent-target')
    for path in [occupied,link]:
        run(build[:-1]+[str(path)],1,'build '+path.name+' refusal')
        run(render[:-1]+[str(path)],1,'render '+path.name+' refusal')
    if occupied.read_text() != 'retain this sentinel\n' or not link.is_symlink() or (output/'absent-target').exists():
        raise SystemExit('Occupied path changed')
    after = {n:tree(output/n) for n in ['build','render']}
    if before != after:
        raise SystemExit('Existing output changed during refusal')
    protected_after = {p:digest(ROOT/p) for p in protected}
    if protected != protected_after:
        raise SystemExit('Committed evidence or manuscript changed')
    (output/'protected-after.json').write_text(json.dumps(protected_after,indent=2)+'\n')
    result = dict(status='PASS',historical_commit=HISTORICAL,protected_files=len(protected),
                  protected_unchanged=True,existing_build_render_unchanged=True,
                  build_files=len(before['build']),render_files=len(before['render']),
                  fresh_successes=2,refusals=8,visual_review='Not performed; capture behavior only')
    (output/'result.json').write_text(json.dumps(result,indent=2)+'\n')
    (output/'output-identities.json').write_text(json.dumps(before,indent=2)+'\n')
    print(json.dumps(result,sort_keys=True))


if __name__ == '__main__':
    main()
