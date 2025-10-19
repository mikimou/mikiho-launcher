package main

import (
	"embed"
	"log"
	"os"
	"os/exec"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"

	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
)

//go:embed all:frontend/dist
var assets embed.FS

const version = "1.0.2"
const repoSlug = "mikimou/mikiho-launcher"

func main() {
	checkAndUpdate()
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:         "mikiho launcher",
		Width:         650,
		Height:        350,
		DisableResize: true,
		Frameless:     true,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 0, G: 0, B: 0, A: 0},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
		LogLevel: logger.DEBUG,
		Debug: options.Debug{
			OpenInspectorOnStartup: false,
		},
		Windows: &windows.Options{
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
		},
		Mac: &mac.Options{
			TitleBar:             mac.TitleBarHiddenInset(),
			Appearance:           mac.NSAppearanceNameDarkAqua,
			WebviewIsTransparent: true,
			WindowIsTranslucent:  false,
			About: &mac.AboutInfo{
				Title:   "Mikiho Launcher",
				Message: "Developed with love \n © 2024 Michal Hicz",
			},
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}

func checkAndUpdate() {
	v := semver.MustParse(version)
	latest, found, err := selfupdate.DetectLatest(repoSlug)
	if err != nil {
		log.Println("Update check failed:", err)
		return
	}
	if !found || latest.Version.LTE(v) {
		log.Println("Already up to date:", version)
		return
	}
	update(v)
}

func update(v semver.Version) {
	latest, err := selfupdate.UpdateSelf(v, repoSlug)
	if err != nil {
		log.Println("Binary update failed:", err)
		return
	}
	if latest.Version.Equals(v) {
		log.Println("Already latest version.")
		return
	}
	log.Printf("Updated to %s successfully.\n", latest.Version)
	restartApp()
}

func restartApp() {
	exe, _ := os.Executable()
	exec.Command(exe).Start()
	os.Exit(0)
}
