package main

import (
	"bufio"
	"context"
	"fmt"
	jbproj "jbextractor/jitterbit/project"
	"os"
	"path/filepath"
	"regexp"
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
		os:      os,
		pathSep: sep,
		eol:     eol,
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
		a.logError(err)
		return ""
	}

	if directory == "" {
		a.logWarning("[SelectProject] User cancelled choosing a directory")
		return directory
	}

	// validate project dir
	file, err := os.Open(fmt.Sprintf("%s%smanifest.jip", directory, a.pathSep))

	if err != nil {
		a.logError(err)
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
		a.logError(err)
		return ""
	}

	return directory
}

// GetEnvs returns the available local environments.
func (a *App) GetEnvs(project string) []string {
	if project == "" {
		a.logWarning("[GetEnvs] Empty project path")
		return nil
	}

	files, err := os.ReadDir(project)

	if err != nil {
		a.logError(err)
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
	targetPath, err := a.copyMetadata(projectPath, env, output)
	if err != nil {
		a.logError(err)
		return false
	}

	envPath := fmt.Sprintf("%s%s%s", projectPath, a.pathSep, env)
	project, err := jbproj.ParseProject(envPath, a.pathSep)
	if err != nil {
		a.logError(err)
		return false
	}

	// Operations
	ops := project.GetEntityType(jbproj.OPERATION)
	err = ops.CreateDirs(targetPath, a.pathSep)
	if err != nil {
		a.logError(err)
		return false
	}

	err = ops.CreateOperations(envPath, targetPath, a.pathSep)
	if err != nil {
		a.logError(err)
		return false
	}

	err = ops.RenameDirs(targetPath)
	if err != nil {
		a.logError(err)
		return false
	}

	// Scripts
	scripts := project.GetEntityType(jbproj.SCRIPT)

	err = scripts.CreateDirs(targetPath, a.pathSep)
	if err != nil {
		a.logError(err)
		return false
	}

	err = scripts.CreateScripts(envPath, targetPath, a.pathSep)
	if err != nil {
		a.logError(err)
		return false
	}

	err = scripts.RenameDirs(targetPath)
	if err != nil {
		a.logError(err)
		return false
	}

	err = a.resolveScripts(project, targetPath)
	if err != nil {
		a.logError(err)
		return false
	}

	return true
}

// Copies over the project and environment metadata files.
func (a *App) copyMetadata(projectPath string, env string, output string) (string, error) {
	// get project name
	manifest, err := os.Open(fmt.Sprintf("%s%smanifest.jip", projectPath, a.pathSep))
	if err != nil {
		return "", err
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
		a.logWarning("[Extract] Corrupted manifest.jip")
		projectName = filepath.Base(filepath.Dir(projectPath))
	}

	// get environment name
	envPath := fmt.Sprintf("%s%s%s%senvironment.properties", projectPath, a.pathSep, env, a.pathSep)
	envProps, err := os.Open(envPath)
	if err != nil {
		return "", err
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
			return targetPath, err
		}
	}

	manifestCopy, err := os.Create(fmt.Sprintf("%s%sproject.properties", targetPath, a.pathSep))
	if err != nil {
		return targetPath, err
	}
	defer manifestCopy.Close()
	manifestCopy.WriteString(manifestContent)

	envPropsCopy, err := os.Create(fmt.Sprintf("%s%senvironment.properties", targetPath, a.pathSep))
	if err != nil {
		return targetPath, err
	}
	defer envPropsCopy.Close()
	envPropsCopy.WriteString(envPropsContent)
	return targetPath, nil
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

func (a *App) resolveScripts(project *jbproj.Project, rootPath string) error {
	scripts := project.GetEntityType(jbproj.SCRIPT)
	ops := project.GetEntityType(jbproj.OPERATION)

	return filepath.WalkDir(rootPath,
		func(path string, d os.DirEntry, err error) error {
			if !d.IsDir() {
				data, err := os.ReadFile(path)
				if err != nil {
					return err
				}

				script := string(data)
				// RunScript
				regex := regexp.MustCompile(`RunScript\(\"sc\.([0-9a-f]{8}\-[0-9a-f]{4}\-[0-9a-f]{4}\-[0-9a-f]{4}\-[0-9a-f]{12})\".*`)
				matches := regex.FindAllStringSubmatch(script, -1)
				for _, m := range matches {
					ent, dir := scripts.FindEntity(m[1], fmt.Sprintf("%s%sData%s%s", project.EnvPath, a.pathSep, a.pathSep, scripts.Type))
					if ent == nil || dir == "" {
						a.logWarning(fmt.Sprintf("[ResolveScripts] Script %s could not be found", m[1]))
						continue
					}
					cbPath := makeCallablePath(ent, dir, rootPath, scripts.Type, a.pathSep)
					replacement := strings.Replace(m[0], fmt.Sprintf("sc.%s", m[1]), cbPath, 1)
					script = strings.Replace(script, m[0], replacement, 1)
				}

				// RunOperation
				regex = regexp.MustCompile(`RunOperation\(\"op\.([0-9a-f]{8}\-[0-9a-f]{4}\-[0-9a-f]{4}\-[0-9a-f]{4}\-[0-9a-f]{12})\".*`)
				matches = regex.FindAllStringSubmatch(script, -1)
				for _, m := range matches {
					ent, dir := ops.FindEntity(m[1], fmt.Sprintf("%s%sData%s%s", project.EnvPath, a.pathSep, a.pathSep, ops.Type))
					if ent == nil || dir == "" {
						a.logWarning(fmt.Sprintf("[ResolveScripts] Operation %s could not be found", m[1]))
						continue
					}
					cbPath := makeCallablePath(ent, dir, rootPath, ops.Type, a.pathSep)
					replacement := strings.Replace(m[0], fmt.Sprintf("op.%s", m[1]), cbPath, 1)
					script = strings.Replace(script, m[0], replacement, 1)
				}

				// JavaScript tags
				regex = regexp.MustCompile(`^<javascript>\n(.|[\r|\n])*\n</javascript>\z`)
				jsMatch := regex.FindStringSubmatch(script)

				if jsMatch != nil {
					err = os.Remove(path)
					if err != nil {
						return err
					}
					jsMatch[0] = strings.TrimPrefix(jsMatch[0], "<javascript>\n")
					script = strings.TrimSuffix(jsMatch[0], "\n</javascript>")
					path = fmt.Sprintf("%s%s", strings.TrimSuffix(path, ".jb"), ".js")
				}

				file, err := os.Create(path)
				if err != nil {
					return err
				}

				_, err = file.WriteString(script)

				if err != nil {
					return err
				}

				err = file.Close()
				if err != nil {
					return err
				}
			}

			return nil
		},
	)
}

func makeCallablePath(ent *jbproj.Entity, dir string, rootPath string, typeName string, sep string) string {
	basePath := fmt.Sprintf("%s%s%s", rootPath, sep, typeName)
	path := fmt.Sprintf("%s%s%s", dir, sep, ent.Name)
	// diff the paths
	path = strings.Replace(path, basePath, "", 1)
	// truncate the first separator
	path = strings.Replace(path, sep, "", 1)
	// normalize separators
	if sep == "\\" {
		path = strings.Replace(path, sep, "/", -1)
	}

	return fmt.Sprintf("<TAG>%ss/%s</TAG>", typeName, path)
}
