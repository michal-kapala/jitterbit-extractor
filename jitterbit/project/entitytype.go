package project

import (
	"encoding/xml"
	"fmt"
	"os"
	"strings"

	"jbextractor/jitterbit/entity"
)

// Entity type identifier.
const (
	API_ENTITY           string = "ApiEntity"
	DOCUMENT             string = "Document"
	EMAIL                string = "EmailMessage"
	JMS_ACK              string = "JMSAcknowledge"
	JMS_BROWSE           string = "JMSBrowse"
	JMS_ENDPOINT         string = "JMSEndpoint"
	JMS_LISTEN           string = "JMSListen"
	JMS_POLL             string = "JMSPoll"
	JMS_SEND_PUBLISH     string = "JMSSendOrPublish"
	MSCRM_ENDPOINT       string = "MSCrmEndpoint"
	MSCRM_QUERY          string = "MSCrmQuery"
	MSCRM_UPSERT         string = "MSCrmUpsert"
	NETSUITE_ENDPOINT    string = "NetSuiteEndpoint"
	NETSUITE_GET_LIST    string = "NetSuiteGetList"
	NETSUITE_QUERY       string = "NetSuiteQuery"
	NETSUITE_UPSERT      string = "NetSuiteUpsert"
	OPERATION            string = "Operation"
	QUICKBOOKS_CREATE    string = "QuickBooksCreate"
	QUICKBOOKS_ENDPOINT  string = "QuickBooksEndpoint"
	QUICKBOOKS_QUERY     string = "QuickBooksQuery"
	SALESFORCE_CONNECTOR string = "SalesforceConnector"
	SALESFORCE_CREATE    string = "SalesforceCreate"
	SALESFORCE_QUERY     string = "SalesforceQuery"
	SALESFORCE_UPSERT    string = "SalesforceUpsert"
	SAP_ENDPOINT         string = "SapEndpoint"
	SAP_FUNCTION         string = "SapFunction"
	SCHEDULE             string = "Schedule"
	SCRIPT               string = "Script"
	SOURCE               string = "Source"
	TARGET               string = "Target"
	TRANSFORMATION       string = "Transformation"
	VARIABLE             string = "ProjectVariable"
	WS_CALL              string = "WebServiceCall"
	XSLT_TRANSFORM       string = "XsltTransform"
)

// A category of Jitterbit objects, e.g. scripts or operations.
type EntityType struct {
	XMLName  xml.Name `xml:"EntityType"`
	Name     string   `xml:"name,attr"`
	Folders  []Folder `xml:"Folder"`
	Entities []Entity `xml:"Entity"`
	// Root directory name.
	Type string
	// Directory paths.
	Dirs map[string]string
}

// CreateDirs performs DFS creation of ID-named directories.
func (et *EntityType) CreateDirs(path string, sep string) error {
	// scripts
	scriptPath := fmt.Sprintf("%s%s%s", path, sep, et.Type)
	if err := os.Mkdir(scriptPath, os.ModePerm); err != nil {
		return err
	}

	for _, folder := range et.Folders {
		// save and create top folders
		folderPath := fmt.Sprintf("%s%s%s", scriptPath, sep, folder.Id)
		et.Dirs[folder.Id] = folderPath
		if err := os.Mkdir(folderPath, os.ModePerm); err != nil {
			return err
		}
		// save and create subdirectories recursively
		err := folder.createSubfolders(&et.Dirs, folderPath, sep)
		if err != nil {
			return err
		}
	}

	return nil
}

// FindEntity searches for entity ID in EntityType and returns its name with parent directory.
func (et *EntityType) FindEntity(id string, rootDir string) (*Entity, string) {
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
				dir = et.Dirs[folder.Id]
				return entity, dir
			}
		}
		// subdirectories
		entity, dir = folder.findSubfolderEntity(&et.Dirs, id)
		if dir == "" {
			continue
		}
		return entity, dir
	}

	return nil, dir
}

// RenameDirs substitutes folder IDs with real names.
func (et *EntityType) RenameDirs(path string) error {
	oldPath := ""
	newPath := ""
	for _, folder := range et.Folders {
		oldPath = et.Dirs[folder.Id]
		newPath = strings.Replace(oldPath, folder.Id, sanitizeFileName(folder.Name), 1)
		err := os.Rename(oldPath, newPath)
		if err != nil {
			return err
		}
		// update map
		updateDirPaths(&et.Dirs, folder.Id, folder.Name)
		err = folder.renameDirs(&et.Dirs, newPath)
		if err != nil {
			return err
		}
	}
	return nil
}

// CreateScripts creates .jb source code files.
func (scripts *EntityType) CreateScripts(envPath string, targetPath string, sep string) error {
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
		if !entry.IsDir() && strings.Contains(entryName, ".xml") {
			inFilePath = fmt.Sprintf("%s%s%s", inPath, sep, entryName)
			script, err := entity.ParseEntity(inFilePath)
			if err != nil {
				return err
			}

			ent, scriptDir := scripts.FindEntity(script.Header.Id, outPath)

			// script was not found in project.xml
			if ent == nil || scriptDir == "" {
				return fmt.Errorf("[CreateScripts] Corrupted project.xml - script %s was not found", script.Header.Id)
			}
			// example script name from Jitterbit's demo project:
			// jb.sqlServer.table1-&gt;table2 [ETL_log]
			saneName := sanitizeFileName(script.Header.Name)
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

// Copies operation definitions in XML format.
func (ops *EntityType) CreateOperations(envPath string, targetPath string, sep string) error {
	inPath := fmt.Sprintf("%s%sData%sOperation", envPath, sep, sep)
	outPath := fmt.Sprintf("%s%sOperation", targetPath, sep)
	entries, err := os.ReadDir(inPath)
	if err != nil {
		return err
	}

	inFilePath := ""
	entryName := ""
	for _, entry := range entries {
		entryName = entry.Name()
		if !entry.IsDir() && strings.Contains(entryName, ".xml") {
			inFilePath = fmt.Sprintf("%s%s%s", inPath, sep, entryName)
			op, err := entity.ParseEntity(inFilePath)
			if err != nil {
				return err
			}

			ent, opDir := ops.FindEntity(op.Header.Id, outPath)
			// operation was not found in project.xml
			if ent == nil || opDir == "" {
				return fmt.Errorf("[CreateOperations] Corrupted project.xml - operation %s was not found", op.Header.Id)
			}

			// copy xml
			saneName := sanitizeFileName(op.Header.Name)
			outFilePath := fmt.Sprintf("%s%s%s.xml", opDir, sep, saneName)
			data, err := os.ReadFile(inFilePath)
			if err != nil {
				return err
			}

			err = os.WriteFile(outFilePath, data, os.ModePerm)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
