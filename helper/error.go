package helper

type Error struct {
	Code    int `json:"code"`
	Message string `json:"message"`
}

var (
	ErrorSystem = &Error{500, "系统错误，请联系管理员"}
)
