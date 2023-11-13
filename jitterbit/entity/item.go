package entity

import (
	"encoding/xml"
)

type Item struct {
	XMLName	xml.Name	`xml:"Item"`
	Key			string		`xml:"key,attr"`
	Value		string		`xml:"value,attr"`
	Enc			bool			`xml:"enc,attr"`
}
