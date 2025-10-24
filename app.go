package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/hashicorp/go-getter"
	wailsrt "github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

const modpackName = "geccopack"
const manifestUrl = "https://archive.hicz.net/geccopack/manifest.json"

type ModpackManifest struct {
	Version string   `json:"version"`
	URL     string   `json:"url"`
	Command []string `json:"command,omitempty"`
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) Version() string {
	return version
}

type Options struct {
	Nickname string `json:"nickname"`
	Ram      int    `json:"ram"`
}

func (a *App) SaveOptions(nickname string, ram int) error {
	opts := Options{
		Nickname: nickname,
		Ram:      ram,
	}

	appData, err := os.UserConfigDir()
	if err != nil {
		return fmt.Errorf("failed to get config dir: %w", err)
	}

	configDir := filepath.Join(appData, "mikiho-launcher")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config dir: %w", err)
	}

	optionsFile := filepath.Join(configDir, "options.json")
	data, err := json.Marshal(opts)
	if err != nil {
		return fmt.Errorf("failed to marshal options: %w", err)
	}

	if err := os.WriteFile(optionsFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write options file: %w", err)
	}

	return nil
}

func (a *App) LoadOptions() (*Options, error) {
	appData, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get config dir: %w", err)
	}

	optionsFile := filepath.Join(appData, "mikiho-launcher", "options.json")
	data, err := os.ReadFile(optionsFile)
	if err != nil {
		if os.IsNotExist(err) {
			return &Options{}, nil
		}
		return nil, fmt.Errorf("failed to read options file: %w", err)
	}

	var opts Options
	if err := json.Unmarshal(data, &opts); err != nil {
		return nil, fmt.Errorf("failed to parse options file: %w", err)
	}

	return &opts, nil
}

func (a *App) UpdateModpack(manifest *ModpackManifest) error {
	appData, _ := os.UserConfigDir()
	modpackDir := filepath.Join(appData, "mikiho-launcher", modpackName)

	// Create directories
	if err := os.MkdirAll(modpackDir, 0755); err != nil {
		return fmt.Errorf("failed to create modpack dir: %w", err)
	}

	// // Download modpack
	// resp, err := http.Get(manifest.URL)
	// if err != nil {
	// 	return fmt.Errorf("failed to download modpack: %w", err)
	// }
	// defer resp.Body.Close()

	// // Save modpack zip
	// zipPath := filepath.Join(modpackDir, "modpack.zip")
	// out, err := os.Create(zipPath)
	// if err != nil {
	// 	return fmt.Errorf("failed to create zip file: %w", err)
	// }
	// defer out.Close()

	// if _, err := io.Copy(out, resp.Body); err != nil {
	// 	return fmt.Errorf("failed to save modpack: %w", err)
	// }
	client := getter.Client{DisableSymlinks: true}
	client.Dst = modpackDir
	client.Dir = true
	client.Src = manifest.URL
	err := client.Get()
	if err != nil {
		fmt.Println("Error")
	}
	if err == nil {
		fmt.Println("Finished")
	}

	// Save manifest
	manifestData, err := json.Marshal(manifest)
	if err != nil {
		return fmt.Errorf("failed to marshal manifest: %w", err)
	}

	manifestPath := filepath.Join(modpackDir, "manifest.json")
	if err := os.WriteFile(manifestPath, manifestData, 0644); err != nil {
		return fmt.Errorf("failed to save manifest: %w", err)
	}

	return nil
}

func (a *App) CheckModpack() error {
	// Read local manifest if exists
	appData, _ := os.UserConfigDir()
	localManifestPath := filepath.Join(appData, "mikiho-launcher", modpackName, "manifest.json")

	var localVersion string
	localManifest, err := os.ReadFile(localManifestPath)
	if err == nil {
		var m ModpackManifest
		if err := json.Unmarshal(localManifest, &m); err == nil {
			localVersion = m.Version
		}
	}

	// Fetch remote manifest
	resp, err := http.Get(manifestUrl)
	if err != nil {
		return fmt.Errorf("failed to fetch manifest: %w", err)
	}
	defer resp.Body.Close()

	var remoteManifest ModpackManifest
	if err := json.NewDecoder(resp.Body).Decode(&remoteManifest); err != nil {
		return fmt.Errorf("failed to parse manifest: %w", err)
	}

	if remoteManifest.Version == localVersion {
		return nil // no update needed
	} else {
		if err := a.UpdateModpack(&remoteManifest); err != nil {
			return fmt.Errorf("failed to update modpack: %w", err)
		}
	}

	return nil
}

func (a *App) LaunchGame(nickname string, ram int) error {
	appData, err := os.UserConfigDir()
	if err != nil {
		return fmt.Errorf("failed to get config dir: %w", err)
	}
	modpackDir := filepath.Join(appData, "mikiho-launcher", modpackName)
	manifestPath := filepath.Join(modpackDir, "manifest.json")

	b, err := os.ReadFile(manifestPath)
	if err != nil {
		return fmt.Errorf("failed to read manifest: %w", err)
	}
	var manifest ModpackManifest
	if err := json.Unmarshal(b, &manifest); err != nil {
		return fmt.Errorf("failed to parse manifest: %w", err)
	}
	if len(manifest.Command) == 0 {
		return fmt.Errorf("manifest has no command template")
	}

	// substitute placeholders
	args := make([]string, len(manifest.Command))
	for i, p := range manifest.Command {
		p = strings.ReplaceAll(p, "{{nick}}", nickname)
		p = strings.ReplaceAll(p, "{{ram}}", fmt.Sprintf("%d", ram))
		p = strings.ReplaceAll(p, "{{mcdir}}", modpackDir)
		// allow manifest to include {{java}} if you want to override java path
		javaPath := filepath.Join(modpackDir, "runtime", "java-runtime-gamma", "windows-x64", "java-runtime-gamma", "bin", "java.exe")
		p = strings.ReplaceAll(p, "{{java}}", javaPath)
		args[i] = p
	}

	exe := args[0]
	execArgs := args[1:]

	// cmdStr := strings.Join(args, " ")
	// wailsrt.LogInfo(a.ctx, "Launching with command: "+cmdStr)

	cmd := exec.CommandContext(context.Background(), exe, execArgs...)
	cmd.Dir = modpackDir
	// stdout, _ := cmd.StdoutPipe()
	// stderr, _ := cmd.StderrPipe()

	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start game: %w", err)
	}

	// // stream logs to frontend
	// go streamAndEmit(a.ctx, stdout)
	// go streamAndEmit(a.ctx, stderr)

	// // wait in background and emit exit event
	// go func() {
	// 	_ = cmd.Wait()
	// 	wailsrt.LogInfo(a.ctx, "Game exited")
	// }()

	return nil
}

func streamAndEmit(ctx context.Context, r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		wailsrt.LogInfo(ctx, scanner.Text()) // Use LogInfo instead of EventsEmit for logs
	}
}
