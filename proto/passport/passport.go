package passport

type DeviceInfo struct {
	DevicePlatform string `form:"device_platform" json:"device_platform" binding:"required"` // 设备平台
	DeviceId       string `form:"device_id" json:"device_id" binding:"required"`             // 设备ID
	DeviceName     string `form:"device_name" json:"device_name" binding:"required"`         // 设备名称
	DeviceModel    string `form:"device_model" json:"device_model" binding:"required"`       // 机型
	OsVersion      string `form:"os_version" json:"os_version"`                              // 操作系统版本
	Screen         string `form:"screen" json:"screen"`                                      // 屏幕
	Channel        string `form:"channel" json:"channel"`                                    // 渠道
}

type TokenInfo struct {
	Uid          int64  `json:"uid"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
}

type LoginRequest struct {
	DeviceInfo
	Type    int    `form:"type" json:"type" binding:"required"`
	Account string `form:"account" json:"account" binding:"required"`
	Passwd  string `form:"passwd" json:"passwd" binding:"required"`
}

type SmsRequest struct {
	Mobile string `form:"mobile" json:"mobile" binding:"required"`
}

type SmsResponse struct {
	Code string `json:"code"`
}

type SmsLoginRequest struct {
	DeviceInfo
	Mobile string `form:"mobile" json:"mobile" binding:"required"`
	Code   string `form:"code" json:"code" binding:"required"`
}

type SmsLoginResponse struct {
	*TokenInfo
	Register int `json:"register"`
}

type SetPwdRequest struct {
	Passwd string `form:"passwd" json:"passwd" binding:"required"`
}

type ChangePwdRequest struct {
	OldPasswd string `form:"old_passwd" json:"old_passwd" binding:"required"`
	NewPasswd string `form:"new_passwd" json:"new_passwd" binding:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `form:"refresh_token" json:"refresh_token" binding:"required"`
}

type GetWxInfoRequest struct {
	App  string `form:"app" json:"app"`
	Code string `form:"code" json:"code"`
}

type GetWxInfoResponse struct {
	Unionid    string `json:"unionid"`
	Openid     string `json:"openid"`
	Nickname   string `json:"nickname"`
	Headimgurl string `json:"headimgurl"`
	Uid        int64  `json:"uid"`
	Mobile     string `json:"mobile"`
}

type OAuthLoginRequest struct {
	DeviceInfo
	Platform string `form:"platform" json:"platform" binding:"required"`
	Openid   string `form:"openid" json:"openid" binding:"required"`
	Mobile   string `form:"mobile" json:"mobile"`
	Code     string `form:"code" json:"code"`
}

type OAuthLoginResponse struct {
	*TokenInfo
	Register int `json:"register"`
}

type PreChangeMobileRequest struct {
	Mobile string `form:"mobile" json:"mobile"`
	Code   string `form:"code" json:"code"`
}

type ChangeMobileRequest struct {
	Mobile string `form:"mobile" json:"mobile"`
	Code   string `form:"code" json:"code"`
}

type PreBindEmailRequest struct {
	Email string `form:"email" json:"email"`
}

type BindEmailRequest struct {
	Email string `form:"email" json:"email"`
	Code  string `form:"code" json:"code"`
}
