package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	appTitle   = "Spolszczenie BATMAN: The Telltale Series"
	appVersion = "1.2"
	appAuthor  = "Januzyne"
	gameDir    = "Batman The Telltale Series"
	checkFile  = "BM_pc_Batman101_data.ttarch2"
	backupDir  = "_spolszczenie_kopia"
)

type manifest struct {
	Version string          `json:"version"`
	Files   map[string]bool `json:"files"`
}

var payloadZip []byte
var setStatus = func(s string, val, max int) {}

func validGame(p string) bool {
	_, err := os.Stat(filepath.Join(p, "Archives", checkFile))
	return err == nil
}

func normalize(p string) string {
	p = strings.TrimSpace(strings.Trim(p, `"`))
	if strings.EqualFold(filepath.Base(p), "Archives") && validGame(filepath.Dir(p)) {
		return filepath.Dir(p)
	}
	return p
}

func loadManifest(game string) *manifest {
	m := &manifest{Files: map[string]bool{}}
	b, err := os.ReadFile(filepath.Join(game, backupDir, "manifest.json"))
	if err == nil {
		json.Unmarshal(b, m)
		if m.Files == nil {
			m.Files = map[string]bool{}
		}
	}
	return m
}

func saveManifest(game string, m *manifest) error {
	b, _ := json.MarshalIndent(m, "", "  ")
	return os.WriteFile(filepath.Join(game, backupDir, "manifest.json"), b, 0644)
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		return err
	}
	return out.Close()
}

func lockErr(err error) string {
	return fmt.Sprintf("%v\n\nUpewnij się, że gra jest zamknięta, i spróbuj ponownie.", err)
}

func install(game string) error {
	zr, err := zip.NewReader(bytes.NewReader(payloadZip), int64(len(payloadZip)))
	if err != nil {
		return err
	}
	arch := filepath.Join(game, "Archives")
	bdir := filepath.Join(game, backupDir)
	if err := os.MkdirAll(bdir, 0755); err != nil {
		return err
	}
	m := loadManifest(game)
	total := len(zr.File)
	for i, f := range zr.File {
		name := filepath.Base(f.Name)
		setStatus("Instalowanie: "+name, i, total)
		dst := filepath.Join(arch, name)
		if _, known := m.Files[name]; !known {
			backed := false
			if _, err := os.Stat(dst); err == nil {
				if err := copyFile(dst, filepath.Join(bdir, name)); err != nil {
					return fmt.Errorf("kopia zapasowa %s: %v", name, err)
				}
				backed = true
			}
			m.Files[name] = backed
			if err := saveManifest(game, m); err != nil {
				return err
			}
		}
		rc, err := f.Open()
		if err != nil {
			return err
		}
		tmp := dst + ".pltmp"
		out, err := os.Create(tmp)
		if err != nil {
			rc.Close()
			return fmt.Errorf("%s: %v", name, err)
		}
		_, err = io.Copy(out, rc)
		rc.Close()
		out.Close()
		if err != nil {
			os.Remove(tmp)
			return err
		}
		if err := os.Rename(tmp, dst); err != nil {
			os.Remove(dst)
			if err2 := os.Rename(tmp, dst); err2 != nil {
				os.Remove(tmp)
				return fmt.Errorf("%s: %v", name, err2)
			}
		}
	}
	m.Version = appVersion
	saveManifest(game, m)
	setStatus("Gotowe! Spolszczenie zostało zainstalowane.", total, total)
	return nil
}

func uninstall(game string) error {
	m := loadManifest(game)
	if len(m.Files) == 0 {
		return fmt.Errorf("Nie znaleziono zainstalowanego spolszczenia w tym folderze.")
	}
	arch := filepath.Join(game, "Archives")
	bdir := filepath.Join(game, backupDir)
	total, i := len(m.Files), 0
	for name, backed := range m.Files {
		setStatus("Przywracanie: "+name, i, total)
		dst := filepath.Join(arch, name)
		if err := os.Remove(dst); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("%s: %v", name, err)
		}
		if backed {
			if err := copyFile(filepath.Join(bdir, name), dst); err != nil {
				return fmt.Errorf("przywracanie %s: %v", name, err)
			}
		}
		delete(m.Files, name)
		saveManifest(game, m)
		i++
	}
	os.RemoveAll(bdir)
	setStatus("Spolszczenie zostało odinstalowane. Przywrócono oryginalne pliki.", total, total)
	return nil
}
