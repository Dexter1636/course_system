package data

import (
	"course_system/common"
)

func InitDataForCourseCommon() {
	// insert students
	// 插入的数据不符合user限制，在windows+mysql8测试中无法进行，已修改 @author：陈朝
	common.GetDB().Exec("INSERT INTO user(user_name, nick_name, password, role_id, enabled) VALUES ('AmyWong', 'Amyy', '123456Ab', '2', 1)")
	common.GetDB().Exec("INSERT INTO user(user_name, nick_name, password, role_id, enabled) VALUES ('DexterPeng', 'Dexter', '123456Ab', '2', 1)")
	common.GetDB().Exec("INSERT INTO user(user_name, nick_name, password, role_id, enabled) VALUES ('SanZhang', 'Sannn', '123456Ab', '2', 1)")
	common.GetDB().Exec("INSERT INTO user(user_name, nick_name, password, role_id, enabled) VALUES ('SiLilisi', 'Sili', '123456Ab', '2', 1)")

	// insert courses
	common.GetDB().Exec("INSERT INTO course(name, avail, cap) VALUES ('test1', 1, 1)")
	common.GetDB().Exec("INSERT INTO course(name, avail, cap) VALUES ('test2', 3, 3)")
	common.GetDB().Exec("INSERT INTO course(name, avail, cap) VALUES ('test3', 0, 100)")
	common.GetDB().Exec("INSERT INTO course(name, avail, cap) VALUES ('test4', 100, 100)")

	common.InitRedisData()
}

func InitDataForUnbing() {
	common.GetDB().Exec("INSERT INTO course(name, avail, cap,teacher_id) VALUES ('test1', 1, 1,810)")
	common.GetDB().Exec("INSERT INTO course(name, avail, cap,teacher_id) VALUES ('test2', 3, 3,893)")
	common.GetDB().Exec("INSERT INTO course(name, avail, cap,teacher_id) VALUES ('test3', 0, 100,810)")
	common.GetDB().Exec("INSERT INTO course(name, avail, cap,teacher_id) VALUES ('test4', 100, 100,893)")
	common.InitRedisData()
}
