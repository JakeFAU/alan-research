from pathlib import Path
import subprocess,shutil,json,time
r=Path.cwd();scratch=Path('/tmp/alan-phase3c2/task3');records=[]
for n in ['alan-dissertation.tex','references.bib']:shutil.copyfile(r/'paper'/n,scratch/'current'/n)
for name in ['base','current']:
 d=scratch/name;argv=['tectonic','--only-cached','--keep-logs','alan-dissertation.tex'];start=time.monotonic()
 p=subprocess.run(argv,cwd=d,stdout=subprocess.PIPE,stderr=subprocess.STDOUT);(d/'build-output.log').write_bytes(p.stdout);records.append({'argv':argv,'cwd':str(d),'exit_code':p.returncode,'elapsed_seconds':round(time.monotonic()-start,3)});assert p.returncode==0,p.stdout
 for argv in [['pdfinfo',str(d/'alan-dissertation.pdf')],['pdftotext','-layout',str(d/'alan-dissertation.pdf'),str(d/'text.txt')]]:
  p=subprocess.run(argv,stdout=subprocess.PIPE,stderr=subprocess.STDOUT);records.append({'argv':argv,'exit_code':p.returncode});assert p.returncode==0
  if argv[0]=='pdfinfo':(d/'pdfinfo.txt').write_bytes(p.stdout)
(scratch/'build-commands.json').write_text(json.dumps(records,indent=2)+'\n');shutil.copyfile(scratch/'current/alan-dissertation.pdf',r/'paper/alan-dissertation.pdf')
print(json.dumps(records,indent=2))
