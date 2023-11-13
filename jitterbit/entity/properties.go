package entity

import (
	"encoding/xml"
)

// Universal store for entity parameters and metadata.
type Properties struct {
	XMLName xml.Name `xml:"Properties"`
	Items   []Item   `xml:"Item"`
}
