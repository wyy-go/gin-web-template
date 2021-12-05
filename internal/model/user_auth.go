package model

// 用户实名认证
type UserAuth struct {
	Id       int64  `json:"id" gorm:"primaryKey;column:id;type:bigint(20)"`
	Uid      int64  `json:"uid" gorm:"column:uid;type:bigint(20);not null"`
	CertType int    `json:"cert_type" gorm:"column:cert_type;type:tinyint(4);not null;default:0"`
	CertNo   string `json:"cert_no" gorm:"column:cert_no;type:varchar(50);not null"`
	Name     string `json:"name" gorm:"column:name;type:varchar(50);not null"`
}

func (_ *UserAuth) TableName() string {
	return "user_auth"
}
