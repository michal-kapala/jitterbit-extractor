package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	jbproj "jbextractor/jitterbit/project"
	jbscript "jbextractor/jitterbit/script"

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
		runtime.LogPrint(a.ctx, err.Error())
		return ""
	}

	if directory == "" {
		runtime.LogPrint(a.ctx, "[SelectProject] User cancelled choosing a directory")
		return directory
	}

	// validate project dir
	file, err := os.Open(fmt.Sprintf("%s\\manifest.jip", directory))

	if err != nil {
		runtime.LogPrint(a.ctx, err.Error())
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
		runtime.LogPrint(a.ctx, err.Error())
		return ""
	}

	return directory
}

// GetEnvs returns the available local environments.
func (a *App) GetEnvs(project string) []string {
	if project == "" {
		runtime.LogPrint(a.ctx, "[GetEnvs] Empty project path")
		return nil
	}

	files, err := os.ReadDir(project)

	if err != nil {
		runtime.LogPrint(a.ctx, err.Error())
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

// Extract performs convertions from Jitterbit Studio format to a more readable project structure.
func (a *App) Extract(projectPath string, env string, output string) bool {
	// get project name
	manifest, err := os.Open(fmt.Sprintf("%s\\manifest.jip", projectPath))
	if err != nil {
		runtime.LogPrint(a.ctx, err.Error())
		return false
	}
	defer manifest.Close()

	projectName := ""
	scanner := bufio.NewScanner(manifest)
	manifestContent := ""
	for scanner.Scan() {
		line := scanner.Text()
		manifestContent += fmt.Sprintf("%s\r\n", line)
		if strings.Contains(line, "project-name=") {
			projectName = strings.Replace(line, "project-name=", "", 1)
		}
	}

	// corrupted manifest, use dir name instead
	if projectName == "" {
		runtime.LogPrint(a.ctx, "[Extract] Corrupted manifest.jip")
		projectName = filepath.Base(filepath.Dir(projectPath))
	}

	// get environment name
	envPath := fmt.Sprintf("%s\\%s\\environment.properties", projectPath, env)
	envProps, err := os.Open(envPath)
	if err != nil {
		runtime.LogPrint(a.ctx, err.Error())
		return false
	}
	defer envProps.Close()

	envName := ""
	scanner = bufio.NewScanner(envProps)
	envPropsContent := ""
	for scanner.Scan() {
		line := scanner.Text()
		envPropsContent += fmt.Sprintf("%s\r\n", line)
		if strings.Contains(line, "environment-name=") {
			envName = strings.Replace(line, "environment-name=", "", 1)
		}
	}

	// corrupted properties, use dir name instead
	if envName == "" {
		runtime.LogPrint(a.ctx, "[Extract] Corrupted environment.properties")
		envName = env
	}

	targetDirName := fmt.Sprintf("%s %s", projectName, envName)
	targetPath := fmt.Sprintf("%s\\%s", output, targetDirName)
	if err := os.Mkdir(targetPath, os.ModePerm); err != nil {
		runtime.LogPrint(a.ctx, err.Error())
		return false
	}

	// copy metadata
	manifestCopy, err := os.Create(fmt.Sprintf("%s\\project.properties", targetPath))
	if err != nil {
		runtime.LogPrint(a.ctx, err.Error())
		return false
	}
	defer manifestCopy.Close()
	manifestCopy.WriteString(manifestContent)

	envPropsCopy, err := os.Create(fmt.Sprintf("%s\\environment.properties", targetPath))
	if err != nil {
		runtime.LogPrint(a.ctx, err.Error())
		return false
	}
	defer envPropsCopy.Close()
	envPropsCopy.WriteString(envPropsContent)

	envPath = fmt.Sprintf("%s\\%s", projectPath, env)
	project, err := jbproj.ParseProject(envPath)
	if err != nil {
		runtime.LogPrint(a.ctx, err.Error())
		return false
	}

	// Scripts
	scripts := jbproj.GetEntityType(project, "Script")
	scriptDirs := make(map[string]string)

	err = jbproj.CreateDirs(scripts, &scriptDirs, targetPath)
	if err != nil {
		runtime.LogPrint(a.ctx, err.Error())
		return false
	}

	err = jbscript.CreateScripts(scripts, &scriptDirs, envPath)
	if err != nil {
		runtime.LogPrint(a.ctx, err.Error())
		return false
	}
	
	return true
}
