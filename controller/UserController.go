package controller

import (
	"context"
	"course_system/common"
	"course_system/model"
	"course_system/vo"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"net/http"
	"regexp"
	"strconv"
)

type IUserController interface {
	Create(c *gin.Context)
	Member(c *gin.Context)
	List(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type UserController struct {
	DB  *gorm.DB
	RDB *redis.Client
	Ctx context.Context
}

func NewUserController() IUserController {
	return UserController{DB: common.GetDB(), RDB: common.GetRDB(), Ctx: common.GetCtx()}
}

func (ctl UserController) Create(c *gin.Context) {
	var req vo.CreateMemberRequest
	var user, u model.User
	code := vo.OK

	defer func() {
		c.JSON(http.StatusOK, vo.CreateMemberResponse{
			Code: code,
			Data: struct{ UserID string }{UserID: strconv.FormatInt(user.Uuid, 10)},
		})
	}()

	if err := c.ShouldBindJSON(&req); err != nil {
		code = vo.UnknownError
		panic(err.Error())
		return
	}

	//权限检查
	//获取cookie
	cookie, err := c.Cookie("camp-session")
	if err != nil {
		code = vo.LoginRequired
		return
	}
	uuidT, err := strconv.ParseInt(cookie, 10, 64)
	//redis检查usertype
	val, err := ctl.RDB.Get(ctl.Ctx, fmt.Sprintf("user:%s", strconv.FormatInt(uuidT, 10))).Result()
	if err == redis.Nil {
		//用户不存在
		code = vo.UserNotExisted
		return
	} else if err != nil {
		//Redis错误
		code = vo.UnknownError
		panic(err.Error())
		return
	} else {
		var user model.User
		if err := json.Unmarshal([]byte(val), &user); err != nil {
			//JSON解析错误
			code = vo.UnknownError
			panic(err.Error())
			return
		}
		if user.RoleId != "1" {
			code = vo.PermDenied
			return
		}
	}

	//if err := ctl.DB.Where("uuid = ?", uuidT).Take(&user).Error; err != nil {
	//	if errors.Is(err, gorm.ErrRecordNotFound) {
	//		code = vo.UserNotExisted
	//		log.Println("Create: uuid not existed")
	//		return
	//	} else {
	//		code = vo.UnknownError
	//		panic(err.Error())
	//		return
	//	}
	//}
	//if user.UserName != "JudgeAdmin" {
	//	code = vo.PermDenied
	//	log.Println("Create: PermDenied")
	//	return
	//}

	//参数校验
	tmpStr := req.Password
	r1, _ := regexp.MatchString("^(\\w*[0-9]+\\w*[a-z]+\\w*[A-Z]+)|(\\w*[0-9]+\\w*[A-Z]+\\w*[a-z]+)$", tmpStr)
	r2, _ := regexp.MatchString("^(\\w*[a-z]+\\w*[0-9]+\\w*[A-Z]+)|(\\w*[a-z]+\\w*[A-Z]+\\w*[0-9]+)$", tmpStr)
	r3, _ := regexp.MatchString("^(\\w*[A-Z]+\\w*[a-z]+\\w*[0-9]+)|(\\w*[A-Z]+\\w*[0-9]+\\w*[a-z]+)$", tmpStr)
	ru, _ := regexp.MatchString("^([a-z]|[A-Z])*$", req.Username)
	rp := r1 || r2 || r3

	if (len(req.Password) > 20 || len(req.Password) < 8 || !rp) ||
		(len(req.Nickname) < 4 || len(req.Nickname) > 20) ||
		(len(req.Username) < 8 || len(req.Username) > 20 || !ru) ||
		(req.UserType > 3 || req.UserType < 1) {
		code = vo.ParamInvalid
		return
	}

	rid := strconv.FormatInt(int64(int(req.UserType)), 10)
	user = model.User{Uuid: 0, UserName: req.Username, NickName: req.Nickname,
		Password: req.Password, RoleId: rid, Enabled: 1}

	if err := ctl.DB.Where("user_name = ?", req.Username).Take(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctl.DB.Create(&user)
			val, err := json.Marshal(user)
			if err != nil {
				//JSON解析错误
				code = vo.UnknownError
				panic(err.Error())
				return
			}
			//存入redis
			err = ctl.RDB.Set(ctl.Ctx, fmt.Sprintf("user:%d", user.Uuid), val, 0).Err()
			if err != nil {
				code = vo.UnknownError
				panic(err.Error())
				return
			}
			return
		} else {
			code = vo.UnknownError
			panic(err.Error())
			return
		}
	}

	//用户已经存在
	code = vo.UserHasExisted
	return
}

func (ctl UserController) Member(c *gin.Context) {
	code := vo.OK
	var user model.User
	defer func() {
		RoleID, _ := strconv.Atoi(user.RoleId)
		c.JSON(http.StatusOK, vo.GetMemberResponse{
			Code: code,
			Data: struct {
				UserID   string
				Nickname string
				Username string
				UserType vo.UserType
			}{UserID: strconv.FormatInt(user.Uuid, 10), Nickname: user.NickName, Username: user.UserName, UserType: vo.UserType(RoleID)},
		})
	}()
	var req vo.GetMemberRequest

	req.UserID = c.Query("UserID")

	val, err := ctl.RDB.Get(ctl.Ctx, fmt.Sprintf("user:%s", req.UserID)).Result()
	if err == redis.Nil {
		//用户不存在
		code = vo.UserNotExisted
		return
	} else if err != nil {
		//Redis错误
		code = vo.UnknownError
		panic(err.Error())
		return
	} else {
		if err := json.Unmarshal([]byte(val), &user); err != nil {
			//JSON解析错误
			code = vo.UnknownError
			panic(err.Error())
			return
		}

		//检查用户已删除
		if user.Enabled == 0 {
			code = vo.UserHasDeleted
			return
		}

		//返回TMember
		return
	}
}

func (ctl UserController) List(c *gin.Context) {
	var req vo.GetMemberListRequest

	//读取数据
	val, err := strconv.Atoi(c.Query("Limit"))
	if err != nil {
		c.JSON(http.StatusOK, vo.GetMemberListResponse{
			Code: vo.UnknownError,
			Data: struct {
				MemberList []vo.TMember
			}{MemberList: []vo.TMember{}},
		})
		panic(err.Error())
		return
	}
	req.Limit = val
	val, err = strconv.Atoi(c.Query("Offset"))
	if err != nil {
		c.JSON(http.StatusOK, vo.GetMemberListResponse{
			Code: vo.UnknownError,
			Data: struct {
				MemberList []vo.TMember
			}{MemberList: []vo.TMember{}},
		})
		panic(err.Error())
		return
	}
	req.Offset = val

	//Limit忽略时，不能存在Offset参数
	if req.Limit <= 0 && req.Offset > 0 {
		c.JSON(http.StatusOK, vo.GetMemberListResponse{
			Code: vo.ParamInvalid,
			Data: struct {
				MemberList []vo.TMember
			}{MemberList: []vo.TMember{}},
		})
		return
	}

	//查询数据库
	var users []model.User
	if err := ctl.DB.Offset(req.Offset).Limit(req.Limit).Find(&users).Error; err != nil {
		c.JSON(http.StatusOK, vo.GetMemberListResponse{
			Code: vo.UnknownError,
			Data: struct {
				MemberList []vo.TMember
			}{MemberList: []vo.TMember{}},
		})
		panic(err.Error())
		return
	}

	//获取数据
	var MemberList []vo.TMember
	for i := 0; i < len(users); i++ {
		UserType, err := strconv.Atoi(users[i].RoleId)
		if err != nil {
			c.JSON(http.StatusOK, vo.GetMemberListResponse{
				Code: vo.UnknownError,
				Data: struct {
					MemberList []vo.TMember
				}{MemberList: []vo.TMember{}},
			})
			panic(err.Error())
			return
		}
		MemberList = append(MemberList, vo.TMember{
			UserID: strconv.FormatInt(users[i].Uuid, 10), Nickname: users[i].NickName, Username: users[i].UserName, UserType: vo.UserType(UserType)})
	}

	//防止返回NULL
	if len(MemberList) == 0 {
		MemberList = make([]vo.TMember, 0)
	}

	//返回参数
	c.JSON(http.StatusOK, vo.GetMemberListResponse{
		Code: vo.OK,
		Data: struct {
			MemberList []vo.TMember
		}{MemberList: MemberList},
	})
}

func (ctl UserController) Update(c *gin.Context) {
	var req vo.UpdateMemberRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, vo.UpdateMemberResponse{Code: vo.UnknownError})
		panic(err.Error())
		return
	}

	//检查数据合法性
	if len(req.Nickname) < 4 || len(req.Nickname) > 20 {
		c.JSON(http.StatusOK, vo.UpdateMemberResponse{Code: vo.ParamInvalid})
		return
	}

	val, err := ctl.RDB.Get(ctl.Ctx, fmt.Sprintf("user:%s", req.UserID)).Result()
	if err == redis.Nil {
		//用户不存在
		c.JSON(http.StatusOK, vo.UpdateMemberResponse{Code: vo.UserNotExisted})
		return
	} else if err != nil {
		//Redis错误
		c.JSON(http.StatusOK, vo.UpdateMemberResponse{Code: vo.UnknownError})
		panic(err.Error())
		return
	} else {
		var user model.User
		if err := json.Unmarshal([]byte(val), &user); err != nil {
			//JSON解析错误
			c.JSON(http.StatusOK, vo.UpdateMemberResponse{Code: vo.UnknownError})
			panic(err.Error())
			return
		}

		/*检查用户不存在(old)
		if err := ctl.DB.Take(&user, req.UserID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusOK, vo.UpdateMemberResponse{Code: vo.UserNotExisted})
				return
			} else {
				panic(err.Error())
			}
		}*/

		//检查用户已删除
		if user.Enabled == 0 {
			c.JSON(http.StatusOK, vo.UpdateMemberResponse{Code: vo.UserHasDeleted})
			return
		}

		//修改用户名
		user.NickName = req.Nickname
		val, err := json.Marshal(user)
		if err != nil {
			//JSON解析错误
			c.JSON(http.StatusOK, vo.UpdateMemberResponse{Code: vo.UnknownError})
			panic(err.Error())
			return
		}
		//存入redis
		err = ctl.RDB.Set(ctl.Ctx, fmt.Sprintf("user:%d", user.Uuid), val, 0).Err()
		if err != nil {
			c.JSON(http.StatusOK, vo.UpdateMemberResponse{Code: vo.UnknownError})
			panic(err.Error())
			return
		}
		//存入mysql
		if err := ctl.DB.Model(&user).Update("nick_name", req.Nickname).Error; err != nil {
			c.JSON(http.StatusOK, vo.UpdateMemberResponse{Code: vo.UnknownError})
			panic(err.Error())
			return
		}

		c.JSON(http.StatusOK, vo.UpdateMemberResponse{Code: vo.OK})
	}
}

func (ctl UserController) Delete(c *gin.Context) {
	var req vo.DeleteMemberRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, vo.DeleteMemberResponse{Code: vo.UnknownError})
		panic(err.Error())
		return
	}

	val, err := ctl.RDB.Get(ctl.Ctx, fmt.Sprintf("user:%s", req.UserID)).Result()
	if err == redis.Nil {
		//用户不存在
		c.JSON(http.StatusOK, vo.DeleteMemberResponse{Code: vo.UserNotExisted})
		return
	} else if err != nil {
		//Redis错误
		c.JSON(http.StatusOK, vo.DeleteMemberResponse{Code: vo.UnknownError})
		panic(err.Error())
		return
	} else {
		var user model.User
		if err := json.Unmarshal([]byte(val), &user); err != nil {
			//JSON解析错误
			c.JSON(http.StatusOK, vo.DeleteMemberResponse{Code: vo.UnknownError})
			panic(err.Error())
			return
		}

		//检查用户已删除
		if user.Enabled == 0 {
			c.JSON(http.StatusOK, vo.DeleteMemberResponse{Code: vo.UserHasDeleted})
			return
		}

		//删除用户，将状态设置为0
		user.Enabled = 0
		val, err := json.Marshal(user)
		if err != nil {
			//JSON解析错误
			c.JSON(http.StatusOK, vo.DeleteMemberResponse{Code: vo.UnknownError})
			panic(err.Error())
			return
		}
		//存入redis
		err = ctl.RDB.Set(ctl.Ctx, fmt.Sprintf("user:%d", user.Uuid), val, 0).Err()
		if err != nil {
			c.JSON(http.StatusOK, vo.DeleteMemberResponse{Code: vo.UnknownError})
			panic(err.Error())
			return
		}
		//存入mysql
		if err := ctl.DB.Model(&user).Update("enabled", "0").Error; err != nil {
			c.JSON(http.StatusOK, vo.DeleteMemberResponse{Code: vo.UnknownError})
			panic(err.Error())
			return
		}

		c.JSON(http.StatusOK, vo.DeleteMemberResponse{Code: vo.OK})
	}
}
