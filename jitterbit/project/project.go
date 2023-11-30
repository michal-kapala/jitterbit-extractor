package project

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/exp/slices"
)

// The project structure saved in project.xml file.
type Project struct {
	XMLName     xml.Name     `xml:"Project"`
	Id          string       `xml:"projectId,attr"`
	Name        string       `xml:"name,attr"`
	EntityTypes []EntityType `xml:"EntityType"`
	// Environment path.
	EnvPath string
}

// ParseProject reads project.xml file with the project structure.
func ParseProject(envPath string, sep string) (*Project, error) {
	projectFile, err := os.Open(fmt.Sprintf("%s%sproject.xml", envPath, sep))
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

	for idx := range project.EntityTypes {
		project.EntityTypes[idx].Dirs = make(map[string]string)
	}

	project.EnvPath = envPath
	return &project, nil
}

// GetEntityType returns a specified EntityType.
func (project *Project) GetEntityType(name string) *EntityType {
	scriptIdx := slices.IndexFunc(project.EntityTypes, func(et EntityType) bool { return et.Name == name })
	et := &project.EntityTypes[scriptIdx]
	et.Type = name
	return et
}

// updateDirPaths updates all paths with a restored directory name.
func updateDirPaths(dirs *map[string]string, id string, name string) {
	for dirId, path := range *dirs {
		(*dirs)[dirId] = strings.Replace(path, id, name, 1)
	}
}
