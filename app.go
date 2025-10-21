package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
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
}

func (a *App) SaveOptions(nickname string) error {
	opts := Options{
		Nickname: nickname,
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
