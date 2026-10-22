//go:build windows

package main

import (
	"bytes"
	"io"
	_ "embed"
	"image"
	_ "image/png"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/lxn/walk"
	"github.com/ulikunitz/xz"
	. "github.com/lxn/walk/declarative"
	"golang.org/x/sys/windows/registry"
)

//go:embed banner.png
var bannerData []byte

// Nazwa pliku z danymi spolszczenia, leżącego obok instalatora.
const dataFileName = "spolszczenie_dane.bin"



var (
	mw       *walk.MainWindow
	pathEdit *walk.LineEdit
	status   *walk.Label
	progress *walk.ProgressBar
	btnInst  *walk.PushButton
	btnUninst *walk.PushButton
)

func steamLibraries() []string {
	var roots []string
	add := func(p string) {
		if p == "" {
			return
		}
		p = filepath.Clean(strings.ReplaceAll(p, "/", `\`))
		for _, r := range roots {
			if strings.EqualFold(r, p) {
				return
			}
		}
		roots = append(roots, p)
	}
	for _, k := range []struct {
		root registry.Key
		path, val string
	}{
		{registry.CURRENT_USER, `Software\Valve\Steam`, "SteamPath"},
		{registry.LOCAL_MACHINE, `SOFTWARE\WOW6432Node\Valve\Steam`, "InstallPath"},
		{registry.LOCAL_MACHINE, `SOFTWARE\Valve\Steam`, "InstallPath"},
	} {
		if key, err := registry.OpenKey(k.root, k.path, registry.QUERY_VALUE); err == nil {
			if v, _, err := key.GetStringValue(k.val); err == nil {
				add(v)
			}
			key.Close()
		}
	}
	re := regexp.MustCompile(`"path"\s+"([^"]+)"`)
	for _, r := range append([]string{}, roots...) {
		b, err := os.ReadFile(filepath.Join(r, "steamapps", "libraryfolders.vdf"))
		if err != nil {
			continue
		}
		for _, m := range re.FindAllStringSubmatch(string(b), -1) {
			add(strings.ReplaceAll(m[1], `\\`, `\`))
		}
	}
	for c := 'C'; c <= 'Z'; c++ {
		d := string(c) + `:\`
		for _, s := range []string{"SteamLibrary", `Steam`, `Program Files (x86)\Steam`, `Program Files\Steam`, `Games\Steam`, `Gry\Steam`} {
			if _, err := os.Stat(filepath.Join(d, s, "steamapps")); err == nil {
				add(filepath.Join(d, s))
			}
		}
	}
	return roots
}

func findGame() string {
	for _, lib := range steamLibraries() {
		p := filepath.Join(lib, "steamapps", "common", gameDir)
		if validGame(p) {
			return p
		}
	}
	return ""
}

func guiStatus(s string, val, max int) {
	mw.Synchronize(func() {
		status.SetText(s)
		progress.SetRange(0, max)
		progress.SetValue(val)
	})
}

func refreshButtons() {
	g := normalize(pathEdit.Text())
	ok := validGame(g)
	btnInst.SetEnabled(ok)
	btnUninst.SetEnabled(ok && len(loadManifest(g).Files) > 0)
}

func run(action func(string) error, okMsg string) {
	g := normalize(pathEdit.Text())
	if !validGame(g) {
		walk.MsgBox(mw, appTitle, "W tym folderze nie ma gry Batman: The Telltale Series.\nWskaż folder gry (ten, w którym jest folder Archives).", walk.MsgBoxIconWarning)
		return
	}
	btnInst.SetEnabled(false)
	btnUninst.SetEnabled(false)
	go func() {
		err := action(g)
		mw.Synchronize(func() {
			if err != nil {
				status.SetText("Wystąpił błąd.")
				walk.MsgBox(mw, appTitle, "Błąd: "+lockErr(err), walk.MsgBoxIconError)
			} else {
				walk.MsgBox(mw, appTitle, okMsg, walk.MsgBoxIconInformation)
			}
			refreshButtons()
		})
	}()
}

func loadPayload() error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	raw, err := os.ReadFile(filepath.Join(filepath.Dir(exe), dataFileName))
	if err != nil {
		return err
	}
	r, err := xz.NewReader(bytes.NewReader(raw))
	if err != nil {
		return err
	}
	b, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	payloadZip = b
	return nil
}

func main() {
	setStatus = guiStatus
	if err := loadPayload(); err != nil {
		walk.MsgBox(nil, appTitle,
			"Nie znaleziono pliku z danymi spolszczenia („"+dataFileName+"”).\n\n"+
				"Rozpakuj całe archiwum ZIP i uruchom instalator z tego samego folderu, "+
				"w którym leży plik „"+dataFileName+"”. Nie przenoś samego pliku .exe.",
			walk.MsgBoxIconError)
		return
	}
	found := findGame()
	initial := "Nie znaleziono gry automatycznie. Kliknij „Przeglądaj…” i wskaż folder gry."
	if found != "" {
		initial = "Znaleziono grę. Kliknij „Zainstaluj”."
	}
	var bannerBmp walk.Image
	if img, _, err := image.Decode(bytes.NewReader(bannerData)); err == nil {
		if bm, err := walk.NewBitmapFromImageForDPI(img, 96); err == nil {
			bannerBmp = bm
		}
	}
	err := MainWindow{
		AssignTo: &mw,
		Title:    appTitle + " " + appVersion + " – autor: " + appAuthor,
		Size:     Size{Width: 680, Height: 620},
		MinSize:  Size{Width: 680, Height: 600},
		Layout:   VBox{Margins: Margins{Left: 16, Top: 14, Right: 16, Bottom: 14}, Spacing: 8},
		Children: []Widget{
			ImageView{Image: bannerBmp, Mode: ImageViewModeZoom, MinSize: Size{Width: 640, Height: 200}, MaxSize: Size{Width: 0, Height: 200}},
			Label{Text: "BATMAN: The Telltale Series – spolszczenie", Font: Font{Family: "Segoe UI", PointSize: 14, Bold: true}},
			Label{Text: "Fanowskie tłumaczenie wszystkich 5 epizodów. Wersja " + appVersion + ".\nOryginalne pliki gry są zapisywane w kopii zapasowej i można je przywrócić przyciskiem „Odinstaluj”."},
			Label{Text: "Autor spolszczenia: " + appAuthor, Font: Font{Family: "Segoe UI", PointSize: 10, Bold: true}},
			Label{Text: "© 2026 " + appAuthor + ". Zakaz rozpowszechniania, modyfikowania i podpisywania jako własne bez zgody autora.\nUdostępniając, podawaj autora i link do oryginalnego wątku na Steam."},
			VSpacer{Size: 6},
			Label{Text: "Folder gry:", Font: Font{Family: "Segoe UI", PointSize: 9, Bold: true}},
			Composite{
				Layout: HBox{MarginsZero: true},
				Children: []Widget{
					LineEdit{AssignTo: &pathEdit, Text: found, OnTextChanged: func() { refreshButtons() }},
					PushButton{Text: "Przeglądaj…", OnClicked: func() {
						dlg := new(walk.FileDialog)
						dlg.Title = "Wskaż folder gry Batman The Telltale Series"
						if ok, _ := dlg.ShowBrowseFolder(mw); ok {
							pathEdit.SetText(normalize(dlg.FilePath))
							if validGame(normalize(dlg.FilePath)) {
								status.SetText("Folder gry poprawny. Kliknij „Zainstaluj”.")
							} else {
								status.SetText("To nie jest folder gry – brak pliku Archives\\" + checkFile + ".")
							}
						}
					}},
				},
			},
			VSpacer{Size: 6},
			Label{AssignTo: &status, Text: initial},
			ProgressBar{AssignTo: &progress},
			VSpacer{},
			Composite{
				Layout: HBox{MarginsZero: true},
				Children: []Widget{
					HSpacer{},
					PushButton{AssignTo: &btnInst, Text: "Zainstaluj", MinSize: Size{Width: 110}, OnClicked: func() {
						run(install, "Spolszczenie zostało zainstalowane.\nMiłej gry!")
					}},
					PushButton{AssignTo: &btnUninst, Text: "Odinstaluj", MinSize: Size{Width: 110}, OnClicked: func() {
						if walk.MsgBox(mw, appTitle, "Na pewno odinstalować spolszczenie i przywrócić oryginalne pliki?", walk.MsgBoxYesNo|walk.MsgBoxIconQuestion) == walk.DlgCmdYes {
							run(uninstall, "Spolszczenie zostało odinstalowane.")
						}
					}},
					PushButton{Text: "Zamknij", MinSize: Size{Width: 90}, OnClicked: func() { mw.Close() }},
				},
			},
		},
	}.Create()
	if err != nil {
		walk.MsgBox(nil, appTitle, err.Error(), walk.MsgBoxIconError)
		return
	}
	refreshButtons()
	mw.Run()
}
