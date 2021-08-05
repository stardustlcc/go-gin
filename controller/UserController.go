package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"godev/common"
	"godev/dto"
	"godev/model"
	"godev/response"
	"godev/util"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

func Register(c *gin.Context)  {

	DB := common.GetDb()

	//var requestMap = make(map[string]string)
	//json.NewDecoder(c.Request.Body).Decode(&requestMap)

	var requestUser = model.User{}
	//json.NewDecoder(c.Request.Body).Decode(&requestUser)
	c.Bind(&requestUser)

	//获取参数
	name 		:= requestUser.Name
	telephone	:= requestUser.Telephone
	password 	:= requestUser.Password

	//数据校验
	if len(telephone) != 11 {
		//gin.H 可以直接换位 map map[string]interface{}{"code":422, "msg":"手机号必须为11为"}
		//gin.H 为别名
		c.JSON(http.StatusUnprocessableEntity, gin.H{"code":422, "msg":"手机号必须为11为" })
		return
	}
	if len(password) < 6 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"code":422, "msg":"密码不能小于6位"})
		return
	}
	if len(name) == 0 {
		name = util.RandomString(10)
	}
	//判断手机号是否存在
	if !isTelephoneExists(DB, telephone) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"code":422,"msg":"用户已经存在"})
		return
	}

	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"code":500,"msg":"加密错误"})
		return
	}

	//创建用户
	newUser := model.User{
		Name:      name,
		Telephone: telephone,
		Password:  string(hasedPassword),
	}

	DB.Create(&newUser)

	//发放token
	token, err := common.ReleaseToken(newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code":500, "msg":"系统错误"})
		return
	}

	//返回结果
	response.Success(c, gin.H{"token": token}, "注册成功")

}

func isTelephoneExists(db *gorm.DB, tele string) bool {
	var user model.User
	db.Where("telephone = ?", tele).First(&user)
	if user.ID != 0 {
		return false
	}
	db.AutoMigrate(&model.User{})
	return true
}

func Login(c *gin.Context) {
	db := common.InitDb()

	var requestUser = model.User{}
	c.Bind(&requestUser)

	//获取参数
	telephone := requestUser.Telephone
	password := requestUser.Password

	//数据校验
	//数据校验
	if len(telephone) != 11 {
		//gin.H 可以直接换位 map map[string]interface{}{"code":422, "msg":"手机号必须为11为"}
		//gin.H 为别名
		response.Response(c,http.StatusUnprocessableEntity, 422, gin.H{}, "手机号必须为11为")
		return
	}
	if len(password) < 6 {
		response.Response(c,http.StatusUnprocessableEntity, 422, gin.H{}, "密码不能小于6位")
		return
	}
	//判断手机号是否存在
	var user model.User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID == 0 {
		response.Response(c,http.StatusUnprocessableEntity, 422, gin.H{}, "用户不存在")
		return
	}
	//判断密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		//c.JSON(http.StatusUnprocessableEntity, gin.H{"code":400, "msg":err})
		//return
	}

	//发送toKen
	token, err := common.ReleaseToken(user)
	if err != nil {
		response.Response(c,http.StatusInternalServerError, 500, gin.H{}, "系统错误")
		log.Printf("token generate error : %v", err)
		return
	}
	//返回数据
	response.Success(c,gin.H{"token":token}, "登录成功")
}

func Info(c *gin.Context)  {
	user, _ := c.Get("user")
	response.Success(c,gin.H{
		"user":dto.ToUserDto(user.(model.User)),
	}, "success")
}
