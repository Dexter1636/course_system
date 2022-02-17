package vo

import "strconv"

func ValidateCreateCourseReq(req CreateCourseRequest) (code ErrNo) {
	code = ValidateCourseCap(req.Cap)
	return code
}

func ValidateGetCourseReq(req GetCourseRequest) (code ErrNo) {
	return ValidateStrCourseId(req.CourseID)
}

func ValidateBookCourseReq(req BookCourseRequest) (code ErrNo) {
	code = ValidateStrStudentId(req.StudentID)
	if code == OK {
		code = ValidateStrCourseId(req.CourseID)
	}
	return code
}

func ValidateGetStudentCourseReq(req GetStudentCourseRequest) (code ErrNo) {
	return ValidateStrStudentId(req.StudentID)
}

func ValidateStrStudentId(idStr string) (code ErrNo) {
	stuId, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return ParamInvalid
	}
	return ValidateStudentId(stuId)
}

func ValidateStudentId(id int64) (code ErrNo) {
	if id > 0 {
		return OK
	} else {
		return ParamInvalid
	}
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
