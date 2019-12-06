package app

import (
	"fmt"
	"github.com/astaxie/beego/validation"
)

func MakeErrors(errors []*validation.Error) {
	for _, err := range errors {
		fmt.Println(err.Key, err.Message)
	}
	return
}
