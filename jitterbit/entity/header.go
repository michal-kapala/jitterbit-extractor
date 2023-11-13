package entity

import (
	"encoding/xml"
)

// Universal entity identification and deployment data.
type Header struct {
	XMLName     xml.Name `xml:"Header"`
	Deleted     bool     `xml:"Deleted,attr"`
	DeployDirty bool     `xml:"DeployDirty,attr"`
	Deployed    bool     `xml:"Deployed,attr"`
	HasMoved    bool     `xml:"HasMoved,attr"`
	Id          string   `xml:"ID,attr"`
	Name        string   `xml:"Name,attr"`
	KongaString string   `xml:"konga.string"`
}
