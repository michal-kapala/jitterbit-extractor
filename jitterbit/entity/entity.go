package entity

import (
	"encoding/xml"
	"io"
	"os"
)

type Entity struct {
	XMLName	xml.Name				`xml:"Entity"`
	Type 				string			`xml:"type,attr"`
	Header			Header			`xml:"Header"`
	Props				Properties	`xml:"Properties"`
	KongaString	string			`xml:"konga.string"`
}

type Header struct {
	XMLName			xml.Name	`xml:"Header"`
	Deleted			bool			`xml:"Deleted,attr"`
	DeployDirty	bool			`xml:"DeployDirty,attr"`
	Deployed		bool			`xml:"Deployed,attr"`
	HasMoved		bool			`xml:"HasMoved,attr"`
	Id					string		`xml:"ID,attr"`
	Name				string		`xml:"Name,attr"`
	KongaString	string		`xml:"konga.string"`
}

type Properties struct {
	XMLName	xml.Name	`xml:"Properties"`
	Items		[]Item		`xml:"Item"`
}

type Item struct {
	XMLName	xml.Name	`xml:"Item"`
	Key			string		`xml:"key,attr"`
	Value		string		`xml:"value,attr"`
	Enc			bool			`xml:"enc,attr"`
}

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
