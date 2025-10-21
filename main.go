package main

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"

	sf "github.com/creativeprojects/go-selfupdate"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

const version = "1.3.0"
const repoSlug = "mikimou/mikiho-launcher"

func main() {
	//update()
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:         "mikiho launcher",
		Width:         585,
		Height:        320,
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

func update() error {
	latest, found, err := sf.DetectLatest(context.Background(), sf.ParseSlug(repoSlug))
	if err != nil {
		return fmt.Errorf("error occurred while detecting version: %w", err)
	}
	if !found {
		return fmt.Errorf("latest version for %s/%s could not be found from github repository", runtime.GOOS, runtime.GOARCH)
	}

	if latest.LessOrEqual(version) {
		log.Printf("Current version (%s) is the latest", version)
		return nil
	}

	exe, err := sf.ExecutablePath()
	if err != nil {
		return errors.New("could not locate executable path")
	}
	if err := sf.UpdateTo(context.Background(), latest.AssetURL, latest.AssetName, exe); err != nil {
		return fmt.Errorf("error occurred while updating binary: %w", err)
	}
	log.Printf("Successfully updated to version %s", latest.Version())
	restartApp()
	return nil
}

func restartApp() {
	exe, _ := os.Executable()
	exec.Command(exe).Start()
	os.Exit(0)
}
