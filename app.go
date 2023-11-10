package main

import (
	"bufio"
	"context"
	"fmt"
	jbproj "jbextractor/jitterbit/project"
	jbscript "jbextractor/jitterbit/script"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	// Application context.
	ctx context.Context
	// Operating system.
	os string
	// Platform-based path separator character.
	pathSep string
	// Platform-based End-Of-Line character(s).
	eol string
}

// NewApp creates a new App application struct
func NewApp(os string) *App {
	sep := ""
	eol := ""
	if os == "windows" {
		sep = "\\"
		eol = "\r\n"
	} else {
		sep = "/"
		eol = "\n"
	}
	return &App{
		os: os,
		pathSep: sep,
		eol: eol,
	}
}

// startup is called at application startup
func (a *App) startup(ctx context.Context) {
	// Perform your setup here
	a.ctx = ctx
	runtime.LogPrintf(a.ctx, "OS: %s", a.os)
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
	file, err := os.Open(fmt.Sprintf("%s%smanifest.jip", directory, a.pathSep))

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
	manifest, err := os.Open(fmt.Sprintf("%s%smanifest.jip", projectPath, a.pathSep))
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
		manifestContent += fmt.Sprintf("%s%s", line, a.eol)
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
	envPath := fmt.Sprintf("%s%s%s%senvironment.properties", projectPath, a.pathSep, env, a.pathSep)
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
		envPropsContent += fmt.Sprintf("%s%s", line, a.eol)
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
	targetPath := fmt.Sprintf("%s%s%s", output, a.pathSep, targetDirName)
	if err := os.Mkdir(targetPath, os.ModePerm); err != nil {
		targetPath += fmt.Sprintf(" %s", getDate())
		if err := os.Mkdir(targetPath, os.ModePerm); err != nil {
			runtime.LogPrint(a.ctx, err.Error())
			return false
		}
	}

	// copy metadata
	manifestCopy, err := os.Create(fmt.Sprintf("%s%sproject.properties", targetPath, a.pathSep))
	if err != nil {
		runtime.LogPrint(a.ctx, err.Error())
		return false
	}
	defer manifestCopy.Close()
	manifestCopy.WriteString(manifestContent)

	envPropsCopy, err := os.Create(fmt.Sprintf("%s%senvironment.properties", targetPath, a.pathSep))
	if err != nil {
		runtime.LogPrint(a.ctx, err.Error())
		return false
	}
	defer envPropsCopy.Close()
	envPropsCopy.WriteString(envPropsContent)

	envPath = fmt.Sprintf("%s%s%s", projectPath, a.pathSep, env)
	project, err := jbproj.ParseProject(envPath, a.pathSep)
	if err != nil {
		runtime.LogPrint(a.ctx, err.Error())
		return false
	}

	// Scripts
	scripts := jbproj.GetEntityType(project, "Script")
	scriptDirs := make(map[string]string)

	err = jbproj.CreateDirs(scripts, &scriptDirs, targetPath, a.pathSep)
	if err != nil {
		runtime.LogPrint(a.ctx, err.Error())
		return false
	}

	err = jbscript.CreateScripts(scripts, &scriptDirs, envPath, targetPath, a.pathSep)
	if err != nil {
		runtime.LogPrint(a.ctx, err.Error())
		return false
	}

	err = jbproj.RenameDirs(scripts, &scriptDirs, targetPath)
	if err != nil {
		runtime.LogPrint(a.ctx, err.Error())
		return false
	}
	
	return true
}

// getDate returns a custom time suffix for files and directories.
func getDate() string {
	date := time.Now().Format("2006-01-02 15:04:05")
	replacer := strings.NewReplacer(
		"-", "",
		" ", "_",
		":", "")
	return replacer.Replace(date)
}
