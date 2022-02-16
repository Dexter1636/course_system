package controller

import (
	"context"
	"course_system/common"
	"course_system/model"
	"course_system/utils"
	"course_system/vo"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
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
	DB  *gorm.DB
	RDB *redis.Client
	Ctx context.Context
}

//连接数据库
func NewAuthController() IAuthController {
	return AuthController{DB: common.GetDB(), RDB: common.GetRDB(), Ctx: common.GetCtx()}
}

//用户输入账号和密码后点击登录
//用户名或者密码错误均返回密码错误。
//WrongPassword      ErrNo = 5  // 密码错误
//登录成功后需要设置 Cookie，Cookie 名称为 camp-session。
//response: ErrNo, UserID
func (ctl AuthController) Login(c *gin.Context) {
	//POST方法，传参使用json，无需修改
	var req vo.LoginRequest
	var user model.User
	var resp vo.LoginResponse
	code := vo.OK

	//response, ErrNo, UserID
	defer func() {
		if code == vo.OK {
			resp = vo.LoginResponse{
				Code: code,
				Data: struct {
					UserID string
				}{strconv.FormatInt(user.Uuid, 10)},
			}
		} else {
			resp = vo.LoginResponse{
				Code: code,
				Data: struct {
					UserID string
				}{""},
			}
		}
		c.JSON(http.StatusOK, resp)
		utils.LogReqRespBody(req, resp, "Login")
	}()

	//校验参数， 用户名、密码是否符合要求
	if err := c.ShouldBindJSON(&req); err != nil {
		code = vo.WrongPassword //修改返回的错误码,220208
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
		code = vo.WrongPassword //修改返回的错误码,220208
		log.Println("[login]: ParamInvalid")
		return
	}

	//根据用户名查询, 无用户
	if err := ctl.DB.Where("user_name = ?", req.Username).Take(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			code = vo.WrongPassword
			log.Println("[login]: no such user")
			return
		} else {
			code = vo.UnknownError
			log.Println("[login]: UnkonwnError")
			return
		}
	}
	//用户已被删除, 文档中未要求, 但感觉应该加上这种情况
	if user.Enabled == 0 {
		code = vo.WrongPassword //修改返回的错误码,220208
		log.Println("[login]: user has deleted")
		return
	}
	//密码错误
	if user.Password != req.Password {
		code = vo.WrongPassword
		log.Println("[login]: wrong password")
		return
	}

	//设置cookie，存储uuid
	c.SetCookie("camp-session", strconv.FormatInt(user.Uuid, 10), 0, "/", "", false, true)
	log.Println("[login]:successfully login, uuid:" + strconv.FormatInt(user.Uuid, 10))
}

//当用户点击退出按钮，销毁当前用户 Session 和认证 Cookie
//登出后清除相应的 Cookie。
func (ctl AuthController) Logout(c *gin.Context) {
	code := vo.OK

	//response, ErrNo
	defer func() {
		log.Println("[logout] ErrNo: ", code)
		c.JSON(http.StatusOK, vo.LogoutResponse{
			Code: code,
		})
	}()

	//无cookie, 需要登录
	ck, err := c.Cookie("camp-session")
	if err != nil {
		code = vo.LoginRequired
		log.Println("[logout]: no cookie, login required")
		return
	}
	//将cookie的maxage设置为-1
	log.Println("[logout]: cookievalue: " + ck)
	c.SetCookie("camp-session", "", -1, "/", "", false, true)
}

//登录后访问个人信息页可以查看自己的信息，包括用户ID、用户名称、用户昵称。
//LoginRequired      ErrNo = 6  // 用户未登录
//通过cookie查看
func (ctl AuthController) WhoAmI(c *gin.Context) {
	var user model.User
	var resp vo.WhoAmIResponse
	code := vo.OK

	//response, ErrNo, user
	defer func() {
		log.Println("[WhoAmI] ErrNo:  ", code)
		if code == vo.OK {
			RoleID, _ := strconv.Atoi(user.RoleId)
			resp = vo.WhoAmIResponse{
				Code: code,
				Data: vo.TMember{
					UserID:   strconv.FormatInt(user.Uuid, 10),
					Nickname: user.NickName,
					Username: user.UserName,
					UserType: vo.UserType(RoleID),
				},
			}
		} else {
			resp = vo.WhoAmIResponse{
				Code: code,
				Data: vo.TMember{
					UserID:   "",
					Nickname: "",
					Username: "",
					UserType: 0,
				},
			}
		}
		c.JSON(http.StatusOK, resp)
		utils.LogBody(resp, "WhoAmI.Resp")
	}()

	cookie, err := c.Cookie("camp-session")
	//无cookie, 需要登录
	if err != nil {
		code = vo.LoginRequired
		log.Println(err) //打印err
		log.Println("[WhoAmI]: no cookie, loginrequired")
		return
	}
	//有cookie, 根据存的Uuid获取信息
	//uuidT, err := strconv.ParseInt(cookie, 10, 64)
	//if err := ctl.DB.Where("uuid = ?", uuidT).Take(&user).Error; err != nil {
	//	if errors.Is(err, gorm.ErrRecordNotFound) {
	//		code = vo.UserNotExisted
	//		log.Println("WhoAmI: uuid not existed")
	//		return
	//	} else {
	//		panic(err.Error())
	//	}
	//}
	//rdbreq := fmt.Sprintf("user:%s", cookie)
	val, err := ctl.RDB.Get(ctl.Ctx, fmt.Sprintf("user:%s", cookie)).Result()
	if err == redis.Nil {
		//用户不存在
		code = vo.UserNotExisted
		log.Println("[WhoAmI]: UserNotExisted")
		return
	} else if err != nil {
		//Redis错误
		code = vo.UnknownError
		log.Println("[WhoAmI]: Redis Error")
		log.Println(err) //打印err
		panic(err.Error())
		return
	} else {
		if err := json.Unmarshal([]byte(val), &user); err != nil {
			//JSON解析错误
			code = vo.UnknownError
			log.Println("[WhoAmI]: JSON error")
			panic(err.Error())
			return
		}

		log.Println("[WhoAmI]: Cookie Value: " + cookie)
		if user.Enabled == 0 {
			code = vo.UserHasDeleted
			log.Println("[WhoAmI]: USerHasDeleted")
			return
		}

	}
}
