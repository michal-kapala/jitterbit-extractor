package main

import (
	"context"
	"fmt"
	"os"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called at application startup
func (a *App) startup(ctx context.Context) {
	// Perform your setup here
	a.ctx = ctx
	runtime.WindowMaximise(a.ctx)
}

// domReady is called after the front-end dom has been loaded
func (a App) domReady(ctx context.Context) {
	// Add your action here
}

// shutdown is called at application termination
func (a *App) shutdown(ctx context.Context) {
	// Perform your teardown here
}

// SelectProject displays a Jitterbit project directory dialog.
func (a *App) SelectProject() string {
	options := runtime.OpenDialogOptions{
		Title:                "Select Jitterbit project",
		CanCreateDirectories: false,
	}
	directory, err := runtime.OpenDirectoryDialog(a.ctx, options)

	if err != nil {
		fmt.Print(err.Error())
		return ""
	}

	// cancelled
	if directory == "" {
		return directory
	}

	// validate project dir
	file, err := os.Open(fmt.Sprintf("%s\\manifest.jip", directory))

	if err != nil {
		fmt.Print(err.Error())
		return ""
	}

	defer file.Close()

	return directory
}

// SelectOutput displays a target directory dialog.
func (a *App) SelectOutput() string {
	options := runtime.OpenDialogOptions{
		Title:                "Select target directory",
		CanCreateDirectories: true,
	}
	directory, err := runtime.OpenDirectoryDialog(a.ctx, options)

	if err != nil {
		fmt.Print(err.Error())
		return ""
	}

	return directory
}

// GetEnvs returns the available local environments.
func (a *App) GetEnvs(project string) []string {
	if project == "" {
		return nil
	}

	files, err := os.ReadDir(project)

	if err != nil {
		fmt.Print(err.Error())
		return nil
	}

	envs := []string{}
	for _, entry := range files {
		if entry.IsDir() {
			envs = append(envs, entry.Name())
		}
	}

	return envs
}
