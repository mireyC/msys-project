package common

type BusinessCode int
type Result struct {
	Code BusinessCode `json:"code"`
	Msg  string       `json:"msg"`
	Data any          `json:"data"`
}

func (r *Result) Success(data any) *Result {
	r.Code = 200
	r.Msg = "success"
	r.Data = data
	return r
}

func (r *Result) Fail(code BusinessCode, msg string) *Result {
	r.Code = code
	r.Msg = msg
	return r
}

/*
model
一般用途
在实际项目中，model.go 文件可能包含以下内容：

数据库实体：映射到数据库表的结构体。
API 响应结构：定义与 API 返回格式一致的结构体，如 Result。
DTO 或 VO：数据传输对象或值对象，用于在不同层之间传输数据。
常量和状态码：与业务相关的常量或状态码，如 BusinessCode。


*/
