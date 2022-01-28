package model

type User struct {
	Uuid      int64  `gorm:"primaryKey; autoIncrement"`
	UserName  string `gorm:"varchar(20)"`
	NickName  string `gorm:"varchar(20)"`
	password  string `gorm:"varchar(20)"`
	RoleId    string `gorm:"varchar(1)"`
	Enabled   int8
}

func (User) TableName() string{
	return "user"
}
