package model

// ParamValidator .
type ParamValidator interface {
	Validate() bool
}

type Model struct {
	LastDated string `json:"last_dated"`
	Dated string `json:"dated"`
}

