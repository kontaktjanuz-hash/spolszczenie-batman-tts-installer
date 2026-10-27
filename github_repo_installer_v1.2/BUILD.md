# Build instructions

Requires Go 1.24+. From the `installer/` folder:

    set GOOS=windows
    set GOARCH=amd64
    go build -trimpath -o Spolszczenie_Batman_TTS_v1.2.exe .

- The installer loads an external data file `spolszczenie_dane.bin` placed next
  to the .exe; the executable contains no game data.
- No network access, no persistence, no obfuscation.
- `rsrc_windows_amd64.syso` supplies the icon / version info / manifest.
- Version is set once via `appVersion` in `core.go`; the banner is regenerated
  with `tools/make_banner.py <version> installer/banner.png`.
