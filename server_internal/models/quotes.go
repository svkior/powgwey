package models

// nolint:revive // ok
import _ "github.com/mailru/easyjson/gen"

//go:generate easyjson quotesm.go

// easyjson:json
type Quotes []Quote

type Quote struct {
	Type     string `json:"type"`
	Language string `json:"language"`
	Source   string `json:"source"`
	Quote    string `json:"quote"`
	Author   string `json:"author"`
}
