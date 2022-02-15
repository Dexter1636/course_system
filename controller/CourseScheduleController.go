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
	"log"
	"net/http"
	"strconv"
)

type ICourseScheduleController interface {
	Bind(c *gin.Context)
	Unbind(c *gin.Context)
	Get(c *gin.Context)
	Schedule(c *gin.Context)
}
type CourseScheduleController struct {
	DB  *gorm.DB
	RDB *redis.Client
	Ctx context.Context
}

func NewCourseScheduleController() ICourseScheduleController {
	return CourseScheduleController{DB: common.GetDB(), RDB: common.GetRDB(), Ctx: common.GetCtx()}
}
func (ctl CourseScheduleController) Bind(c *gin.Context) {
	var req vo.BindCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, vo.BindCourseResponse{Code: vo.ParamInvalid})
		return
	}
	var sample model.Course
	number2, err1 := strconv.ParseInt(req.CourseID, 10, 64)
	if err1 != nil {
		c.JSON(http.StatusOK, vo.BindCourseResponse{Code: vo.ParamInvalid})
		return
	}
	number, err2 := strconv.ParseInt(req.TeacherID, 10, 64)
	if err2 != nil {
		c.JSON(http.StatusOK, vo.BindCourseResponse{Code: vo.ParamInvalid})
		return
	}
	log.Println("Bind ", req.CourseID, " ", req.TeacherID)
	val, err4 := ctl.RDB.Get(ctl.Ctx, fmt.Sprintf("course:%s", req.CourseID)).Result()
	if err4 == redis.Nil {
		//course not exist
		log.Println("Bind Case 1")
		c.JSON(http.StatusOK, vo.BindCourseResponse{Code: vo.CourseNotExisted})
		return
	} else if err4 != nil {
		//Redis错误
		log.Println("Bind Case 2")
		c.JSON(http.StatusOK, vo.BindCourseResponse{Code: vo.UnknownError})
		return
	} else {
		log.Println("Bind Case 3")
		if err4 := json.Unmarshal([]byte(val), &sample); err4 != nil {
			log.Println("Bind Case 4")
			//JSON解析错误
			c.JSON(http.StatusOK, vo.BindCourseResponse{Code: vo.UnknownError})
			return
		}
	}

	if sample.TeacherId != 0 {
		log.Println("Bind Case 5")
		c.JSON(http.StatusOK, vo.BindCourseResponse{Code: vo.CourseHasBound})
		return
	}
	sample.TeacherId = number
	val2, err3 := json.Marshal(sample)
	if err3 != nil {
		//JSON解析错误
		log.Println("Bind Case 6")
		c.JSON(http.StatusOK, vo.BindCourseResponse{Code: vo.UnknownError})
		return
	}

	//存入redis

	err5 := ctl.RDB.Set(ctl.Ctx, fmt.Sprintf("course:%s", req.CourseID), val2, 0).Err()
	if err5 != nil {
		log.Println("Bind Case 7")
		c.JSON(http.StatusOK, vo.BindCourseResponse{Code: vo.UnknownError})
		panic(err5.Error())
		return
	}
	//存入mysql
	if err := ctl.DB.Model(&model.Course{}).First(&sample, number2).Update("TeacherId", req.TeacherID).Error; err != nil {
		log.Println("Bind Case 8")
		c.JSON(http.StatusOK, vo.BindCourseResponse{Code: vo.UnknownError})
		return
	}
	c.JSON(http.StatusOK, vo.BindCourseResponse{Code: vo.OK})
	//a := ctl.DB.Model(&model.Course{}).First(&sample, number)
	//if a.RowsAffected == 0 {
	//	c.JSON(http.StatusOK, vo.BindCourseResponse{Code: vo.CourseNotExisted})
	//} else if sample.TeacherId != 0 {
	//	c.JSON(http.StatusOK, vo.BindCourseResponse{Code: vo.CourseHasBound})
	//} else {
	//	ctl.DB.Model(&model.Course{}).First(&sample, number).Update("TeacherId", req.TeacherID)
	//	c.JSON(http.StatusOK, vo.BindCourseResponse{Code: vo.OK})
	//}
}
func (ctl CourseScheduleController) Unbind(c *gin.Context) {
	var req vo.UnbindCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, vo.UnbindCourseResponse{Code: vo.ParamInvalid})
		return
	}
	var sample model.Course
	number, err1 := strconv.ParseInt(req.CourseID, 10, 64)
	if err1 != nil {
		c.JSON(http.StatusOK, vo.UnbindCourseResponse{Code: vo.ParamInvalid})
		return
	}
	number2, err2 := strconv.ParseInt(req.TeacherID, 10, 64)
	if err2 != nil {
		c.JSON(http.StatusOK, vo.UnbindCourseResponse{Code: vo.ParamInvalid})
		return
	}
	log.Println("UnBind ", req.CourseID, " ", req.TeacherID)
	val, err4 := ctl.RDB.Get(ctl.Ctx, fmt.Sprintf("course:%s", req.CourseID)).Result()
	if err4 == redis.Nil {
		//course not exist
		log.Println("UnBind Case 1")
		c.JSON(http.StatusOK, vo.UnbindCourseResponse{Code: vo.CourseNotExisted})
		return
	} else if err4 != nil {
		//Redis错误
		log.Println("UnBind Case 2")
		c.JSON(http.StatusOK, vo.UnbindCourseResponse{Code: vo.UnknownError})
		return
	} else {
		log.Println("UnBind Case 3")
		if err := json.Unmarshal([]byte(val), &sample); err != nil {
			log.Println("UnBind Case 4")
			//JSON解析错误
			c.JSON(http.StatusOK, vo.UnbindCourseResponse{Code: vo.UnknownError})
			return
		}
	}
	if sample.TeacherId == 0 {
		log.Println("UnBind Case 5")
		c.JSON(http.StatusOK, vo.UnbindCourseResponse{Code: vo.CourseNotBind})
		return
	}
	if sample.TeacherId != number2 {
		log.Println("UnBind Case 6")
		c.JSON(http.StatusOK, vo.UnbindCourseResponse{Code: vo.UserNotExisted})
		return
	}
	sample.TeacherId = 0
	val2, err3 := json.Marshal(sample)
	if err3 != nil {
		//JSON解析错误
		log.Println("UnBind Case 7")
		c.JSON(http.StatusOK, vo.UnbindCourseResponse{Code: vo.UnknownError})
		return
	}
	//存入redis
	err := ctl.RDB.Set(ctl.Ctx, fmt.Sprintf("course:%s", req.CourseID), val2, 0).Err()
	if err != nil {
		log.Println("UnBind Case 8")
		c.JSON(http.StatusOK, vo.UnbindCourseResponse{Code: vo.UnknownError})
		panic(err.Error())
		return
	}
	//存入mysql
	if err := ctl.DB.Model(&model.Course{}).First(&sample, number).Update("TeacherId", 0).Error; err != nil {
		log.Println("UnBind Case 9")
		c.JSON(http.StatusOK, vo.UnbindCourseResponse{Code: vo.UnknownError})
		panic(err.Error())
		return
	}
	c.JSON(http.StatusOK, vo.UnbindCourseResponse{Code: vo.OK})
	//a := ctl.DB.Model(&model.Course{}).First(&sample, number)
	//if a.RowsAffected == 0 {
	//	c.JSON(http.StatusOK, vo.UnbindCourseResponse{Code: vo.CourseNotExisted})
	//} else if sample.TeacherId == 0 {
	//	c.JSON(http.StatusOK, vo.UnbindCourseResponse{Code: vo.CourseNotBind})
	//} else {
	//	ctl.DB.Model(&model.Course{}).First(&sample, number).Update("TeacherId", 0)
	//	c.JSON(http.StatusOK, vo.UnbindCourseResponse{Code: vo.OK})
	//}
}
func (ctl CourseScheduleController) Get(c *gin.Context) {
	var req vo.GetTeacherCourseRequest
	//if err := c.ShouldBindJSON(&req); err != nil {
	//	c.JSON(http.StatusOK, vo.GetTeacherCourseResponse{Code: vo.ParamInvalid})
	//	return
	//}
	req.TeacherID = c.Query("TeacherID")
	number, err := strconv.ParseInt(req.TeacherID, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, vo.GetTeacherCourseResponse{Code: vo.ParamInvalid})
		return
	}
	var sample model.User
	if err := ctl.DB.Model(&model.User{}).First(&sample, number).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, vo.GetTeacherCourseResponse{Code: vo.UserNotExisted})
			return
		} else {
			c.JSON(http.StatusOK, vo.GetTeacherCourseResponse{Code: vo.UnknownError})
			return
		}
	} else if sample.RoleId != "3" {
		c.JSON(http.StatusOK, vo.GetTeacherCourseResponse{Code: vo.UserNotExisted})
		return
	}
	var rows []model.Course
	var ans vo.GetTeacherCourseResponse
	result := ctl.DB.Model(&model.Course{}).Where("teacher_id = ?", number).Find(&rows)
	ans.Data.CourseList = make([]*vo.TCourse, result.RowsAffected)
	for i := 0; i < int(result.RowsAffected); i++ {
		ans.Data.CourseList[i] = new(vo.TCourse)
		ans.Data.CourseList[i].CourseID = strconv.FormatInt(rows[i].Id, 10)
		ans.Data.CourseList[i].Name = rows[i].Name
		ans.Data.CourseList[i].TeacherID = strconv.FormatInt(rows[i].TeacherId, 10)
	}
	ans.Code = vo.OK
	c.JSON(http.StatusOK, ans)
}

type node struct {
	to, nxt int
}

var tot int = 0
var a []node
var q []int
var v []bool
var match []int
var tid []string
var cid []string

func add(x int, y int) {
	tot++
	a[q[x]].nxt = tot
	q[x] = tot
	a[q[x]].to = y
	tot++
	a[q[y]].nxt = tot
	q[y] = tot
	a[q[y]].to = x
}
func dfs(x int) bool {
	var p int = x
	for a[p].nxt != 0 {
		p = a[p].nxt
		if v[a[p].to] == false {
			v[a[p].to] = true
			if match[a[p].to] == 0 || dfs(match[a[p].to]) == true {
				match[a[p].to] = x
				match[x] = a[p].to
				return true
			}
		}
	}
	return false
}
func (ctl CourseScheduleController) Schedule(c *gin.Context) {
	var req vo.ScheduleCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, vo.ScheduleCourseResponse{
			Code: vo.ParamInvalid,
			Data: nil,
		})
		return
	}
	var tnum map[string]int
	var cnum map[string]int
	tnum = make(map[string]int)
	cnum = make(map[string]int)
	tid = make([]string, 1, len(req.TeacherCourseRelationShip)+10)
	cid = make([]string, 1, len(req.TeacherCourseRelationShip)+10)

	var n, m, nums int = 0, 0, 0
	for i, j := range req.TeacherCourseRelationShip {
		n++
		tnum[i] = n
		tid = append(tid, i)
		for k := 0; k < len(j); k++ {
			nums += 2
			x := j[k]
			value, ok := cnum[x]
			if !ok {
				m++
				cnum[x] = m
				cid = append(cid, x)
			}
			if value == value+1 {
				fmt.Println(value)
			}
		}
	}
	log.Println("schedule:", n, " ", m, " ", nums)
	a = make([]node, nums+10)
	q = make([]int, n+m+10)
	v = make([]bool, nums+10)
	match = make([]int, n+m+10)
	tot = n + m
	for i := 1; i <= n+m; i++ {
		q[i] = i
	}
	for i, j := range req.TeacherCourseRelationShip {
		for k := 0; k < len(j); k++ {
			var x, y int = tnum[i], cnum[j[k]]
			add(x, y+n)
		}
	}
	for i := 1; i <= n; i++ {
		for j := 1; j <= nums; j++ {
			v[j] = false
		}
		dfs(i)
	}
	var ans vo.ScheduleCourseResponse
	ans.Data = make(map[string]string)
	ans.Code = vo.OK

	for i := 1; i <= n; i++ {
		if match[i] != 0 {
			ans.Data[tid[i]] = cid[match[i]-n]
		}
	}
	c.JSON(http.StatusOK, ans)
}
