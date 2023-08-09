package models

import (
	"fmt"
	"time"

	"github.com/toma-gin-chat/utils"
	"gorm.io/gorm"
)

type UserBasic struct {
	gorm.Model
	// NickName      string `json:"nickName" gorm:"column:nick_ame;type:varchar(128);comment:'昵称'"`
	// UserName      string `json:"userName" gorm:"column:user_name;type:varchar(128);comment:'用户名'"`
	Name          string `json:"name" gorm:"column:name;type:varchar(128);comment:'用户名'"`
	Password      string `json:"password" gorm:"column:password;type:varchar(128);comment:'密码'"`
	Phone         string `json:"phone" gorm:"column:phone;type:varchar(128);comment:'手机号'" valid:"matches(^1[3-9]{1}\\d{9})"`
	Email         string `json:"email" gorm:"column:email;type:varchar(128);comment:'邮箱'" valid:"email"`
	Identity      string `json:"identity" gorm:"column:identity;type:varchar(128);comment:'认证'"`
	ClientIP      string `json:"clientIP" gorm:"column:client_ip;type:varchar(128);comment:'客户端ip'"`
	ClientPort    string `json:"clientPort" gorm:"column:client_port;type:varchar(128);comment:'客户端端口'"`
	IsLogout      bool   `json:"isLogout" gorm:"column:is_logout;type:tinyint;comment:'是否登出'"`
	DeviceInfo    string `json:"deviceInfo" gorm:"column:device_info;type:varchar(128);comment:'设备信息'"`
	Salt          string `json:"-" gorm:"column:salt;type:varchar(128);comment:'加密密盐'"`
	LoginTime     *time.Time
	HeartbeatTime *time.Time
	LoginOutTime  *time.Time
}

func (t *UserBasic) TableName() string {
	return "user_basic"
}

// 查询用户列表
func GetUserList() []UserBasic {
	results := []UserBasic{}
	utils.Db.Find(&results)
	return results
}

// 创建用户
func CreateUser(user UserBasic) *gorm.DB {
	return utils.Db.Create(&user)
}

// 更新用户
func UpdateUser(user UserBasic) *gorm.DB {
	return utils.Db.Model(&user).Updates(UserBasic{
		Name:     user.Name,
		Password: user.Password,
		Phone:    user.Phone,
		Email:    user.Email,
	})
}

// 删除用户
func DeleteUser(id int) *gorm.DB {
	return utils.Db.Delete(new(UserBasic), "id = ?", id)
}

func FindUserByName(value string, id int) *UserBasic {
	var result UserBasic
	db := utils.Db
	var err error
	if id == 0 {
		err = db.First(&result, "name = ?", value).Error
	} else {
		err = db.First(&result, "id <> ? AND name = ?", id, value).Error
	}
	if err == nil {
		return &result
	}
	fmt.Printf("通过用户名(%s)查询错误, %s\n", value, err.Error())
	return nil
}

func FindUserByID(id int) *UserBasic {
	var result UserBasic
	err := utils.Db.First(&result, "id = ?", id).Error
	if err == nil {
		return &result
	}

	fmt.Printf("通过ID(%d)查询错误, %s\n", id, err.Error())
	return nil
}

func FindUserByEmail(value string, id int) *UserBasic {
	var result UserBasic
	db := utils.Db
	var err error
	if id == 0 {
		err = db.First(&result, "email = ?", value).Error
	} else {
		err = db.First(&result, "id <> ? AND email = ?", id, value).Error
	}
	if err == nil {
		return &result
	}
	fmt.Printf("通过邮箱(%s)查询错误, %s\n", value, err.Error())
	return nil
}

func FindUserByPhone(value string, id int) *UserBasic {
	var result UserBasic
	db := utils.Db
	var err error
	if id == 0 {
		err = db.First(&result, "phone = ?", value).Error
	} else {
		err = db.First(&result, "id <> ? AND phone = ?", id, value).Error
	}
	if err == nil {
		return &result
	}
	fmt.Printf("通过手机号(%s)查询错误, %s\n", value, err.Error())
	return nil
}

func FindUserByNameAndPWD(name, pwd string) *UserBasic {
	var result UserBasic
	err := utils.Db.First(&result, "(name = ? OR phone =? OR email = ?) AND password = ?", name, name, name, pwd).Error
	if err == nil {
		return &result
	}

	fmt.Printf("通过查询错误, %s\n", err.Error())
	return nil
}
