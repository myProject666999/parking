package utils

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Count   int64       `json:"count,omitempty"`
}

func Success(data interface{}) Response {
	return Response{
		Code:    0,
		Message: "操作成功",
		Data:    data,
	}
}

func SuccessWithCount(data interface{}, count int64) Response {
	return Response{
		Code:    0,
		Message: "操作成功",
		Data:    data,
		Count:   count,
	}
}

func SuccessMsg(message string) Response {
	return Response{
		Code:    0,
		Message: message,
	}
}

func Error(message string) Response {
	return Response{
		Code:    1,
		Message: message,
	}
}

func ErrorCode(code int, message string) Response {
	return Response{
		Code:    code,
		Message: message,
	}
}
