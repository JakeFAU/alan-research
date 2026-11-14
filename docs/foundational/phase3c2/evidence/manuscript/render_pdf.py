from pathlib import Path
import subprocess,json
from PIL import Image,ImageOps,ImageDraw
s=Path('/tmp/alan-phase3c2/task3');pages=[1,2,3,4,5,6,11,12,14,15,22,23,24,25,32,33,34,35,43,44,45,46,47,48];d=s/'renders';d.mkdir(exist_ok=True);records=[]
for n in pages:
 argv=['pdftoppm','-f',str(n),'-l',str(n),'-singlefile','-r','100','-png',str(s/'current/alan-dissertation.pdf'),str(d/f'page-{n:02}')];p=subprocess.run(argv,stdout=subprocess.PIPE,stderr=subprocess.PIPE);assert p.returncode==0;records.append({'argv':argv,'exit_code':p.returncode})
for k in range(0,len(pages),4):
 sheet=Image.new('RGB',(1700,2440),'#ddd');draw=ImageDraw.Draw(sheet)
 for j,n in enumerate(pages[k:k+4]):
  im=Image.open(d/f'page-{n:02}.png');im.thumbnail((830,1180));x=(j%2)*850;y=(j//2)*1220;sheet.paste(im,(x,y+25));draw.text((x+10,y+5),f'PDF page {n}',fill='black')
 sheet.save(d/f'sheet-{k//4+1}.png')
(s/'render-commands.json').write_text(json.dumps(records,indent=2)+'\n');print(pages)
