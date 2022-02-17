package vo

import "strconv"

func ValidateCreateCourseReq(req CreateCourseRequest) (code ErrNo) {
	code = ValidateCourseCap(req.Cap)
	return code
}

func ValidateGetCourseReq(req GetCourseRequest) (code ErrNo) {
	return ValidateStrCourseId(req.CourseID)
}

func ValidateStrCourseId(idStr string) (code ErrNo) {
	courseId, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return ParamInvalid
	}
	return ValidateCourseId(courseId)
}

func ValidateCourseId(id int64) (code ErrNo) {
	if id > 0 {
		return OK
	} else {
		return ParamInvalid
	}
}

func ValidateCourseCap(cap int) (code ErrNo) {
	if cap >= 0 {
		return OK
	} else {
		return ParamInvalid
	}
}
