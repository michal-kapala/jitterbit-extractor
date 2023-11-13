package project

import (
	"encoding/xml"
	"fmt"
	"os"
	"strings"
)

// A virtual Jitterbit directory.
type Folder struct {
	XMLName    xml.Name `xml:"Folder"`
	Id         string   `xml:"entityId,attr"`
	Name       string   `xml:"name,attr"`
	Subfolders []Folder `xml:"Folder"`
	Entities   []Entity `xml:"Entity"`
}

// createSubfolders recursively creates a folder's subdirectories and their subdirectories.
func (parent *Folder) createSubfolders(dirs *map[string]string, parentPath string, sep string) error {
	for _, folder := range (*parent).Subfolders {
		folderPath := fmt.Sprintf("%s%s%s", parentPath, sep, folder.Id)
		(*dirs)[folder.Id] = folderPath
		if err := os.Mkdir(folderPath, os.ModePerm); err != nil {
			return err
		}
		// save and create subdirectories recursively
		if len(folder.Subfolders) > 0 {
			err := folder.createSubfolders(dirs, folderPath, sep)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// findSubfolderEntity recursively searches for entity ID in the folder and its subdirectories.
func (parent *Folder) findSubfolderEntity(dirs *map[string]string, id string) (*Entity, string) {
	var entity *Entity
	dir := ""
	// subdirs
	for _, folder := range parent.Subfolders {
		for _, ent := range folder.Entities {
			if ent.Id == id {
				entity = &ent
				dir = (*dirs)[folder.Id]
				return entity, dir
			}
		}
		resultEnt, resultDir := folder.findSubfolderEntity(dirs, id)
		if resultDir == "" {
			continue
		}
		return resultEnt, resultDir
	}
	// not found
	return nil, ""
}

// renameDirs substitutes subfolder IDs with real names.
func (parent *Folder) renameDirs(dirs *map[string]string, path string) error {
	oldPath := ""
	newPath := ""
	for _, folder := range parent.Subfolders {
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
