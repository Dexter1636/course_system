package controller

import (
	"course_system/common"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IAuthController interface {
	Login(c *gin.Context)
	Logout(c *gin.Context)
	WhoAmI(c *gin.Context)
}

type AuthController struct {
	DB *gorm.DB
}

func NewAuthController() IAuthController {
	db := common.GetDB()
	return AuthController{DB: db}
}

func (ctl AuthController) Login(c *gin.Context) {
//test 220128
}

func (ctl AuthController) Logout(c *gin.Context) {

}

func (ctl AuthController) WhoAmI(c *gin.Context) {

}
