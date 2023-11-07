package project

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"

	"golang.org/x/exp/slices"
)

type Project struct {
	XMLName			xml.Name			`xml:"Project"`
	Id					string				`xml:"projectId,attr"`
	Name				string				`xml:"name,attr"`
	EntityTypes	[]EntityType 	`xml:"EntityType"`
}

type EntityType struct {
	XMLName		xml.Name 	`xml:"EntityType"`
	Name			string		`xml:"name,attr"`
	Folders		[]Folder	`xml:"Folder"`
	Entities	[]Entity	`xml:"Entity"`
}

type Folder struct {
	XMLName			xml.Name 	`xml:"Folder"`
	Id					string		`xml:"entityId,attr"`
	Name				string		`xml:"name,attr"`
	Subfolders	[]Folder	`xml:"Folder"`
	Entities		[]Entity	`xml:"Entity"`
}

type Entity struct {
	XMLName	xml.Name	`xml:"Entity"`
	Id			string		`xml:"entityId,attr"`
	Name		string		`xml:"name,attr"`
}

// ParseProject reads project.xml file with the project structure.
func ParseProject(envPath string) (*Project, error) {
	projectFile, err := os.Open(fmt.Sprintf("%s\\project.xml", envPath))
	if err != nil {
		return nil, err
	}

	defer projectFile.Close()

	bytes, err := io.ReadAll(projectFile)
	if err != nil {
		return nil, err
	}

	var project Project

	err = xml.Unmarshal(bytes, &project)
	if err != nil {
		return nil, err
	}

	return &project, nil
}

// GetEntityType returns a specified EntityType from a Project.
func GetEntityType(project *Project, name string) *EntityType {
	scriptIdx := slices.IndexFunc(project.EntityTypes, func(et EntityType) bool { return et.Name == name} )
	return &project.EntityTypes[scriptIdx]
}

// CreateDirs performs DFS creation of ID-named directories.
func CreateDirs(et *EntityType, dirs *map[string]string, path string) error {
	// scripts
	scriptPath := fmt.Sprintf("%s\\Script", path)
	if err := os.Mkdir(scriptPath, os.ModePerm); err != nil {
		return err
	}

	for _, folder := range et.Folders {
		// save and create top folders
		folderPath := fmt.Sprintf("%s\\%s", scriptPath, folder.Name)
		(*dirs)[folder.Name] = folderPath
		if err := os.Mkdir(folderPath, os.ModePerm); err != nil {
			return err
		}
		// save and create subdirectories recursively
		err := createSubfolders(&folder, dirs, folderPath)
		if err != nil {
			return err
		}
	}

	return nil
}

// createSubfolders recursively creates a folder's subdirectories and their subdirectories.
func createSubfolders(parent *Folder, dirs *map[string]string, parentPath string) error {
	for _, folder := range (*parent).Subfolders {
		folderPath := fmt.Sprintf("%s\\%s", parentPath, folder.Name)
		(*dirs)[folder.Name] = folderPath
		if err := os.Mkdir(folderPath, os.ModePerm); err != nil {
			return err
		}
		// save and create subdirectories recursively
		if len(folder.Subfolders) > 0 {
			err := createSubfolders(&folder, dirs, folderPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// FindEntityDir searches for entity ID in EntityType and returns its directory.
func FindEntityDir(et *EntityType, dirs *map[string]string, id string, rootDir string) (*Entity, string) {
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
				dir = (*dirs)[folder.Name]
				return entity, dir
			}
		}
		// subdirectories
		entity, dir = findSubfolderEntityDir(&folder, dirs, id)
		if dir == "" {
			continue
		}
		return entity, dir
	}

	return nil, dir
}

// findSubfolderEntityDir recursively searches for entity ID in a folder and its subdirectories.
func findSubfolderEntityDir(parent *Folder, dirs *map[string]string, id string) (*Entity, string) {
	var entity *Entity
	dir := ""
	// subdirs
	for _, folder := range parent.Subfolders {
		for _, ent := range folder.Entities {
			if ent.Id == id {
				entity = &ent
				dir = (*dirs)[folder.Name]
				return entity, dir
			}
		}
		resultEnt, resultDir := findSubfolderEntityDir(&folder, dirs, id)
		if resultDir == "" {
			continue
		}
		return resultEnt, resultDir
	}
	// not found
	return nil, ""
}
