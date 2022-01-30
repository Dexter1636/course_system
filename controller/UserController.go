package controller

import (
	"course_system/common"
	"course_system/model"
	"course_system/vo"
	"errors"
	"github.com/gin-gonic/gin"
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
	DB *gorm.DB
}

func NewUserController() IUserController {
	db := common.GetDB()
	return UserController{DB: db}
}

func (ctl UserController) Create(c *gin.Context) {
	var req vo.CreateMemberRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		panic(err.Error())
	}
	//req = vo.CreateMemberRequest{Nickname: "alex", UserType: 1,
	//	Username: "alexander", Password: "12345678"}
	var user model.User

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
		c.JSON(http.StatusOK, vo.CreateMemberResponse{
			Code: vo.ParamInvalid,
		})
		return
	}

	if err := ctl.DB.First("user_name = ?", req.Username).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctl.DB.Create(&user)
			c.JSON(http.StatusOK, vo.CreateMemberResponse{
				Code: vo.OK,
				Data: struct{ UserID string }{UserID: string(user.Uuid)},
			})
		} else {
			c.JSON(http.StatusOK, vo.CreateMemberResponse{
				Code: vo.UserHasExisted,
			})
			return
		}
	} else {
		panic(err.Error())
	}
}

func (ctl UserController) Member(c *gin.Context) {
	var req vo.GetMemberRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		panic(err.Error())
	}

	var user model.User

	//检查用户不存在
	if err := ctl.DB.Where("uuid = ?", req.UserID).Take(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, vo.GetMemberResponse{Code: vo.UserNotExisted})
			return
		} else {
			panic(err.Error())
		}
	}

	//检查用户已删除
	if user.Enabled == 0 {
		c.JSON(http.StatusOK, vo.GetMemberResponse{Code: vo.UserHasDeleted})
		return
	}

	//返回TMember
	RoleID, _ := strconv.Atoi(user.RoleId)
	c.JSON(http.StatusOK, vo.GetMemberResponse{Code: vo.OK,
		Data: vo.TMember{UserID: string(user.Uuid), Nickname: user.NickName, Username: user.UserName, UserType: vo.UserType(RoleID)}})
	return
}

func (ctl UserController) List(c *gin.Context) {
	var req vo.GetMemberListRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		panic(err.Error())
	}

	//查询数据库
	var users []model.User
	if err := ctl.DB.Offset(req.Offset).Limit(req.Limit).Find(&users).Error; err != nil {
		panic(err.Error())
	}

	//获取数据
	var MemberList []vo.TMember
	for i := 0; i < len(users); i++ {
		UserType, err := strconv.Atoi(users[i].RoleId)
		if err != nil {
			panic(err.Error())
		}
		MemberList = append(MemberList, vo.TMember{
			strconv.FormatInt(users[i].Uuid, 10), users[i].NickName, users[i].UserName, vo.UserType(UserType)})
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
		panic(err.Error())
	}

	//检查数据合法性
	if len(req.Nickname) < 4 || len(req.Nickname) > 20 {
		c.JSON(http.StatusOK, vo.UpdateMemberResponse{Code: vo.ParamInvalid})
		return
	}

	var user model.User
	//检查用户不存在
	if err := ctl.DB.Take(&user, req.UserID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, vo.UpdateMemberResponse{Code: vo.UserNotExisted})
			return
		} else {
			panic(err.Error())
		}
	}

	//检查用户已删除
	if user.Enabled == 0 {
		c.JSON(http.StatusOK, vo.UpdateMemberResponse{Code: vo.UserHasDeleted})
		return
	}

	//修改用户名
	if err := ctl.DB.Model(&user).Update("nick_name", req.Nickname).Error; err != nil {
		panic(err.Error())
	}

	c.JSON(http.StatusOK, vo.UpdateMemberResponse{Code: vo.OK})

}

func (ctl UserController) Delete(c *gin.Context) {
	var req vo.DeleteMemberRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		panic(err.Error())
	}

	var user model.User
	//检查用户不存在
	if err := ctl.DB.Take(&user, req.UserID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, vo.DeleteMemberResponse{Code: vo.UserNotExisted})
			return
		} else {
			panic(err.Error())
		}
	}
	//text
	//检查用户已删除
	if user.Enabled == 0 {
		c.JSON(http.StatusOK, vo.DeleteMemberResponse{Code: vo.UserHasDeleted})
		return
	}

	//删除用户，将状态设置为0
	if err := ctl.DB.Model(&user).Update("enabled", "0").Error; err != nil {
		panic(err.Error())
	}

	c.JSON(http.StatusOK, vo.DeleteMemberResponse{Code: vo.OK})

}
