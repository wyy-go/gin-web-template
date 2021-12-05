package errors

const (
	ErrCustom int32 = 999999
)

// TODO：完善错误处理
var (
	ErrSign                = NewError(100001, "签名错误")
	ErrInvalidParam        = NewError(100002, "参数错误")
	ErrInvalidToken        = NewError(100003, "无效的令牌")
	ErrTokenExpired        = NewError(100004, "令牌过期")
	ErrTokenRevoked        = NewError(100005, "令牌已失效")
	ErrUnAuthorized        = NewError(100006, "未授权")
	ErrNoLogin             = NewError(100007, "未登录")
	ErrPassword            = NewError(100008, "用户名或密码错误")
	ErrUserNotExists       = NewError(100009, "用户不存在")
	ErrUserAlreadyExists   = NewError(100010, "用户已存在")
	ErrAccountNotAvailable = NewError(100011, "账户不可用")
	ErrSendCodeTooFrequent = NewError(100012, "验证码发送过于频繁")
	ErrSendCodeLimit       = NewError(100012, "当日发送验证码次数过多")
	ErrRequireSendCode     = NewError(100013, "请先发送短信验证码")
	ErrVerifyCodeExpired   = NewError(100014, "验证码过期")
	ErrVerifyCodeLimit     = NewError(100015, "验证码失效")
	ErrVerifyCode          = NewError(100016, "验证码错误")
	ErrInvalidUserTag      = NewError(100017, "技能已被删除")
	ErrErrorTimesLimit     = NewError(100018, "密码错误次数过多")
	ErrPayPasswd           = NewError(100019, "支付密码错误")
	ErrPayPasswdLock       = NewError(100020, "支付密码已被锁定")
)
