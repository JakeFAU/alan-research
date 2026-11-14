"""Assemble full-resolution Poppler renders for visual review; no PDF changes."""
from pathlib import Path
from PIL import Image, ImageDraw
p = Path('/tmp/alan-phase3c4/task3-pdf/renders')
files = sorted(p.glob('page-*.png'))
for n in range(0, len(files), 4):
    batch = files[n:n+4]
    ims = [Image.open(f).convert('RGB') for f in batch]
    w, h = max(im.width for im in ims), max(im.height for im in ims)
    sheet = Image.new('RGB', (2*w, 2*(h+28)), '#dddddd')
    draw = ImageDraw.Draw(sheet)
    for k, (f, im) in enumerate(zip(batch, ims)):
        x, y = (k % 2)*w, (k//2)*(h+28)
        draw.text((x+10, y+5), f.stem, fill='black')
        sheet.paste(im, (x, y+28))
    sheet.save(p / f'sheet-{n//4:02}.png')
print(len(files), 'pages;', len(list(p.glob('sheet*'))), 'sheets')
