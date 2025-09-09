package errs

type Code int

// 通用
const (
	Success             Code = 200
	InternalServerError Code = 500 // 服务器错误
	Unknown             Code = 1
	ParamInvalid        Code = 400
	MissingParameter    Code = 401 //  缺少参数
	Internal            Code = 3
)

// 游戏模块
const (
	GameNotFound        Code = 2001
	SpinFailed          Code = 2002
	InsufficientBalance Code = 2003 //余额不足
	BetOverLimit        Code = 2004
	CmdNotFound         Code = 2005
	Timeout             Code = 2006
	DataFormatError     Code = 2007 // 数据格式有问题
	ConnClosed          Code = 2008 // 连接已关闭
	WrongBetAmount      Code = 2009 //下注金额错误
)

// 系统/服务模块
const (
	DBError         Code = 9001
	RedisError      Code = 9002
	ThirdPartyError Code = 9003
)
