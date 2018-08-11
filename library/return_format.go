package library

import (
	"encoding/json"
	"reflect"
)

type (
	returnFormat struct {
		RespNo  int64
		RespMsg string
		Data   interface{}
	}
)
/**
return format
 */
func ReturnJsonWithError(errNo int64, errMsg string, data interface{}) (res string, err error) {

	if data == nil || !reflect.ValueOf(data).IsValid(){
		data = ""
	}
	if errMsg == "ref" {
		errMsg = CodeString(errNo)
	}

	formatter := new(returnFormat)
	formatter.RespNo = errNo
	formatter.RespMsg = errMsg
	formatter.Data = data

	result, err := json.Marshal(formatter)

	if err != nil {
		return "", err
	}

	return string(result), nil
}
