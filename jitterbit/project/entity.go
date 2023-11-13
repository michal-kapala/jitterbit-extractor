package project

import (
	"encoding/xml"
)

type Entity struct {
	XMLName	xml.Name	`xml:"Entity"`
	Id			string		`xml:"entityId,attr"`
	Name		string		`xml:"name,attr"`
}
