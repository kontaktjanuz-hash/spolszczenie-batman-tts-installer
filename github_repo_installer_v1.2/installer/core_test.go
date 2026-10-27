package main
import ("os";"path/filepath";"testing";"archive/zip";"bytes")
func TestRoundtrip(t *testing.T){
	d:=t.TempDir(); a:=filepath.Join(d,"Archives"); os.MkdirAll(a,0755)
	os.WriteFile(filepath.Join(a,checkFile),[]byte("x"),0644)
	os.WriteFile(filepath.Join(a,"ui_menu_english.landb"),[]byte("ORIG"),0644)
	var buf bytes.Buffer; w:=zip.NewWriter(&buf)
	for _,n:=range []string{"ui_menu_english.landb","BM_pc_PC_101_Polish_all.ttarch2"}{f,_:=w.Create(n);f.Write([]byte("PL-"+n))}
	w.Close(); payloadZip=buf.Bytes()
	if err:=install(d);err!=nil{t.Fatal(err)}
	if err:=install(d);err!=nil{t.Fatal(err)} // reinstall must keep original backup
	b,_:=os.ReadFile(filepath.Join(a,"ui_menu_english.landb")); if string(b)!="PL-ui_menu_english.landb"{t.Fatal("not installed")}
	if !validGame(normalize(a)) {t.Fatal("normalize")}
	if err:=uninstall(d);err!=nil{t.Fatal(err)}
	b,_=os.ReadFile(filepath.Join(a,"ui_menu_english.landb")); if string(b)!="ORIG"{t.Fatal("not restored: "+string(b))}
	if _,err:=os.Stat(filepath.Join(a,"BM_pc_PC_101_Polish_all.ttarch2"));err==nil{t.Fatal("not removed")}
	if _,err:=os.Stat(filepath.Join(d,backupDir));err==nil{t.Fatal("backup dir left")}
}
