package response

type Response struct {
	Success bool   `json:"success"` // true or false
	Msg     string `json:"msg"`     // error message
	Obj     interface{}
}

func Error(msg string) Response {
	return Response{
		Success: false,
		Msg:     msg,
	}
}

func OK(obj interface{}) Response {
	return Response{
		Success: true,
		Obj:     obj,
	}
}
