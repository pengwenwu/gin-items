package model

import (
	"time"
)

type Model struct {
	LastDated time.Time `json:"last_dated"`
	Dated     time.Time `json:"dated"`
}
