# Spolszczenie Batman TTS — installer (source)

Open-source installer for a **free, fan-made Polish translation** of
**Batman: The Telltale Series** (Steam, PC). This repository contains only the
**installer source code**, published so it can be verified (e.g. for Nexus Mods
file review). It does **not** contain the translation files themselves.

Author: **Januzyne**

---

## What the installer does

- Detects the Steam installation of the game.
- Backs up the original game files it is about to replace.
- Copies the translation files into `<game>\Archives\`.
- Can fully restore the original files ("Odinstaluj" / Uninstall).

It is a small Windows GUI app written in **Go** (GUI library:
[`github.com/lxn/walk`](https://github.com/lxn/walk)).

### What it does NOT do
- No network connections of any kind.
- No persistence, no scheduled tasks, no services.
- No obfuscation or packing.
- The only registry access is **reading** Steam's install path.

At runtime the installer reads an **external** data file
`spolszczenie_dane.bin` (an xz-compressed zip of the translation files) that
sits next to the `.exe`. The executable itself contains no game data.

---

## Build

Requires **Go 1.24+**. From the `installer/` folder:

```sh
set GOOS=windows
set GOARCH=amd64
go build -trimpath -ldflags "-s -w" -o Spolszczenie_Batman_TTS_v1.0.exe .
```

Icon, version info and the application manifest are provided by
`rsrc_windows_amd64.syso`, a standard Windows resource object that `go build`
links in automatically.

See [`BUILD.md`](BUILD.md) for more detail.

---

## Files

| Path | Purpose |
|------|---------|
| `installer/main.go` | GUI + payload loader (reads `spolszczenie_dane.bin`) |
| `installer/core.go` | install / uninstall / backup logic |
| `installer/go.mod`, `go.sum` | dependencies |
| `installer/app.manifest` | Windows application manifest |
| `installer/icon.ico`, `banner.png` | resources |
| `installer/rsrc_windows_amd64.syso` | compiled Windows resource (icon + version info) |

---

## License

Installer source code: **MIT** (see [`LICENSE`](LICENSE)).

The Polish **translation content and fonts are not part of this repository**
and are not covered by this license. The translation is a non-commercial fan
project and is not affiliated with Telltale Games, Warner Bros. or DC Comics.
