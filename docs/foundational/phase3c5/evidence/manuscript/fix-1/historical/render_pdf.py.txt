"""Archive build outputs and render changed pages plus adjacent continuations."""
from pathlib import Path
import gzip,json,re,subprocess,hashlib,shutil
ROOT=Path(__file__).resolve().parents[5]
OUT=Path(__file__).resolve().parent
TMP=Path('/tmp/alan-phase3c5/task3-pdf')
sha=lambda x:hashlib.sha256(x).hexdigest()
shutil.copyfile(TMP/'build-commands.json',OUT/'build-commands.json')
pages={};warnings={}
for name in ['base','current']:
 for source in ['alan-dissertation.log','text.txt','pdfinfo.txt']:
  (OUT/(name+'-'+source+'.gz')).write_bytes(gzip.compress((TMP/name/source).read_bytes(),mtime=0))
 pages[name]=(TMP/name/'text.txt').read_text().split('\f')
 if not pages[name][-1].strip():pages[name].pop()
 warnings[name]=[x for x in (TMP/name/'alan-dissertation.log').read_text().splitlines() if 'Warning:' in x or x.startswith(('Underfull','Overfull'))]
(OUT/'base-alan-dissertation.pdf.gz').write_bytes(gzip.compress((TMP/'base/alan-dissertation.pdf').read_bytes(),mtime=0))
changed=[i+1 for i,t in enumerate(pages['current']) if i>=len(pages['base']) or t!=pages['base'][i]]
coverage=sorted({j for i in changed for j in [i-1,i,i+1] if 1<=j<=len(pages['current'])})
normalize=lambda s:re.sub(r'(lines? )\d+(?:--\d+)?',r'\1N',s)
comparison={'page_counts':{k:len(v) for k,v in pages.items()},'warnings':warnings,'normalized_warnings_equal':[normalize(x) for x in warnings['base']]==[normalize(x) for x in warnings['current']],'changed_pages':changed,'rendered_pages':coverage}
(OUT/'comparison.json').write_text(json.dumps(comparison,indent=2)+'\n')
renders=TMP/'renders';renders.mkdir(exist_ok=True);records=[]
for n in coverage:
 argv=['pdftoppm','-f',str(n),'-l',str(n),'-r','100','-png','-singlefile',str(ROOT/'paper/alan-dissertation.pdf'),str(renders/f'page-{n:02}')]
 r=subprocess.run(argv,capture_output=True,text=True);records.append(dict(argv=argv,exit_code=r.returncode,stdout=r.stdout,stderr=r.stderr));assert r.returncode==0
(OUT/'render-commands.json').write_text(json.dumps(records,indent=2)+'\n')
from PIL import Image,ImageDraw
for start in range(0,len(coverage),4):
 ns=coverage[start:start+4];imgs=[Image.open(renders/f'page-{n:02}.png').convert('RGB') for n in ns];w,h=imgs[0].size
 sheet=Image.new('RGB',(2*w,2*(h+28)),'#cccccc');d=ImageDraw.Draw(sheet)
 for j,(n,im) in enumerate(zip(ns,imgs)):
  x=(j%2)*w;y=(j//2)*(h+28);d.text((x+8,y+5),f'Page {n}',fill='black');sheet.paste(im,(x,y+28))
 sheet.save(renders/f'sheet-{start//4+1:02}.png')
print(json.dumps(comparison,indent=2))
