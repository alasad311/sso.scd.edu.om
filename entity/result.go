package entity

type Result struct {
	Code    int         `json:"Error"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data"`
}

func (res *Result) SetCode(code int) *Result {
	res.Code = code
	return res
}

func (res *Result) SetMessage(msg string) *Result {
	res.Message = msg
	return res
}

func (res *Result) SetData(data interface{}) *Result {
	res.Data = data
	return res
}
