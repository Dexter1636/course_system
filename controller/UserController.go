package controller

import (
	"course_system/common"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

	var user model.User

	if err := ctl.DB.Where("user_name= ?", req.Username).Take(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			u := &model.User{Uuid: 0, UserName: req.Username, NickName: req.Nickname,
				Password: req.Password, RoleId: string(req.UserType), Enabled: 1}
			if err := ctl.DB.Create(u).Error; err != nil {
				panic(err.Error())
			}
		} else {
			panic(err.Error())
		}
	}
}

func (ctl UserController) Member(c *gin.Context) {
	var req vo.GetMemberRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		panic(err.Error())
	}
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

	var user model.User
	//检查用户不存在
	if err := ctl.DB.Where("uuid = ?", req.UserID).Take(&user).Error; err != nil {
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
	if err := ctl.DB.Model(&user).Where("uuid = ?", req.UserID).Update("nick_name", req.Nickname).Error; err != nil {
		panic(err.Error())
	}

	c.JSON(http.StatusOK, vo.DeleteMemberResponse{Code: vo.OK})

}

func (ctl UserController) Delete(c *gin.Context) {
	var req vo.DeleteMemberRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		panic(err.Error())
	}

	var user model.User
	//检查用户不存在
	if err := ctl.DB.Where("uuid = ?", req.UserID).Take(&user).Error; err != nil {
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
	if err := ctl.DB.Model(&user).Where("uuid = ?", req.UserID).Update("enabled", "0").Error; err != nil {
		panic(err.Error())
	}

	c.JSON(http.StatusOK, vo.DeleteMemberResponse{Code: vo.OK})

}
