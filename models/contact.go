package models

import "gorm.io/gorm"

// 人员关系
type Contact struct {
	gorm.Model
	OwnrID   uint // 谁的关系消息
	TargetID uint //  对应的谁
	Type     int  // 对应类型 0 1 2
	Desc     string
}

func (t *Contact) TableName() string {
	return "contact"
}
