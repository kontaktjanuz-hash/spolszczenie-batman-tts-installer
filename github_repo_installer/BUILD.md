# Build instructions — Spolszczenie Batman TTS installer

Open-source installer for a free Polish fan-translation of
"Batman: The Telltale Series". GUI written in Go (github.com/lxn/walk).

## What it does
- Finds the Steam game folder, backs up the original files it replaces,
  copies the translation files into `<game>\Archives\`, and can restore
  the originals ("Uninstall").
- At runtime it reads an EXTERNAL data file `spolszczenie_dane.bin`
  (xz-compressed zip of the translation files) placed next to the .exe.
  The .exe itself contains no game data.
- No network access, no registry writes beyond READING Steam's install
  path, no persistence, no obfuscation.

## Build (reproducible)
Requires Go 1.24+. On Windows (or cross-compiling):

    set GOOS=windows
    set GOARCH=amd64
    go build -trimpath -ldflags "-s -w" -o Spolszczenie_Batman_TTS_v1.0.exe .

Version info / icon / app manifest are provided by `rsrc_windows_amd64.syso`
(a standard Windows resource object linked automatically by `go build`).

## Files
- main.go   — GUI, payload loader (reads spolszczenie_dane.bin)
- core.go   — install/uninstall/backup logic
- go.mod / go.sum — dependencies
- app.manifest, icon.ico, banner.png, rsrc_windows_amd64.syso — resources
