package model

type Course struct {
	Id        int64  `gorm:"primaryKey; autoIncrement"`
	Name      string `gorm:"varchar(40)"`
	Cap       int
	Avail     int
	TeacherId int64
}

func (Course) TableName() string {
	return "course"
}
