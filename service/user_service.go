package service

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/toma-gin-chat/models"
	"github.com/toma-gin-chat/utils"
)

// UserLogin
// @Tags 用户模块
// @Summary 登录
// @param name formData string false "name"
// @param password formData string false "password"
// @Success 200 {string} json{"code","msg","data"}
// @Router /user/login [post]
func UserLogin(c *gin.Context) {
	name := c.PostForm("name")
	pwd := c.PostForm("password")

	user := models.FindUserByName(name, 0)
	if user == nil {
		FailWithMsg("账户或密码错误", c)
		return
	}

	pwdHash := utils.MakePassword(pwd, user.Salt)
	findUser := models.FindUserByNameAndPWD(name, pwdHash)
	if findUser == nil {
		FailWithMsg("账户或密码错误", c)
		return
	}
	j := utils.NewJWT()
	claims := j.CreateClaims(utils.BaseClaims{
		ID:       user.ID,
		Username: user.Name,
	})

	token, err := j.CreateToken(claims)
	if err != nil {
		FailWithMsg("生成token错误", c)
		return
	}

	OkWithDetail(
		map[string]interface{}{"user": user, "token": token},
		"登录成功",
		c,
	)
}

// GetUserList
// @Tags 用户模块
// @Summary 查询所有用户
// @Success 200 {string} json{"code","data"}
// @Router /user/getUserList [get]
func GetUserList(c *gin.Context) {
	list := models.GetUserList()

	OkWithDetail(list, "查询成功", c)
}

// CreateUser
// @Tags 用户模块
// @Summary 新增用户
// @param name formData string false "name"
// @param password formData string false "password"
// @param repassword formData string false "repassword"
// @Success 200 {string} json{"code","msg"}
// @Router /user/createUser [post]
func CreateUser(c *gin.Context) {
	user := models.UserBasic{}
	user.Name = c.PostForm("name")
	pwd := c.PostForm("password")
	repassword := c.PostForm("repassword")
	if pwd != repassword {
		FailWithMsg("两次密码不一致", c)
		return
	}

	if findUser := models.FindUserByName(user.Name, 0); findUser != nil {
		FailWithMsg("用户名重复", c)
		return
	}

	user.Salt = fmt.Sprintf("%06d", rand.Int31())
	user.Password = utils.MakePassword(pwd, user.Salt)

	createUserDB := models.CreateUser(user)
	if err := createUserDB.Error; err != nil {
		FailWithMsg("创建失败", c)
		return
	}

	OkWithMsg("创建成功", c)
}

// DeleteUser
// @Tags 用户模块
// @Summary 删除用户
// @param id formData string false "id"
// @Success 200 {string} json{"code","msg"}
// @Router /user/deleteUser [post]
func DeleteUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.PostForm("id"))
	if id == 0 {
		FailWithMsg("id 不能为空!", c)
		return
	}
	db := models.DeleteUser(id)
	if err := db.Error; err != nil {
		FailWithMsg("删除失败", c)
		return
	}
	OkWithMsg("删除成功", c)
}

// UpdateUser
// @Tags 用户模块
// @Summary 更新用户
// @param id formData string false "id"
// @param name formData string false "name"
// @param password formData string false "password"
// @param phone formData string false "phone"
// @param email formData string false "email"
// @Success 200 {string} json{"code","message"}
// @Router /user/updateUser [post]
func UpdateUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.PostForm("id"))
	if id == 0 {
		FailWithMsg("id 不能为空!", c)
		return
	}
	user.ID = uint(id)
	user.Name = c.PostForm("name")
	user.Password = c.PostForm("password")
	user.Phone = c.PostForm("phone")
	user.Email = c.PostForm("email")
	if _, err := govalidator.ValidateStruct(user); err != nil {
		FailWithMsg("修改参数不匹配", c)
		return
	}
	if findUser := models.FindUserByName(user.Name, id); findUser != nil {
		FailWithMsg("用户名重复", c)
		return
	}
	if findUser := models.FindUserByPhone(user.Phone, id); findUser != nil {
		FailWithMsg("手机号重复", c)
		return
	}
	if findUser := models.FindUserByEmail(user.Email, id); findUser != nil {
		FailWithMsg("邮箱重复", c)
		return
	}

	if findUser := models.FindUserByID(int(user.ID)); findUser == nil {
		FailWithMsg("无效用户", c)
		return
	} else {
		user.Password = utils.MakePassword(user.Password, findUser.Salt)
	}

	db := models.UpdateUser(user)
	if err := db.Error; err != nil {
		FailWithMsg("更新失败", c)
		return
	}
	OkWithMsg("更新成功", c)
}
