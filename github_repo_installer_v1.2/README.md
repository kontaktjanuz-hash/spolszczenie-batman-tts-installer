# Spolszczenie Batman TTS — installer (source)

Open-source installer for a **free, fan-made Polish translation** of
**Batman: The Telltale Series** (Steam / Epic, PC). This repository contains
only the **installer source code**, published so it can be verified (e.g. for
Nexus Mods file review). It does **not** contain the translation files.

Author: **Januzyne** · Current version: **1.2**

## What the installer does
- Detects the game folder (Steam auto-detect; Epic via "Browse").
- Backs up the original game files it replaces.
- Copies the translation files into `<game>\Archives\`.
- Can fully restore the originals ("Odinstaluj" / Uninstall).

Small Windows GUI app written in **Go** (GUI: `github.com/lxn/walk`).

### What it does NOT do
- No network connections of any kind.
- No persistence, no services, no scheduled tasks.
- No obfuscation or packing.
- Only registry access is **reading** Steam's install path.

At runtime it reads an **external** data file `spolszczenie_dane.bin`
(an xz-compressed zip of the translation files) next to the `.exe`. The
executable itself contains no game data.

## Build
Requires **Go 1.24+**. From the `installer/` folder:

```sh
set GOOS=windows
set GOARCH=amd64
go build -trimpath -o Spolszczenie_Batman_TTS_v1.2.exe .
```

Icon, version info and the app manifest come from `rsrc_windows_amd64.syso`,
a standard Windows resource object linked automatically by `go build`.

The version shown in the installer comes from one constant: `appVersion` in
`installer/core.go`. The banner image (`installer/banner.png`) is produced by
`tools/make_banner.py <version> <out.png>`, so the version in the banner is
regenerated with a single command.

## Files
| Path | Purpose |
|------|---------|
| `installer/main.go` | GUI + payload loader (reads `spolszczenie_dane.bin`) |
| `installer/core.go` | install / uninstall / backup logic; `appVersion` constant |
| `installer/go.mod`, `go.sum` | dependencies |
| `installer/app.manifest` | Windows application manifest |
| `installer/icon.ico`, `banner.png` | resources |
| `installer/rsrc_windows_amd64.syso` | compiled Windows resource (icon + version info) |
| `tools/make_banner.py` | regenerates the banner with a given version number |

## License
Installer source code: **MIT** (see `LICENSE`). The Polish translation content
and fonts are **not** part of this repository and are not covered by it. This
is a non-commercial fan project, not affiliated with Telltale Games, Warner
Bros. or DC Comics.
