package script

import (
	"fmt"
	"jbextractor/jitterbit/entity"
	"jbextractor/jitterbit/project"
	"os"
	"strings"
)

// CreateScripts creates .jb source code files.
func CreateScripts(scripts *project.EntityType, dirs *map[string]string, envPath string, targetPath string, sep string) error {
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
			
			ent, scriptDir := project.FindEntityDir(scripts, dirs, script.Header.Id, outPath)

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
