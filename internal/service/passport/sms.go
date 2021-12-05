package passport

type SmsProvider interface {
	GetName() string
	SendCode(mobile string)
	VerifyCode(mobile, code string) (bool, error)
}
