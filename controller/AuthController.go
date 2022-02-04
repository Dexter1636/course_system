package controller

import (
	"course_system/common"
	"course_system/model"
	"course_system/vo"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

type IAuthController interface {
	Login(c *gin.Context)
	Logout(c *gin.Context)
	WhoAmI(c *gin.Context)
}

type AuthController struct {
	DB *gorm.DB
}

//连接数据库
func NewAuthController() IAuthController {
	db := common.GetDB()
	return AuthController{DB: db}
}

//用户输入账号和密码后点击登录
//用户名或者密码错误均返回密码错误。
//WrongPassword      ErrNo = 5  // 密码错误
//登录成功后需要设置 Cookie，Cookie 名称为 camp-session。
//response: ErrNo, UserID
func (ctl AuthController) Login(c *gin.Context) {

	var req vo.LoginRequest
	var user model.User
	code := vo.OK

	//response, ErrNo, UserID
	defer func() {
		c.JSON(http.StatusOK, vo.LoginResponse{
			Code: code,
			Data: struct {
				UserID string
			}{strconv.FormatInt(user.Uuid, 10)},
		})
	}()

	//是否需要校验参数?  指校验用户名以及密码是否符合要求?
	//直接复制其他同学的user管理里的参数校验了
	if err := c.ShouldBindJSON(&req); err != nil {
		code = vo.ParamInvalid
		return
	}
	tmpStr := req.Password
	r1, _ := regexp.MatchString("^(\\w*[0-9]+\\w*[a-z]+\\w*[A-Z]+)|(\\w*[0-9]+\\w*[A-Z]+\\w*[a-z]+)$", tmpStr)
	r2, _ := regexp.MatchString("^(\\w*[a-z]+\\w*[0-9]+\\w*[A-Z]+)|(\\w*[a-z]+\\w*[A-Z]+\\w*[0-9]+)$", tmpStr)
	r3, _ := regexp.MatchString("^(\\w*[A-Z]+\\w*[a-z]+\\w*[0-9]+)|(\\w*[A-Z]+\\w*[0-9]+\\w*[a-z]+)$", tmpStr)
	ru, _ := regexp.MatchString("^([a-z]|[A-Z])*$", req.Username)
	rp := r1 || r2 || r3
	if (len(req.Password) > 20 || len(req.Password) < 8 || !rp) ||
		(len(req.Username) < 8 || len(req.Username) > 20 || !ru) {
		code = vo.ParamInvalid
		return
	}

	//查询, 用户名或者密码错误均返回密码错误, WrongPassword
	//根据用户名查询, 无用户
	if err := ctl.DB.Where("user_name = ?", req.Username).Take(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			code = vo.WrongPassword
			log.Println("login: no such user")
			return
		}
	}
	//用户已被删除, 文档中未要求, 但感觉应该加上这种情况
	if user.Enabled == 0 {
		code = vo.UserHasDeleted
		log.Println("login: user has deleted")
		return
	}
	//密码错误
	if user.Password != req.Password {
		code = vo.WrongPassword
		log.Println("login: wrong password")
		return
	}

	//设置cookie, 不知道对不对嗷
	c.SetCookie("camp-session", string(user.Uuid), 0, "/", "localhost", false, true)

}

//当用户点击退出按钮，销毁当前用户 Session 和认证 Cookie
//登出后清除相应的 Cookie。
func (ctl AuthController) Logout(c *gin.Context) {
	//var user model.User
	code := vo.OK

	//response, ErrNo, user
	defer func() {
		c.JSON(http.StatusOK, vo.WhoAmIResponse{
			Code: code,
		})
	}()

	//无cookie, 需要登录
	_, err := c.Cookie("camp-session")
	if err != nil {
		code = vo.LoginRequired
		log.Println("logout: no cookie, login required")
		return
	}
	//将cookie的maxage设置为-1
	c.SetCookie("camp-session", "", -1, "/", "localhost", false, true)
}

//登录后访问个人信息页可以查看自己的信息，包括用户ID、用户名称、用户昵称。
//LoginRequired      ErrNo = 6  // 用户未登录
//通过cookie查看
func (ctl AuthController) WhoAmI(c *gin.Context) {
	var user model.User
	code := vo.OK

	//response, ErrNo, user
	defer func() {
		RoleID, _ := strconv.Atoi(user.RoleId)
		c.JSON(http.StatusOK, vo.WhoAmIResponse{
			Code: code,
			Data: vo.TMember{
				UserID:   strconv.FormatInt(user.Uuid, 10),
				Nickname: user.NickName,
				Username: user.UserName,
				UserType: vo.UserType(RoleID),
			},
		})
	}()

	cookie, err := c.Cookie("camp-session")
	//无cookie, 需要登录
	if err != nil {
		code = vo.LoginRequired
		log.Println("WhoAmI: no cookie, loginrequired")
		return
	}
	//有cookie, 根据存的Uuid获取信息
	uuidT, err := strconv.ParseInt(cookie, 10, 64)
	if err := ctl.DB.Where("uuid = ?", uuidT).Take(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			code = vo.UserNotExisted
			log.Println("WhoAmI: uuid not existed")
			return
		} else {
			panic(err.Error())
		}
	}
}
