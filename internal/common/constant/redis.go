package constant

const (
	RedisKeyToken        = "PASSPORT:TOKEN:%d:%s:%d"         // "PASSPORT:TOKEN:appid:plat:uid"
	RedisKeyRefreshToken = "PASSPORT:REFRESH_TOKEN:%d:%s:%d" // "PASSPORT:REFRESH_TOKEN:appid:plat:uid"
)
