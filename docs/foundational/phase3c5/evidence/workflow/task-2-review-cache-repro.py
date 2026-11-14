from pathlib import Path
import hashlib, os, py_compile, subprocess, tempfile, sys
root=Path.cwd()
with tempfile.TemporaryDirectory(prefix='phase3c5-cache-review-') as tmp:
 target=Path(tmp).resolve();sys.pycache_prefix=str(target/'cache')
 for phase in ('phase3c2','phase3c3','phase3c4','phase3c5'):
  rel=Path('docs/foundational')/phase/'tools/audit.py'; p=target/rel;p.parent.mkdir(parents=True,exist_ok=True);p.write_bytes((root/rel).read_bytes())
 p=target/'docs/foundational/phase3c4/tools/audit.py';original=p.read_bytes();before=p.stat();old=b'ALPHABET = prior.ALPHABET';new=b"ALPHABET = ('CACHE',)".ljust(len(old));assert len(new)==len(old)
 p.write_bytes(original.replace(old,new,1));os.utime(p,(before.st_atime,before.st_mtime));py_compile.compile(str(p),doraise=True)
 p.write_bytes(original);os.utime(p,(before.st_atime,before.st_mtime));assert hashlib.sha256(p.read_bytes()).hexdigest()=='163fb14ea65af59a507e52b88c188c0b4fe2df65b40e8d2d5c7aef7f076fc88e'
 code="import importlib.util; s=importlib.util.spec_from_file_location('review', 'docs/foundational/phase3c5/tools/audit.py'); m=importlib.util.module_from_spec(s); s.loader.exec_module(m); print('Auditor imported; effective ALPHABET:', m.ALPHABET)"
 result=subprocess.run(['python3','-B','-X','pycache_prefix='+sys.pycache_prefix,'-c',code],cwd=target,capture_output=True,text=True)
 print('All repository files untouched; temporary predecessor source restored to pinned SHA256.');print('Exit:',result.returncode);print(result.stdout,end='');print(result.stderr,end='')
