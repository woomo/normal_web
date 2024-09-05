package model

import (
	"normal_web/util"
	"time"
)

type User struct {
	ID                 int       `gorm:"column:user_id;primaryKey"`
	Email              string    `gorm:"column:email;type:varchar(50)" `
	NickName           string    `gorm:"column:nick_name;type:varchar(30)"`
	Avatar             string    `gorm:"column:avatar;type:varchar(255)"`
	Password           string    `gorm:"column:password;type:varchar(60)"`
	Sex                string    `gorm:"column:sex;type:tinyint(4)"`
	LastUseDeviceID    string    `gorm:"column:last_use_device_id;varchar(32)"`
	LastUseDeviceBrand string    `gorm:"column:last_use_device_brand;varchar(30)"`
	LastLoginIP        string    `gorm:"column:last_login_ip;varchar(128)"`
	Status             int       `gorm:"column:status;tinyint(1)"`
	Phone              string    `gorm:"column:phone;type:varchar(11)"`
	City               string    `gorm:"column:city;type:varchar(255)"`
	Introduction       string    `gorm:"column:introduction;type:varchar(255)"`
	JoinTime           time.Time `gorm:"column:join_time;datetime"`
	LastLoginTime      time.Time `gorm:"column:last_login_time;datetime"`
}

func (User) TableName() string { //显示指定表名
	return "user"
}

var (
	_all_user_field = util.GetGormFields(User{})
)
