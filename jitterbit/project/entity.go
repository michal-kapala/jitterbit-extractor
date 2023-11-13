package project

import (
	"encoding/xml"
)

// A data entity location.
type Entity struct {
	XMLName xml.Name `xml:"Entity"`
	Id      string   `xml:"entityId,attr"`
	Name    string   `xml:"name,attr"`
}
