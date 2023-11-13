package entity

import (
	"encoding/xml"
)

// Operation-specific tag, stores activity data (operation steps).
type Pipeline struct {
	XMLName    xml.Name   `xml:"Pipeline"`
	OpType     string     `xml:"opType,attr"`
	Activities Activities `xml:"Activities"`
}

// All steps of an operation.
type Activities struct {
	XMLName    xml.Name   `xml:"Activities"`
	Activities []Activity `xml:",any"`
}

// A single step of an operation.
type Activity struct {
	Id        string `xml:"activityId,attr"`
	ContentId string `xml:"contentId,attr"`
	// Identifies the type of activity (string).
	Role string `xml:"role,attr"`
	// Identifies the type of activity (integer).
	Type string `xml:"type,attr"`
}
