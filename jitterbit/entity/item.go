package entity

import (
	"encoding/xml"
)

// Universal KV data item. Enc indicates encrypted values.
type Item struct {
	XMLName xml.Name `xml:"Item"`
	Key     string   `xml:"key,attr"`
	Value   string   `xml:"value,attr"`
	Enc     bool     `xml:"enc,attr"`
}
