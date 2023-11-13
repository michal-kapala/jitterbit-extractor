package project

import (
	"encoding/xml"
	"fmt"
	"os"
	"strings"

	"jbextractor/jitterbit/entity"
)

type EntityType struct {
	XMLName		xml.Name 	`xml:"EntityType"`
	Name			string		`xml:"name,attr"`
	Folders		[]Folder	`xml:"Folder"`
	Entities	[]Entity	`xml:"Entity"`
}

// CreateDirs performs DFS creation of ID-named directories.
func (et *EntityType) CreateDirs(dirs *map[string]string, path string, sep string) error {
	// scripts
	scriptPath := fmt.Sprintf("%s%sScript", path, sep)
	if err := os.Mkdir(scriptPath, os.ModePerm); err != nil {
		return err
	}

	for _, folder := range et.Folders {
		// save and create top folders
		folderPath := fmt.Sprintf("%s%s%s", scriptPath, sep, folder.Id)
		(*dirs)[folder.Id] = folderPath
		if err := os.Mkdir(folderPath, os.ModePerm); err != nil {
			return err
		}
		// save and create subdirectories recursively
		err := folder.createSubfolders(dirs, folderPath, sep)
		if err != nil {
			return err
		}
	}

	return nil
}

// FindEntityDir searches for entity ID in EntityType and returns its directory.
func (et *EntityType) FindEntityDir(dirs *map[string]string, id string, rootDir string) (*Entity, string) {
	var entity *Entity 
	dir := ""
	// top-level entities
	for _, ent := range (*et).Entities {
		if ent.Id == id {
			entity = &ent
			dir = rootDir
			return entity, dir
		}
	}

	for _, folder := range et.Folders {
		// first-level folder entities
		for _, ent := range folder.Entities {
			if ent.Id == id {
				entity = &ent
				dir = (*dirs)[folder.Id]
				return entity, dir
			}
		}
		// subdirectories
		entity, dir = folder.findSubfolderEntityDir(dirs, id)
		if dir == "" {
			continue
		}
		return entity, dir
	}

	return nil, dir
}

// RenameDirs substitutes folder IDs with real names.
func (et *EntityType) RenameDirs(dirs *map[string]string, path string) error {
	oldPath := ""
	newPath := ""
	for _, folder := range et.Folders {
		oldPath = (*dirs)[folder.Id]
		newPath = strings.Replace(oldPath, folder.Id, folder.Name, 1)
		err := os.Rename(oldPath, newPath)
		if err != nil {
			return err
		}
		// update map
		updateDirPaths(dirs, folder.Id, folder.Name)
		err = folder.renameDirs(dirs, newPath)
		if err != nil {
			return err
		}
	}
	return nil
}

// CreateScripts creates .jb source code files.
func (scripts *EntityType) CreateScripts(dirs *map[string]string, envPath string, targetPath string, sep string) error {
	inPath := fmt.Sprintf("%s%sData%sScript", envPath, sep, sep)
	outPath := fmt.Sprintf("%s%sScript", targetPath, sep)
	entries, err := os.ReadDir(inPath)
	if err != nil {
		return err
	}

	inFilePath := ""
	entryName := ""
	for _, entry := range entries {
		entryName = entry.Name()
		if !entry.IsDir() && strings.Contains(entryName, ".xml"){
			inFilePath = fmt.Sprintf("%s%s%s", inPath, sep, entryName)
			inFile, err := os.Open(inFilePath)
			if err != nil {
				return err
			}
			defer inFile.Close()

			script, err := entity.ParseEntity(inFilePath)
			if err != nil {
				return err
			}
			
			ent, scriptDir := scripts.FindEntityDir(dirs, script.Header.Id, outPath)

			// script was not found in project.xml
			if ent == nil || scriptDir == "" {
				return fmt.Errorf("[CreateScripts] Corrupted project.xml - script %s was not found", script.Header.Id)
			}
			// example script name from Jitterbit's demo project:
			// jb.sqlServer.table1-&gt;table2 [ETL_log]
			saneName := sanitizeScriptName(script.Header.Name)
			outFilePath := fmt.Sprintf("%s%s%s.jb", scriptDir, sep, saneName)
			outFile, err := os.Create(outFilePath)
			if err != nil {
				return err
			}

			_, err = outFile.WriteString(script.KongaString)
			if err != nil {
				return err
			}
			outFile.Close()
		}
	}

	return nil
}

// sanitizeScriptName cleanses the script names of special characters disallowed by Windows file system.
func sanitizeScriptName(name string) string {
	replacer := strings.NewReplacer(
		"<", "_",
		">", "_",
		"/", "_",
		"\\", "_",
		"?", "_",
		":", "_",
		"*", "_",
		"|", "_",
		"\"", "_")
	return replacer.Replace(name)
}
