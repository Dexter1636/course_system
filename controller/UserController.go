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
	rp, _ := regexp.MatchString("^([0-9]|[a-z]|[A-Z])*$", req.Password)
	ru, _ := regexp.MatchString("^([a-z]|[A-Z])*$", req.Username)

	if (len(req.Password) > 20 || len(req.Password) < 4) ||
		(len(req.Nickname) < 4 || len(req.Nickname) > 20 || !rp) ||
		(len(req.Username) < 8 || len(req.Username) > 20 || !ru) {
		c.JSON(http.StatusOK, vo.CreateMemberResponse{
			Code: vo.ParamInvalid,
		})
		return
	}

	user = model.User{Uuid: 0, UserName: req.Username, NickName: req.Nickname,
		Password: req.Password, RoleId: string(req.UserType), Enabled: 1}
	if err := ctl.DB.Create(&user).Error; err != nil {
		fmt.pr
	} else {
		panic(err.Error())
	}

	if err := ctl.DB.Where("user_name = ?", req.Username).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
<<<<<<< Updated upstream
			ctl.DB.Create(&user)
=======

>>>>>>> Stashed changes
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
<<<<<<< Updated upstream
	c.JSON(http.StatusOK, vo.GetMemberResponse{Code: vo.OK,
		Data: vo.TMember{UserID: string(user.Uuid), Nickname: user.NickName, Username: user.UserName, UserType: vo.UserType(RoleID)}})
=======
	c.JSON(http.StatusOK, vo.GetMemberResponse{
		Code: vo.OK,
		Data: struct {
			UserID   string
			Nickname string
			Username string
			UserType vo.UserType
		}{UserID: strconv.FormatInt(user.Uuid, 10), Nickname: user.NickName, Username: user.UserName, UserType: vo.UserType(RoleID)},
	})
>>>>>>> Stashed changes
	return
}

func (ctl UserController) List(c *gin.Context) {
	var req vo.GetMemberListRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		panic(err.Error())
	}

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

	//检查用户已删除
	if user.Enabled == 0 {
		c.JSON(http.StatusOK, vo.DeleteMemberResponse{Code: vo.UserHasDeleted})
		return
	}

	//删除用户
	if err := ctl.DB.Model(&user).Update("enabled", "0").Error; err != nil {
		panic(err.Error())
	}

	c.JSON(http.StatusOK, vo.DeleteMemberResponse{Code: vo.OK})

}
