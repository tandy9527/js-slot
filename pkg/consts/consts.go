package consts

import "time"

// 响应相关常量
const (
	RESP_ERROR = "error"
	RESP_DATA  = "data"
	LOGIN_CMD  = "login" //  登陆命令
)

// 业务相关常量
const (
	// JWT key
	JWT_SECRET_KEY = "jti"
	// JWT 过期时间 2小时，单位秒
	JWT_EXPIRE_TIME = 2 * 60 * 60 * time.Second

	// 验证码过期时间  5分钟，单位秒
	CODE_EXPIRE_TIME = 5 * 60 * time.Second
	// 验证码长度
	CODE_LENGTH = 6
)

const (
	// 游戏状态-正常
	GAME_STATUS_ONLINE = 0
	//  不可见
	GAME_STATUS_OFFLINE = 1
)

// HTTP 头常量
const (
	AUTHORIZATION_HEADER = "Authorization"
	CONTENT_TYPE_JSON    = "application/json"
	// TOKEN 类型
	TOKEN_TYPE = "bearer"
)

// 用户状态
const (
	//  被锁定
	USER_STATUS_LOCK = 1
	//  正常
	USER_STATUS_INACTIVE = 0
)

//------------------------- Redis 相关----------------------------------------
const (
	//  手机验证码ke'y
	REDIS_VERIFY_CODE_PHONE = "verify_code:phone:"
	//  邮箱验证码key
	REDIS_VERIFY_CODE_EMAIL = "verify_code:email:"
	// 用户ID自增key
	REDIS_USER_ID = "user:id"
	// 用户信息缓存前缀
	REDIS_USER_KEY = "user:"
	// 用户状态字段
	REDIS_USER_STATE = "state"
	// 用户Token缓存前缀
	REDIS_TOKEN_KEY = "token:"
	// 邀请码集合key
	REDIS_INVITATION_CODES = "invitation_codes"
	// 连接数
	REDIS_CONN_TOTAL = "conn:total"
	// 游戏连接数
	REDIS_GAME_TOTAL = "conn:game:"
	//  待处理数据队列
	REDIS_DATA_QUEUE_PENDING = "data_queue:pending"
	//  正在处理的数据队列
	REDIS_DATA_QUEUE_PROCESSING = "data_queue:processing"
	//  处理失败的数据队列
	REDIS_DATA_QUEUE_DEAD = "data_queue:dead"
)
