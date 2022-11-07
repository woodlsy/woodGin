package helper

type Error struct {
	Code    int
	Message string
}

var (
	ErrorSystem = &Error{500, "系统错误，请联系管理员"}
)
