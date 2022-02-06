package data

import "course_system/common"

func InitDataForCourseCommon() {
	// insert students
	common.GetDB().Exec("INSERT INTO user(user_name, nick_name, password, role_id, enabled) VALUES ('Amy Wong', 'Amy', '123456', '2', 1)")
	common.GetDB().Exec("INSERT INTO user(user_name, nick_name, password, role_id, enabled) VALUES ('Dexter Peng', 'Dexter', '123456', '2', 1)")
	common.GetDB().Exec("INSERT INTO user(user_name, nick_name, password, role_id, enabled) VALUES ('San Zhang', 'San', '123456', '2', 1)")
	common.GetDB().Exec("INSERT INTO user(user_name, nick_name, password, role_id, enabled) VALUES ('Si Li', 'Si', '123456', '2', 1)")

	// insert courses
	common.GetDB().Exec("INSERT INTO course(name, avail, cap) VALUES ('test1', 1, 1)")
	common.GetDB().Exec("INSERT INTO course(name, avail, cap) VALUES ('test2', 3, 3)")
	common.GetDB().Exec("INSERT INTO course(name, avail, cap) VALUES ('test3', 0, 100)")
	common.GetDB().Exec("INSERT INTO course(name, avail, cap) VALUES ('test4', 100, 100)")
}
