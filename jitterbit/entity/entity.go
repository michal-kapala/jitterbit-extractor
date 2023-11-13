package entity

import (
	"encoding/xml"
	"io"
	"os"
)

// Universal Jitterbit object, e.g. a script, operation, source or target.
type Entity struct {
	XMLName     xml.Name   `xml:"Entity"`
	Type        string     `xml:"type,attr"`
	Header      Header     `xml:"Header"`
	Props       Properties `xml:"Properties"`
	KongaString string     `xml:"konga.string"`
	Pipeline    Pipeline   `xml:"Pipeline"`
}

// Parses an XML file from <project>/<environment>/Data/<any>.
func ParseEntity(filePath string) (*Entity, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var entity Entity

	err = xml.Unmarshal(bytes, &entity)
	if err != nil {
		return nil, err
	}

	return &entity, nil
}
