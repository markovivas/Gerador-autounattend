package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// Programas oferecidos na aba "Pós-instalação" (IDs do winget).
type WingetApp struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Cat   string `json:"cat"`
	Basic bool   `json:"basic"`
}

var WingetApps = []WingetApp{
	// navegadores
	{"Google.Chrome", "Google Chrome", "browsers", true},
	{"Mozilla.Firefox", "Mozilla Firefox", "browsers", false},
	{"Brave.Brave", "Brave", "browsers", false},
	// compactadores, mídia e utilitários
	{"7zip.7zip", "7-Zip", "tools", true},
	{"RARLab.WinRAR", "WinRAR", "tools", false},
	{"VideoLAN.VLC", "VLC media player", "media", true},
	{"CodecGuide.K-LiteCodecPack.Standard", "K-Lite Codec Pack", "media", false},
	{"OBSProject.OBSStudio", "OBS Studio", "media", false},
	{"Audacity.Audacity", "Audacity", "media", false},
	{"HandBrake.HandBrake", "HandBrake", "media", false},
	{"IrfanSkiski.IrfanView", "IrfanView", "media", false},
	{"Spotify.Spotify", "Spotify", "media", false},
	{"voidtools.Everything", "Everything (busca de arquivos)", "tools", false},
	{"ShareX.ShareX", "ShareX (capturas de tela)", "tools", false},
	{"qBittorrent.qBittorrent", "qBittorrent", "tools", false},
	{"Microsoft.PowerToys", "Microsoft PowerToys", "tools", false},
	{"Rufus.Rufus", "Rufus", "tools", false},
	{"Ventoy.Ventoy", "Ventoy", "tools", false},
	// escritório e PDF
	{"Adobe.Acrobat.Reader.64-bit", "Adobe Acrobat Reader", "office", true},
	{"SumatraPDF.SumatraPDF", "SumatraPDF", "office", false},
	{"Foxit.FoxitReader", "Foxit PDF Reader", "office", false},
	{"TheDocumentFoundation.LibreOffice", "LibreOffice", "office", false},
	{"ONLYOFFICE.DesktopEditors", "ONLYOFFICE", "office", false},
	{"Google.GoogleDrive", "Google Drive", "office", false},
	// comunicação
	{"WhatsApp.WhatsApp", "WhatsApp", "comm", true},
	{"Telegram.TelegramDesktop", "Telegram", "comm", false},
	{"Discord.Discord", "Discord", "comm", false},
	{"Zoom.Zoom", "Zoom", "comm", false},
	// acesso remoto
	{"AnyDeskSoftwareGmbH.AnyDesk", "AnyDesk", "remote", true},
	{"TeamViewer.TeamViewer", "TeamViewer", "remote", false},
	{"RustDesk.RustDesk", "RustDesk", "remote", false},
	// desenvolvimento
	{"Notepad++.Notepad++", "Notepad++", "dev", true},
	{"Microsoft.VisualStudioCode", "Visual Studio Code", "dev", false},
	{"Git.Git", "Git", "dev", false},
	{"Python.Python.3.12", "Python 3.12", "dev", false},
	{"Microsoft.WindowsTerminal", "Windows Terminal", "dev", false},
	// jogos
	{"Valve.Steam", "Steam", "games", false},
	{"EpicGames.EpicGamesLauncher", "Epic Games Launcher", "games", false},
	// bibliotecas e runtimes
	{"Microsoft.VCRedist.2015+.x64", "Visual C++ Redistributable 2015-2022 (x64)", "runtimes", true},
	{"Microsoft.VCRedist.2015+.x86", "Visual C++ Redistributable 2015-2022 (x86)", "runtimes", true},
	{"Microsoft.DotNet.DesktopRuntime.8", ".NET Desktop Runtime 8", "runtimes", true},
	{"Microsoft.DotNet.DesktopRuntime.6", ".NET Desktop Runtime 6", "runtimes", false},
	{"Microsoft.DirectX", "DirectX (runtime)", "runtimes", false},
	{"Oracle.JavaRuntimeEnvironment", "Java Runtime Environment", "runtimes", false},
	{"Microsoft.EdgeWebView2Runtime", "WebView2 Runtime", "runtimes", false},
	// hardware e segurança
	{"CPUID.CPU-Z", "CPU-Z", "hw", false},
	{"REALiX.HWiNFO", "HWiNFO", "hw", false},
	{"CrystalDewWorld.CrystalDiskInfo", "CrystalDiskInfo", "hw", false},
	{"Malwarebytes.Malwarebytes", "Malwarebytes", "hw", false},
}

func BasicPostApps() []string {
	var r []string
	for _, a := range WingetApps {
		if a.Basic {
			r = append(r, a.ID)
		}
	}
	return r
}

var (
	reWingetID = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._+\-]{1,99}$`)
	reChocoPkg = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._+\-]{0,99}$`)
)

const (
	postBatPath = `C:\Users\Public\Desktop\Pos-instalacao.bat`
	postPs1Path = `C:\ProgramData\Winhance\Unattend\Scripts\PosInstalacao.ps1`
)

type postApp struct{ ID, Name string }

// resolve a seleção: programas winget marcados ou digitados + pacotes chocolatey
func postSelection(o *Options) ([]postApp, []string, error) {
	sel := map[string]bool{}
	for _, id := range o.PostApps {
		sel[id] = true
	}
	var wingetApps []postApp
	seenW := map[string]bool{}
	for _, a := range WingetApps {
		if sel[a.ID] {
			wingetApps = append(wingetApps, postApp{a.ID, a.Name})
			seenW[a.ID] = true
		}
	}
	wingetSrc := o.PostWingetCustom
	if strings.TrimSpace(wingetSrc) == "" && strings.TrimSpace(o.PostCustom) != "" {
		wingetSrc = o.PostCustom
	}
	for _, ln := range strings.Split(strings.ReplaceAll(wingetSrc, "\r", ""), "\n") {
		if i := strings.Index(ln, "#"); i >= 0 {
			ln = ln[:i]
		}
		ln = strings.TrimSpace(ln)
		if ln == "" {
			continue
		}
		if !reWingetID.MatchString(ln) {
			return nil, nil, fmt.Errorf(tr(o.Lang, "ID do winget inválido: %q (use o formato Empresa.Programa)"), ln)
		}
		if !seenW[ln] {
			seenW[ln] = true
			wingetApps = append(wingetApps, postApp{ln, ln})
		}
	}

	var chocoApps []string
	seenC := map[string]bool{}
	for _, ln := range strings.Split(strings.ReplaceAll(o.PostChocoCustom, "\r", ""), "\n") {
		if i := strings.Index(ln, "#"); i >= 0 {
			ln = ln[:i]
		}
		ln = strings.TrimSpace(ln)
		if ln == "" {
			continue
		}
		if !reChocoPkg.MatchString(ln) {
			return nil, nil, fmt.Errorf(tr(o.Lang, "Pacote do Chocolatey inválido: %q"), ln)
		}
		if !seenC[strings.ToLower(ln)] {
			seenC[strings.ToLower(ln)] = true
			chocoApps = append(chocoApps, ln)
		}
	}

	return wingetApps, chocoApps, nil
}

type postText struct {
	Title, NoAdmin, Count, WaitWinget, NoWinget, Enter, Done, Failed, Retry string
}

var postTexts = map[string]postText{
	"pt-BR": {"Pós-instalação — instalando programas", "Não foi possível pedir permissão de administrador; algumas instalações podem falhar.", "Serão instalados {0} programas via Winget. Aguarde, isso pode demorar.", "Aguardando o winget (Instalador de Aplicativo) ficar disponível… precisa de internet.", "O winget não está disponível neste Windows. Atualize o \"Instalador de Aplicativo\" na Microsoft Store e execute este script de novo.", "Pressione Enter para fechar", "Concluído: {0} instalado(s) ou já presente(s), {1} com falha.", "Programas que falharam:", "Verifique a internet e execute o script novamente: ele só instala o que falta."},
	"en":    {"Post-install — installing programs", "Could not ask for administrator permission; some installs may fail.", "{0} programs will be installed via Winget. Please wait, this can take a while.", "Waiting for winget (App Installer) to become available… internet is required.", "winget is not available on this Windows. Update \"App Installer\" in the Microsoft Store and run this script again.", "Press Enter to close", "Done: {0} installed or already present, {1} failed.", "Programs that failed:", "Check your internet connection and run the script again: it only installs what is missing."},
	"es":    {"Posinstalación — instalando programas", "No se pudo pedir permiso de administrador; algunas instalaciones pueden fallar.", "Se instalarán {0} programas con Winget. Espere, puede tardar.", "Esperando a que winget (Instalador de aplicación) esté disponible… se necesita internet.", "winget no está disponible en este Windows. Actualice el \"Instalador de aplicación\" en Microsoft Store y ejecute este script de nuevo.", "Pulse Enter para cerrar", "Terminado: {0} instalado(s) o ya presente(s), {1} con error.", "Programas con error:", "Compruebe internet y ejecute el script de nuevo: solo instala lo que falta."},
}

func psq(s string) string { return "'" + strings.ReplaceAll(s, "'", "''") + "'" }

func buildPostScript(o *Options, apps []postApp, chocoApps []string) string {
	tx, ok := postTexts[o.Lang]
	if !ok {
		tx = postTexts["pt-BR"]
	}
	q := func(s string) string { return strings.ReplaceAll(s, "'", "''") }

	var wingetBlock string
	if len(apps) > 0 {
		var al []string
		for _, a := range apps {
			al = append(al, fmt.Sprintf("    @{ Id = %s; Name = %s }", psq(a.ID), psq(a.Name)))
		}
		wingetBlock = fmt.Sprintf(`$apps = @(
%s
)

Write-Host ''
Write-Host '[Winget] %s' -ForegroundColor Cyan
Write-Host ('%s' -f $apps.Count)
Write-Host ''

function Find-Winget { Get-Command winget.exe -ErrorAction SilentlyContinue }
$tries = 0
while (-not (Find-Winget) -and $tries -lt 20) {
    if ($tries -eq 0) { Write-Host '%s' -ForegroundColor Yellow }
    try { Add-AppxPackage -RegisterByFamilyName -MainPackage 'Microsoft.DesktopAppInstaller_8wekyb3d8bbwe' -ErrorAction Stop } catch { }
    Start-Sleep -Seconds 15
    $tries++
}
if (-not (Find-Winget)) {
    Write-Host '%s' -ForegroundColor Red
} else {
    $i = 0
    foreach ($a in $apps) {
        $i++
        Write-Host ('[Winget {0}/{1}] {2}' -f $i, $apps.Count, $a.Name) -ForegroundColor Cyan
        $code = 1
        for ($try = 1; $try -le 2 -and $code -ne 0; $try++) {
            & winget install --id $a.Id -e --source winget --silent --accept-package-agreements --accept-source-agreements
            $code = $LASTEXITCODE
            if ($code -eq -1978335135 -or $code -eq -1978335189) { $code = 0 }
            if ($code -ne 0 -and $try -lt 2) { Start-Sleep -Seconds 10 }
        }
        if ($code -eq 0) { $wOk += $a.Name } else { $wFail += $a.Name }
    }
    Write-Host ''
    Write-Host ('Winget: {0} instalado(s), {1} com falha.' -f $wOk.Count, $wFail.Count) -ForegroundColor Green
    if ($wFail.Count -gt 0) {
        Write-Host 'Pacotes do Winget que falharam:' -ForegroundColor Yellow
        foreach ($f in $wFail) { Write-Host (' - ' + $f) -ForegroundColor Yellow }
    }
}
`, strings.Join(al, "\n"), q(tx.Title), q(tx.Count), q(tx.WaitWinget), q(tx.NoWinget))
	}

	var chocoBlock string
	if o.PostChocoEnabled || len(chocoApps) > 0 {
		var cl []string
		for _, c := range chocoApps {
			cl = append(cl, fmt.Sprintf("    %s", psq(c)))
		}
		chocoBlock = fmt.Sprintf(`Write-Host ''
Write-Host '[Chocolatey] Gerenciador de Pacotes' -ForegroundColor Cyan
function Find-Choco { Get-Command choco.exe -ErrorAction SilentlyContinue }
if (-not (Find-Choco)) {
    Write-Host 'Instalando Chocolatey...' -ForegroundColor Cyan
    try {
        Set-ExecutionPolicy Bypass -Scope Process -Force
        [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072
        iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))
    } catch {
        Write-Host "Falha ao instalar Chocolatey: $($_.Exception.Message)" -ForegroundColor Red
    }
    $cElapsed = 0
    while (-not (Find-Choco) -and $cElapsed -lt 60) {
        Start-Sleep -Seconds 3
        $cElapsed += 3
        $env:Path = [System.Environment]::GetEnvironmentVariable("Path","Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path","User")
    }
}
if (Find-Choco) {
    Write-Host 'Chocolatey pronto.' -ForegroundColor Green
    choco feature enable -n allowGlobalConfirmation 2>$null
    $chocoPkgs = @(
%s
    )
    if ($chocoPkgs.Count -gt 0) {
        $cIdx = 0
        foreach ($pkg in $chocoPkgs) {
            $cIdx++
            Write-Host ('[Chocolatey {0}/{1}] Instalando {2}...' -f $cIdx, $chocoPkgs.Count, $pkg) -ForegroundColor Cyan
            & choco install $pkg -y --no-progress
            if ($LASTEXITCODE -eq 0) {
                $cOk += $pkg
                Write-Host "$pkg instalado com sucesso." -ForegroundColor Green
            } else {
                $cFail += $pkg
                Write-Host "$pkg falhou (código $LASTEXITCODE)." -ForegroundColor Yellow
            }
        }
        Write-Host ''
        Write-Host ('Chocolatey: {0} instalado(s), {1} com falha.' -f $cOk.Count, $cFail.Count) -ForegroundColor Green
        if ($cFail.Count -gt 0) {
            Write-Host 'Pacotes do Chocolatey que falharam:' -ForegroundColor Yellow
            foreach ($f in $cFail) { Write-Host (' - ' + $f) -ForegroundColor Yellow }
        }
    }
} else {
    Write-Host 'Chocolatey não pôde ser instalado (verifique sua conexão com a internet).' -ForegroundColor Red
}
`, strings.Join(cl, "\n"))
	}

	del := ""
	if o.PostDelete {
		del = "if (($wFail.Count -eq 0) -and ($cFail.Count -eq 0)) { Remove-Item -LiteralPath '" + postBatPath + "' -Force -ErrorAction SilentlyContinue; Remove-Item -LiteralPath $PSCommandPath -Force -ErrorAction SilentlyContinue }\n"
	}

	res := fmt.Sprintf(`$ErrorActionPreference = 'Continue'
try { [Console]::OutputEncoding = [System.Text.Encoding]::UTF8 } catch { }
$Host.UI.RawUI.WindowTitle = '%s'

$isAdmin = ([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
if (-not $isAdmin) {
    try {
        Start-Process -FilePath 'powershell.exe' -Verb RunAs -ArgumentList @('-NoProfile', '-ExecutionPolicy', 'Bypass', '-File', ('"{0}"' -f $PSCommandPath))
        exit
    } catch {
        Write-Host '%s' -ForegroundColor Yellow
    }
}

$wOk = @()
$wFail = @()
$cOk = @()
$cFail = @()

%s
%s

Write-Host ''
Write-Host '========================================' -ForegroundColor Green
Write-Host 'Pós-instalação concluída.' -ForegroundColor Green
Write-Host '========================================' -ForegroundColor Green
%sRead-Host '%s' | Out-Null
`, q(tx.Title), q(tx.NoAdmin), wingetBlock, chocoBlock, del, q(tx.Enter))

	return strings.ReplaceAll(res, "\r", "")
}

const postBat = "@echo off\npowershell.exe -NoProfile -ExecutionPolicy Bypass -File \"" + postPs1Path + "\"\n"

// trechos <File> adicionais do bloco Extensions
func postFiles(o *Options, apps []postApp, chocoApps []string) string {
	ps := strings.ReplaceAll(buildPostScript(o, apps, chocoApps), "]]>", "]]]]><![CDATA[>")
	return "    <File path=\"" + postPs1Path + "\"><![CDATA[\n" + ps + "]]></File>\n" +
		"    <File path=\"" + postBatPath + "\"><![CDATA[\n" + postBat + "]]></File>\n"
}

var errNoApps = errors.New("no apps")
