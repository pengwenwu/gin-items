package helper

import (
	"strings"

	mapSet "github.com/deckarep/golang-set"
)

func GetVerifyField(fields []string, getField string) (verifyField string) {
	if len(fields) == 0 || getField == ""{
		return
	}
	getField = strings.Replace(getField, "，", ",", -1)
	getField = strings.Replace(getField, " ", "", -1)
	getFieldSlice := strings.Split(getField, ",")
	fieldsSet := mapSet.NewSet(fields)
	getFieldSet := mapSet.NewSet(getFieldSlice)
	intersectSet := fieldsSet.Intersect(getFieldSet)
	verifyFieldSet := make(map[string]interface{})
	intersectSet.Each(func(i interface{}) bool {
		tmp := strings.Split(, ".")
	})

	return
}
