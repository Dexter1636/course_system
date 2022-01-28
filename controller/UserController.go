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

}

func (ctl UserController) Member(c *gin.Context) {

}

func (ctl UserController) List(c *gin.Context) {

}

func (ctl UserController) Update(c *gin.Context) {

}

func (ctl UserController) Delete(c *gin.Context) {

}
