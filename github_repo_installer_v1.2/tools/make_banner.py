# -*- coding: utf-8 -*-
"""Generator bannera instalatora z parametrem wersji.
Uzycie: python3 make_banner.py <wersja> <plik_wyjsciowy>
np. python3 make_banner.py 1.2 /root/bat/installer/src/banner.png
Dzieki temu numer wersji w banerze zmienia sie jedna komenda."""
import sys, random, math
from PIL import Image, ImageDraw, ImageFont

W, H = 1280, 400
YELLOW = (240, 196, 40)
WHITE = (245, 247, 250)
GREY = (176, 184, 198)

DEJAVU = '/usr/share/fonts/truetype/dejavu/'
def font(name, size):
    return ImageFont.truetype(DEJAVU + name, size)

def make(version, out):
    random.seed(7)
    img = Image.new('RGB', (W, H))
    px = img.load()
    # pionowy gradient granatu nocnego nieba
    for y in range(H):
        t = y / H
        r = int(10 + 8 * (0.5 - abs(t - 0.35)))
        g = int(13 + 9 * (0.5 - abs(t - 0.35)))
        b = int(28 + 18 * (1 - t))
        for x in range(W):
            px[x, y] = (max(4, r), max(6, g), max(12, b))
    d = ImageDraw.Draw(img, 'RGBA')
    # ksiezyc + poswiata
    mx, my, mr = 1015, 78, 34
    for rr in range(120, mr, -2):
        a = int(26 * (1 - (rr - mr) / (120 - mr)))
        d.ellipse([mx - rr, my - rr, mx + rr, my + rr], fill=(230, 232, 210, a))
    d.ellipse([mx - mr, my - mr, mx + mr, my + mr], fill=(236, 238, 220))
    # smuga swiatla ksiezyca
    d.polygon([(mx, my), (mx - 120, H), (mx + 60, H)], fill=(220, 222, 200, 10))
    # sylwetka miasta: kilka warstw budynkow
    def skyline(base_y, col, n, maxh, seed):
        random.seed(seed); x = -20
        while x < W + 20:
            bw = random.randint(40, 95); bh = random.randint(maxh // 3, maxh)
            d.rectangle([x, base_y - bh, x + bw, H], fill=col)
            # anteny
            if random.random() < 0.3:
                ax = x + bw // 2
                d.line([(ax, base_y - bh), (ax, base_y - bh - random.randint(15, 40))], fill=col, width=3)
            x += bw + random.randint(-6, 10)
    skyline(H, (18, 22, 40), 1, 150, 1)
    skyline(H, (9, 12, 24), 1, 210, 2)
    skyline(H, (4, 6, 13), 1, 170, 3)
    # okna (zolte kropki) na dolnej warstwie
    random.seed(11)
    for _ in range(260):
        wx = random.randint(300, W - 10); wy = random.randint(H - 150, H - 8)
        if random.random() < 0.5:
            s = random.choice([3, 3, 4])
            d.rectangle([wx, wy, wx + s, wy + s], fill=(*YELLOW, random.randint(120, 230)))
    # deszcz
    random.seed(21)
    for _ in range(300):
        rx = random.randint(0, W); ry = random.randint(0, H); ln = random.randint(8, 18)
        d.line([(rx, ry), (rx - 3, ry + ln)], fill=(150, 160, 180, 40), width=1)
    # --- teksty ---
    d.text((70, 95), 'BATMAN', font=font('DejaVuSans-Bold.ttf', 92), fill=WHITE)
    d.text((72, 200), 'THE TELLTALE SERIES', font=font('DejaVuSans-Bold.ttf', 34), fill=YELLOW)
    d.line([(74, 250), (330, 250)], fill=YELLOW, width=3)
    info = f'Spolszczenie   •   wszystkie 5 epizodów   •   wersja {version}'
    d.text((72, 285), info, font=font('DejaVuSans.ttf', 26), fill=GREY)
    d.text((72, 330), 'autor: Januzyne', font=font('DejaVuSans-Bold.ttf', 24), fill=YELLOW)
    img.save(out)
    print('zapisano', out, img.size, 'wersja', version)

if __name__ == '__main__':
    make(sys.argv[1], sys.argv[2])
