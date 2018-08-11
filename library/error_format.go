package library

import "fmt"

const CommenErr int64 = 10001
const CodeErrApi int64 = 40000
const DataError int64 = 10002
const AddPostFail int64 = 20001
const InternalError int64 = 50000
const HttpError int64 = 30001
const CodeSucc int64 = 10000
const ParamFail int64 = 10003


func CodeString(errorNo int64) string {
	switch errorNo {
	case CodeSucc:
		return "success"
	case DataError:
		return "data error"
	case HttpError:
		return "http request fail"
	case ParamFail:
		return "param unexpected"
	case CommenErr:
		return "system error"
	default:
		return "system error"
	}
}

type WhaleErr struct {
	errorCode int64
	errorMsg string
}

func WhaleError(errorCode int64, errorMessage string) *WhaleErr {
	return &WhaleErr{errorCode: errorCode, errorMsg: errorMessage}
}

func (e WhaleErr) Error() string{
	return fmt.Sprintf("%s", e.errorMsg)
}
func (e WhaleErr) Code() string{
	return fmt.Sprintf("%d", e.errorCode)
}



