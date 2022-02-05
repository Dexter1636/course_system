package dto

import (
	"course_system/model"
	"course_system/vo"
	"strconv"
)

func ToTCourse(course model.Course) vo.TCourse {
	if course.Id == 0 {
		return vo.TCourse{}
	}
	return vo.TCourse{
		CourseID:  strconv.FormatInt(course.Id, 10),
		Name:      course.Name,
		TeacherID: strconv.FormatInt(course.TeacherId, 10),
	}
}
