package main

import (
	"os"
	"strings"
	"testing"
)

func TestDefaultsMatchGolden(t *testing.T) {
	if err := loadAssets(); err != nil {
		t.Fatal(err)
	}
	g, err := os.ReadFile("golden_default.xml")
	if err != nil {
		t.Skip("golden_default.xml ausente")
	}
	r, err := Generate(DefaultOptions())
	if err != nil {
		t.Fatal(err)
	}
	if r.XML != string(g) {
		t.Fatalf("XML padrão mudou (len %d vs %d)", len(r.XML), len(g))
	}
}

func TestNoNetworkOrShortcutStuff(t *testing.T) {
	if err := loadAssets(); err != nil {
		t.Fatal(err)
	}
	o := DefaultOptions()
	o.DisableNet, o.HideWireless = true, true // mesmo que venham de um perfil antigo
	r, err := Generate(o)
	if err != nil {
		t.Fatal(err)
	}
	for _, bad := range []string{"Disable-NetAdapter", "Enable-NetAdapter", "Install Winhance", "WinhanceInstall", "HideWirelessSetupInOOBE",
		"TcpAckFrequency", "TCPNoDelay", "NetworkThrottlingIndex", `Services\RasMan`, `Services\RasAuto`, `Services\icssvc`, `Services\BITS`,
		"AllowWiFiHotSpotReporting", "AllowAutoConnectToWiFiSenseHotspots", "DODownloadMode", "CreateShortcut"} {
		if strings.Contains(r.XML, bad) {
			t.Errorf("o XML ainda contém %q", bad)
		}
	}
}

func TestPostInstall(t *testing.T) {
	if err := loadAssets(); err != nil {
		t.Fatal(err)
	}
	for _, lang := range []string{"pt-BR", "en", "es"} {
		for _, run := range []bool{true, false} {
			o := DefaultOptions()
			o.Lang = lang
			o.RunScript = run
			o.PostEnabled = true
			o.PostAutoRun = true
			o.PostDelete = true
			o.PostCustom = "Foo.Bar # comentario\n\nMicrosoft.VCRedist.2015+.x64\n"
			r, err := Generate(o)
			if err != nil {
				t.Fatal(err)
			}
			for _, want := range []string{postBatPath, postPs1Path, "Id = 'Google.Chrome'", "Id = 'Foo.Bar'", "winget install --id", `start "" "C:\Users\Public\Desktop\Pos-instalacao.bat"`, "Remove-Item -LiteralPath"} {
				if !strings.Contains(r.XML, want) {
					t.Errorf("[%s run=%v] falta %q", lang, run, want)
				}
			}
			if strings.Count(r.XML, "Id = 'Microsoft.VCRedist.2015+.x64'") != 1 {
				t.Errorf("duplicado ou ausente: VCRedist")
			}
			if run != strings.Contains(r.XML, "Winhancements.ps1") {
				t.Errorf("Winhancements.ps1 inconsistente com RunScript=%v", run)
			}
			if !strings.HasSuffix(r.XML, "  </Extensions>\n</unattend>\n") {
				t.Errorf("fim do XML inesperado")
			}
		}
	}
	o := DefaultOptions()
	o.PostEnabled = true
	o.PostCustom = "x y; calc"
	if _, err := Generate(o); err == nil {
		t.Errorf("ID inválido deveria falhar")
	}
	o = DefaultOptions()
	o.PostEnabled = true
	o.PostApps = nil
	r, _ := Generate(o)
	if r == nil || len(r.Warnings) == 0 {
		t.Errorf("esperava aviso de lista vazia")
	}
}

func TestMultipleUserAccounts(t *testing.T) {
	if err := loadAssets(); err != nil {
		t.Fatal(err)
	}
	o := DefaultOptions()
	o.CreateUser = true
	o.UserAccounts = []UserAccount{
		{Name: "Marco", DisplayName: "Marco Vivas", Password: "secretPassword123", Group: "Administrators"},
		{Name: "Usuario", DisplayName: "Usuário Comum", Password: "", Group: "Users"},
	}
	o.AutoLogon = true
	o.AutoLogonUser = "Marco"
	o.AutoLogonN = 2

	r, err := Generate(o)
	if err != nil {
		t.Fatalf("falha ao gerar com múltiplas contas: %v", err)
	}
	if !strings.Contains(r.XML, "<Name>Marco</Name>") {
		t.Errorf("falta conta Marco")
	}
	if !strings.Contains(r.XML, "<DisplayName>Marco Vivas</DisplayName>") {
		t.Errorf("falta DisplayName Marco Vivas")
	}
	if !strings.Contains(r.XML, "<Group>Administrators</Group>") {
		t.Errorf("falta Group Administrators")
	}
	if !strings.Contains(r.XML, "<Value>secretPassword123</Value>") {
		t.Errorf("falta senha do Marco")
	}
	if !strings.Contains(r.XML, "<Name>Usuario</Name>") {
		t.Errorf("falta conta Usuario")
	}
	if !strings.Contains(r.XML, "<Group>Users</Group>") {
		t.Errorf("falta Group Users")
	}
	if !strings.Contains(r.XML, "<Username>Marco</Username>") {
		t.Errorf("falta AutoLogon Username Marco")
	}
	if !strings.Contains(r.XML, "<LogonCount>2</LogonCount>") {
		t.Errorf("falta LogonCount 2")
	}
}

func TestCustomWallpaper(t *testing.T) {
	if err := loadAssets(); err != nil {
		t.Fatal(err)
	}
	o := DefaultOptions()
	o.WallpaperEnabled = true
	o.WallpaperMode = "url"
	o.WallpaperURL = "https://marcovivas.com/docs/DEV9.png"

	r, err := Generate(o)
	if err != nil {
		t.Fatalf("falha ao gerar com wallpaper: %v", err)
	}
	if !strings.Contains(r.XML, getWallpaperPath) {
		t.Errorf("falta getWallpaperPath no XML")
	}
	if !strings.Contains(r.XML, setWallpaperPath) {
		t.Errorf("falta setWallpaperPath no XML")
	}
	if !strings.Contains(r.XML, "https://marcovivas.com/docs/DEV9.png") {
		t.Errorf("falta URL da imagem no script")
	}
}

func TestChocolateyPostInstall(t *testing.T) {
	if err := loadAssets(); err != nil {
		t.Fatal(err)
	}
	o := DefaultOptions()
	o.PostEnabled = true
	o.PostChocoEnabled = true
	o.PostChocoCustom = "googlechrome\n7zip\nvlc\n# comentario\ngit"
	o.PostWingetCustom = "Google.Chrome\nVideoLAN.VLC"

	r, err := Generate(o)
	if err != nil {
		t.Fatalf("falha ao gerar com chocolatey: %v", err)
	}
	if !strings.Contains(r.XML, "community.chocolatey.org/install.ps1") {
		t.Errorf("falta script de instalação do Chocolatey")
	}
	if !strings.Contains(r.XML, "choco install 'googlechrome'") && !strings.Contains(r.XML, "choco install $pkg") {
		t.Errorf("falta comando choco install")
	}
	if !strings.Contains(r.XML, "'googlechrome'") {
		t.Errorf("falta pacote googlechrome")
	}
	if !strings.Contains(r.XML, "'7zip'") {
		t.Errorf("falta pacote 7zip")
	}
	if !strings.Contains(r.XML, "Google.Chrome") {
		t.Errorf("falta pacote Winget Google.Chrome")
	}
}

func TestXMLWellFormedWithAllFeatures(t *testing.T) {
	if err := loadAssets(); err != nil {
		t.Fatal(err)
	}
	o := DefaultOptions()
	o.CreateUser = true
	o.UserAccounts = []UserAccount{
		{Name: "Marco", DisplayName: "Marco", Password: "site@1989", Group: "Administrators"},
		{Name: "Usuario", DisplayName: "Usuario", Password: "", Group: "Users"},
	}
	o.AutoLogon = true
	o.AutoLogonUser = "Marco"
	o.AutoLogonN = 1
	o.WallpaperEnabled = true
	o.WallpaperURL = "https://marcovivas.com/docs/DEV9.png"
	o.PostEnabled = true
	o.PostAutoRun = true
	o.PostChocoEnabled = true
	o.PostChocoCustom = "googlechrome\n7zip"
	o.PostWingetCustom = "Google.Chrome\nVideoLAN.VLC"
	o.FirstLogonCmds = "echo 1 & echo 2"

	r, err := Generate(o)
	if err != nil {
		t.Fatalf("falha ao gerar XML completo: %v", err)
	}

	// Verifica se a substituição de & por &amp; ocorreu nos elementos CommandLine
	if strings.Contains(r.XML, `<CommandLine>powershell.exe -NoProfile -ExecutionPolicy Bypass -Command "try { $b = & '`) {
		t.Errorf("encontrado & sem escape no elemento CommandLine do wallpaper!")
	}
	if !strings.Contains(r.XML, `<CommandLine>powershell.exe -NoProfile -ExecutionPolicy Bypass -Command "try { $b = &amp; '`) {
		t.Errorf("falta &amp; escapado no elemento CommandLine do wallpaper!")
	}
	if !strings.Contains(r.XML, "<CommandLine>echo 1 &amp; echo 2</CommandLine>") {
		t.Errorf("falta &amp; escapado no FirstLogonCmds!")
	}
}
