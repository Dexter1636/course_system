package repository

import (
	"course_system/common"
	"course_system/model"
	"course_system/vo"
	"errors"
	"gorm.io/gorm"
	"log"
)

type ICourseRepository interface {
	GetCourseById(id int64, course *model.Course) (code vo.ErrNo)
	GetCourseListByStudentId(stuId int64, courseList *[]model.Course) (code vo.ErrNo)
	CreateCourse(course *model.Course) (code vo.ErrNo)
}

type CourseRepository struct {
	DB *gorm.DB
}

func NewCourseRepository() ICourseRepository {
	return CourseRepository{DB: common.GetDB()}
}

func (cr CourseRepository) GetCourseById(id int64, course *model.Course) (code vo.ErrNo) {
	if err := cr.DB.First(&course, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("repository.GetCourseById: CourseNotExisted")
			return vo.CourseNotExisted
		} else {
			panic(err.Error())
		}
	}
	return vo.OK
}

func (cr CourseRepository) GetCourseListByStudentId(stuId int64, courseList *[]model.Course) (code vo.ErrNo) {
	//subQuery := cr.DB.Table("sc").Select("course_id").Where("student_id", stuId)
	//if err := subQuery.Error; err != nil {
	//
	//}
	//if err := cr.DB.Find(&courseList, subQuery).Error; err != nil {
	//	panic(err.Error())
	//}
	if err := cr.DB.Raw("SELECT * FROM course WHERE id IN (SELECT course_id FROM sc WHERE student_id = ?)", stuId).Scan(&courseList).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("repository.GetCourseListByStudentId: CourseNotExisted")
			return vo.StudentHasNoCourse
		} else {
			panic(err.Error())
		}
	}
	return vo.OK
}

func (cr CourseRepository) CreateCourse(course *model.Course) (code vo.ErrNo) {
	if err := cr.DB.Create(&course).Error; err != nil {
		panic(err.Error())
	}
	return vo.OK
}
