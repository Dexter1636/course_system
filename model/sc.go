package model

type Sc struct {
	StudentId int64
	CourseId  int64
}

func (Sc) TableName() string {
	return "sc"
}
