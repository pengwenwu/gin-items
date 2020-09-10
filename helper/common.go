package helper

import (
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/validation"
	mapSet "github.com/deckarep/golang-set"

	"gin-items/library/ecode"
)

// 校验字段
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

// 设置校验参数错误
func GetEcodeValidParam(validErrors []*validation.Error) (err error) {
	if validErrors == nil {
		return
	}
	for _, validErr := range validErrors {
		ecode.Register(map[int]string{ecode.IllegalParams.Code(): validErr.Key + " " + validErr.Message})
		err = ecode.IllegalParams
		return
	}
	return
}

// 绑定json错误
func GetEcodeBindJson(bindErr error) (err error) {
	if bindErr == nil {
		return
	}
	ecode.Register(map[int]string{ecode.RequestErr.Code(): bindErr.Error()})
	err = ecode.RequestErr
	return
}

func FormatDateTimeNow() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func FormatDateTimeZero() string {
	return "0000-00-00 00:00:00"
}
