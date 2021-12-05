package model

import (
	"gorm.io/plugin/soft_delete"
	"time"
)

// 用户设备
type UserDevice struct {
	Id             int64                 `json:"id" gorm:"primaryKey;column:id;type:bigint(20)"`            // 系统编号
	Appid          int                   `json:"appid" gorm:"column:appid;type:int(11);not null;default:0"` // 应用ID
	Ua             string                `json:"ua" gorm:"column:ua;type:varchar(50);not null"`
	Uid            int64                 `json:"uid" gorm:"column:uid;type:bigint(20);not null;default:0"`                // 用户ID
	DevicePlatform string                `json:"device_platform" gorm:"column:device_platform;type:varchar(50);not null"` // 设备平台类型
	DeviceId       string                `json:"device_id" gorm:"column:device_id;type:varchar(64);not null"`             // 设备ID
	DeviceName     string                `json:"device_name" gorm:"column:device_name;type:varchar(50);not null"`         // 设备名称
	Vendor         string                `json:"vendor" gorm:"column:vendor;type:varchar(50);not null"`                   // 设备厂商
	PushClientid   string                `json:"push_clientid" gorm:"column:push_clientid;type:varchar(50);not null"`     // 推送设备客户端标识
	Oaid           string                `json:"oaid" gorm:"column:oaid;type:varchar(50);not null"`                       // 移动智能设备标识公共服务平台提供的匿名设备标识符(OAID)
	Idfa           string                `json:"idfa" gorm:"column:idfa;type:varchar(50);not null"`                       // iOS平台配置应用使用广告标识(IDFA)
	Imei           string                `json:"imei" gorm:"column:imei;type:varchar(50);not null"`                       // 国际移动设备识别码IMEI(International Mobile Equipment Identity)
	Model          string                `json:"model" gorm:"column:model;type:varchar(50);not null"`                     // 设备型号
	CreatedAt      time.Time             `json:"created_at" gorm:"column:created_at;type:datetime;not null"`              // 创建时间
	UpdatedAt      time.Time             `json:"updated_at" gorm:"column:updated_at;type:datetime;not null"`              // 更新时间
	DeletedAt      soft_delete.DeletedAt `json:"deleted_at" gorm:"column:deleted_at;type:bigint(20);not null;default:0"`  // 删除时间
	LastLoginTime  time.Time             `json:"last_login_time" gorm:"column:last_login_time;type:datetime"`             // 最近登录时间
	LastLoginIp    string                `json:"last_login_ip" gorm:"column:last_login_ip;type:varchar(50);not null"`     // 最近登录IP
}

func (_ *UserDevice) TableName() string {
	return "user_device"
}
