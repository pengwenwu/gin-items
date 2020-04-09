package helper

import (
	"fmt"
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
	fieldsSet := mapSet.NewSet()
	for _,v := range fields {
		fieldsSet.Add(v)
	}
	getFieldSet := mapSet.NewSet()
	for _,v := range getFieldSlice {
		getFieldSet.Add(v)
	}
	intersectSet := fieldsSet.Intersect(getFieldSet)

	var verifyFieldSet []string
	intersectSet.Each(func(i interface{}) bool {
		iToStr := fmt.Sprintf("%v", i)
		if tmp := strings.Split(iToStr, ".");len(tmp) == 2 {
			verifyFieldSet = append(verifyFieldSet, tmp[1])
		} else {
			verifyFieldSet = append(verifyFieldSet, iToStr)
		}
		return false
	})

	verifyField = strings.Join(verifyFieldSet, ",")
	return
}
