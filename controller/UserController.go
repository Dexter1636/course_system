package controller

import (
	"course_system/common"
	"course_system/model"
	"course_system/vo"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
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

}

func (ctl UserController) Member(c *gin.Context) {

}

func (ctl UserController) List(c *gin.Context) {

}

func (ctl UserController) Update(c *gin.Context) {

}

func (ctl UserController) Delete(c *gin.Context) {
	var req vo.DeleteMemberRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		panic(err.Error())
	}

	var user model.User
	//检查用户不存在
	if err := ctl.DB.Where("user_name = ?", req.UserID).Take(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, vo.DeleteMemberResponse{Code: vo.CourseNotExisted})
			return
		} else {
			panic(err.Error())
		}
	}

	//检查用户已删除
	if(user.Enabled == 0){
		c.JSON(http.StatusOK, vo.DeleteMemberResponse{Code: vo.UserHasDeleted})
		return
	}

	//删除用户
	if err := ctl.DB.Model(&user).Where("user_name = ?", req.UserID).Update("enabled","0").Error; err != nil {
		panic(err.Error())
	}

	c.JSON(http.StatusOK,vo.DeleteMemberResponse{Code: vo.OK})

}
