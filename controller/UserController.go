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
		resp := vo.CreateMemberResponse{
			Code: code,
			Data: struct{ UserID string }{UserID: strconv.FormatInt(user.Uuid, 10)},
		}
		c.JSON(http.StatusOK, resp)
		utils.LogReqRespBody(req, resp, "CreateMember")
	}()

	if err := c.ShouldBindJSON(&req); err != nil {
		code = vo.UnknownError
		//panic(err.Error())
		log.Println("CreateMember:ShouldBindJSON error")
		return
	}

	//权限检查
	//获取cookie
	//cookie, err := c.Cookie("camp-session")
	//if err != nil {
	//	code = vo.LoginRequired
	//	log.Println("CreateMember:Login Required")
	//	log.Println(err)
	//	return
	//}
	//获取session
	session, err := Store.Get(c.Request, "camp-session")
	if session.IsNew || err != nil {
		code = vo.LoginRequired
		log.Println("[CreateMember] : no session, 出错了，这里查不到session")
		log.Println(err)
		return
	} else {
		cookie := session.Values["UserType"].(string)
		uuidT, err := strconv.ParseInt(cookie, 10, 64)
		if uuidT != 1 {
			code = vo.PermDenied
			log.Println("CreateMember:PermDenied cause user not admin，非管理员的session")
			log.Println(err)
			return
		}
	}

	//redis检查usertype
	//val, err := ctl.RDB.Get(ctl.Ctx, fmt.Sprintf("user:%s", strconv.FormatInt(uuidT, 10))).Result()
	//if err == redis.Nil {
	//	//用户不存在
	//	code = vo.UserNotExisted
	//	log.Println("CreateMember:UserNotExisted while login check")
	//	return
	//} else if err != nil {
	//	//Redis错误
	//	code = vo.UnknownError
	//	log.Println("CreateMember:redis-error while login check")
	//	//panic(err.Error())
	//	return
	//} else {
	//	var userTmp model.UserTmp
	//	if err := json.Unmarshal([]byte(val), &userTmp); err != nil {
	//		//JSON解析错误
	//		code = vo.UnknownError
	//		log.Println("CreateMember:json-error while login check")
	//		log.Println(err)
	//		//panic(err.Error())
	//		return
	//	}
	//	if user.Enabled == 0 {
	//		code = vo.UserHasDeleted
	//		log.Println("CreateMember:UserHasDeleted While LoginCheck")
	//		return
	//	}
	//	if userTmp.RoleId != "1" {
	//		code = vo.PermDenied
	//		log.Println("CreateMember:PermDenied cause user not admin")
	//		return
	//	}
	//}

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
		log.Println("CreateMember:ParamInvalid")
		return
	}

	rid := strconv.FormatInt(int64(int(req.UserType)), 10)
	user = model.User{UserName: req.Username, NickName: req.Nickname,
		Password: req.Password, RoleId: rid, Enabled: 1}

	if err := ctl.DB.Where("user_name = ?", req.Username).Take(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctl.DB.Create(&user)
			val, err := json.Marshal(user)
			if err != nil {
				//JSON解析错误
				code = vo.UnknownError
				log.Println("CreateMember:JSON-error while creating")
				//panic(err.Error())
				return
			}
			//存入redis
			err = ctl.RDB.Set(ctl.Ctx, fmt.Sprintf("user:%d", user.Uuid), val, 0).Err()
			if err != nil {
				code = vo.UnknownError
				//panic(err.Error())
				log.Println("CreateMember:redis-error while creating")
				return
			}
			log.Println("CreateMember:Successfully create, userid:" + strconv.FormatInt(user.Uuid, 10))
			return
		} else {
			code = vo.UnknownError
			//panic(err.Error())
			log.Println("CreateMember:Unknown-error while creating")
			return
		}
	}

	//用户已经存在
	code = vo.UserHasExisted
	log.Println("CreateMember:UserExisted")
	return
}

func (ctl UserController) Member(c *gin.Context) {
	code := vo.OK
	var user model.User
	defer func() {
		RoleID, _ := strconv.Atoi(user.RoleId)
		resp := vo.GetMemberResponse{
			Code: code,
			Data: struct {
				UserID   string
				Nickname string
				Username string
				UserType vo.UserType
			}{UserID: strconv.FormatInt(user.Uuid, 10), Nickname: user.NickName, Username: user.UserName, UserType: vo.UserType(RoleID)}}
		c.JSON(http.StatusOK, resp)
	}()
	var req vo.GetMemberRequest

	req.UserID = c.Query("UserID")
	log.Print("GetMember: asking for uuid:" + req.UserID)

	val, err := ctl.RDB.Get(ctl.Ctx, fmt.Sprintf("user:%s", req.UserID)).Result()
	if err == redis.Nil {
		//用户不存在
		code = vo.UserNotExisted
		log.Print("GetMember:UserNotExisted")
		return
	} else if err != nil {
		//Redis错误
		code = vo.UnknownError
		log.Println("Member:redis-error")
		//panic(err.Error())
		return
	} else {
		if err := json.Unmarshal([]byte(val), &user); err != nil {
			//JSON解析错误
			code = vo.UnknownError
			//panic(err.Error())
			log.Println("Member:JSON-error")
			return
		}

		//检查用户已删除
		if user.Enabled == 0 {
			code = vo.UserHasDeleted
			log.Print("GetMember:UserHasDeleted")
			return
		}
		log.Print("GetMember:Return Successfully,username:" + user.UserName)
		//返回TMember
		return
	}
}

func (ctl UserController) List(c *gin.Context) {
	var req vo.GetMemberListRequest
	var MemberList []vo.TMember
	code := vo.OK

	defer func() {
		//防止返回NULL
		if len(MemberList) == 0 {
			MemberList = make([]vo.TMember, 0)
		}
		resp := vo.GetMemberListResponse{
			Code: code,
			Data: struct{ MemberList []vo.TMember }{MemberList: MemberList},
		}
		c.JSON(http.StatusOK, resp)
		utils.LogReqRespBody(req, fmt.Sprintf("MemberList len: %d", len(MemberList)), "List")
	}()

	//读取数据
	val, err := strconv.Atoi(c.Query("Limit"))
	if err != nil {
		code = vo.ParamInvalid
		log.Println("List:Limit参数错误")
		log.Println(err.Error())
		return
	}
	req.Limit = val
	val, err = strconv.Atoi(c.Query("Offset"))
	if err != nil {
		code = vo.ParamInvalid
		log.Println("List:Offset参数错误")
		log.Println(err.Error())
		return
	}
	req.Offset = val

	//Limit忽略时，不能存在Offset参数
	if req.Limit <= 0 && req.Offset > 0 {
		code = vo.ParamInvalid
		log.Println("List:Limit不存在，存在Offset")
		return
	}

	//查询数据库
	var users []model.User
	if err := ctl.DB.Offset(req.Offset).Limit(req.Limit).Find(&users).Error; err != nil {
		code = vo.UnknownError
		log.Println("List:查询数据库错误")
		log.Println(err.Error())
		return
	}

	//获取数据
	for i := 0; i < len(users); i++ {
		UserType, err := strconv.Atoi(users[i].RoleId)
		if err != nil {
			code = vo.UnknownError
			log.Println("List:Atoi(RoleId)错误")
			log.Println(err.Error())
			return
		}
		MemberList = append(MemberList, vo.TMember{
			UserID: strconv.FormatInt(users[i].Uuid, 10), Nickname: users[i].NickName, Username: users[i].UserName, UserType: vo.UserType(UserType)})
	}

	//返回参数
	log.Println("List:查询成功，Code=0")
}

func (ctl UserController) Update(c *gin.Context) {
	var req vo.UpdateMemberRequest
	code := vo.OK

	defer func() {
		resp := vo.UpdateMemberResponse{
			Code: code,
		}
		c.JSON(http.StatusOK, resp)
		utils.LogReqRespBody(req, resp, "Update")
	}()

	if err := c.ShouldBindJSON(&req); err != nil {
		code = vo.UnknownError
		log.Println("Update:参数BindJSON错误")
		log.Println(err.Error())
		return
	}

	//检查数据合法性
	if len(req.Nickname) < 4 || len(req.Nickname) > 20 {
		code = vo.ParamInvalid
		log.Println("Update:Nickname不合法")
		return
	}

	val, err := ctl.RDB.Get(ctl.Ctx, fmt.Sprintf("user:%s", req.UserID)).Result()
	if err == redis.Nil {
		//用户不存在
		code = vo.UserNotExisted
		log.Println("Update:用户不存在")
		return
	} else if err != nil {
		//Redis错误
		code = vo.UnknownError
		log.Println("Update:Redis get错误")
		log.Println(err.Error())
		return
	} else {
		var user model.User
		if err := json.Unmarshal([]byte(val), &user); err != nil {
			//JSON解析错误
			code = vo.UnknownError
			log.Println("Update:Unmarshal解析错误")
			log.Println(err.Error())
			return
		}

		//检查用户已删除
		if user.Enabled == 0 {
			code = vo.UserHasDeleted
			log.Println("Update:用户已删除")
			return
		}

		//修改用户名
		user.NickName = req.Nickname
		val, err := json.Marshal(user)
		if err != nil {
			//JSON解析错误
			code = vo.UnknownError
			log.Println("Update:Marshal解析错误")
			log.Println(err.Error())
			return
		}
		//存入redis
		err = ctl.RDB.Set(ctl.Ctx, fmt.Sprintf("user:%d", user.Uuid), val, 0).Err()
		if err != nil {
			code = vo.UnknownError
			log.Println("Update:Redis set错误")
			log.Println(err.Error())
			return
		}
		//存入mysql
		if err := ctl.DB.Model(&user).Update("nick_name", req.Nickname).Error; err != nil {
			code = vo.UnknownError
			log.Println("Update:Mysql写入错误")
			log.Println(err.Error())
			return
		}

		log.Println("Update:修改成功，Code=0")
	}
}

func (ctl UserController) Delete(c *gin.Context) {
	var req vo.DeleteMemberRequest
	code := vo.OK

	defer func() {
		resp := vo.DeleteMemberResponse{
			Code: code,
		}
		c.JSON(http.StatusOK, resp)
		utils.LogReqRespBody(req, resp, "Delete")
	}()

	if err := c.ShouldBindJSON(&req); err != nil {
		code = vo.UnknownError
		log.Println("Delete:参数BindJSON错误")
		log.Println(err.Error())
		return
	}

	val, err := ctl.RDB.Get(ctl.Ctx, fmt.Sprintf("user:%s", req.UserID)).Result()
	if err == redis.Nil {
		//用户不存在
		code = vo.UserNotExisted
		log.Println("Delete:用户不存在")
		return
	} else if err != nil {
		//Redis错误
		code = vo.UnknownError
		log.Println("Delete:Redis get错误")
		log.Println(err.Error())
		return
	} else {
		var user model.User
		if err := json.Unmarshal([]byte(val), &user); err != nil {
			//JSON解析错误
			code = vo.UnknownError
			log.Println("Delete:Unmarshal解析错误")
			log.Println(err.Error())
			return
		}

		//检查用户已删除
		if user.Enabled == 0 {
			code = vo.UserHasDeleted
			log.Println("Delete:用户已删除")
			return
		}

		//删除用户，将状态设置为0
		user.Enabled = 0
		val, err := json.Marshal(user)
		if err != nil {
			//JSON解析错误
			code = vo.UnknownError
			log.Println("Delete:Marshal解析错误")
			log.Println(err.Error())
			return
		}
		//存入redis
		err = ctl.RDB.Set(ctl.Ctx, fmt.Sprintf("user:%d", user.Uuid), val, 0).Err()
		if err != nil {
			code = vo.UnknownError
			log.Println("Delete:Redis set错误")
			log.Println(err.Error())
			return
		}
		//存入mysql
		if err := ctl.DB.Model(&user).Update("enabled", "0").Error; err != nil {
			code = vo.UnknownError
			log.Println("Delete:Mysql写入错误")
			log.Println(err.Error())
			return
		}

		log.Println("Delete:删除成功，Code=0")
	}
}
