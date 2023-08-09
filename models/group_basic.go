package models

import "gorm.io/gorm"

// 群消息
type GroupBasic struct {
	gorm.Model
	Name    string // 群名
	OwnerID uint   // 群拥有者
	Icon    string // 头像
	Type    int    // 群类型
	Desc    string
}

func (t *GroupBasic) TableName() string {
	return "group_basic"
}
