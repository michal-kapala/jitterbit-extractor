package entity

import (
	"encoding/xml"
)

type Properties struct {
	XMLName	xml.Name	`xml:"Properties"`
	Items		[]Item		`xml:"Item"`
}
