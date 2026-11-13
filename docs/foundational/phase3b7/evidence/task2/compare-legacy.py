import hashlib,json,subprocess
from pathlib import Path
root=Path('/tmp/alan-phase3b5/corpus')
m=json.loads((root/'manifest.json').read_bytes());result=[]
for mode in ['normal','race']:
 for case in ['N07','applied-unknown-revoked-before-fresh-gate']:
  path=root/('raw-'+mode)/m['raw-'+mode+'.tar.gz']['cases'][case]
  old=subprocess.run(['/tmp/alan-phase3b7/old-effect-check','-case',str(path)],capture_output=True)
  new=subprocess.run(['/tmp/alan-phase3b7/new-effect-check','-case',str(path)],capture_output=True)
  assert old.returncode==new.returncode==0,(case,old.returncode,new.returncode)
  assert old.stdout==new.stdout and old.stderr==new.stderr==b'',case
  result.append({'mode':mode,'case':case,'bytes':len(new.stdout),'sha256':hashlib.sha256(new.stdout).hexdigest(),'old_exit':old.returncode,'new_exit':new.returncode,'identical':True})
Path('docs/foundational/phase3b7/evidence/task2/historical-byte-regression.json').write_text(json.dumps({'base':'f238b778ac786d3efbd86f94e98068439b9e4ed2','reports':result},indent=2)+'\n')
print('4 historical normal/race exports byte-identical')
