package constant

const (
	SignKey    = "a74db8b7-3b97-4653-8e80-ae90ba0e81b3"
	AccountKey = "Account"
)

const (
	Online     = 1
	PushOnline = 2
	Offline    = 3
)

const (
	PushOnlineKeepDays = 30 // 推送在线状态保持天数
)

const (
	RetryCD = 3600 // 登录错误次数超限后再次登录冷却时间
)