package errs

type APIError struct {
	Code Code   `json:"code"`
	Msg  string `json:"msg"`
}

func (e *APIError) Error() string {
	return e.Msg
}

// 常用错误实例
var (
	//成功
	ErrSuccess = &APIError{Code: Success, Msg: "ok"}
	// 没有找到游戏
	ErrGameNotFound = &APIError{Code: GameNotFound, Msg: "game not found"}
	// 未知命令
	ErrCmdNotFound = &APIError{Code: CmdNotFound, Msg: "unknown command"}
	// 超时
	ErrTimeout             = &APIError{Code: Timeout, Msg: "request timeout"}
	ErrInternalServerError = &APIError{Code: InternalServerError, Msg: "internal server error,Please contact the administrator"} // 服务器内部错误
	ErrDataFormatError     = &APIError{Code: DataFormatError, Msg: "data format error"}                                          // 数据格式有问题
	ErrConnClosed          = &APIError{Code: ConnClosed, Msg: "connection closed"}                                               // 连接已关闭
)
