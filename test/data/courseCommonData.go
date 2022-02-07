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
}

func InitDataForUser() {
	//insert students
	common.GetDB().Exec("INSERT INTO user(uuid, user_name, nick_name, password, role_id, enabled) " +
		"VALUES (1, 'JudgeAdmin', 'JudgeAdmin', 'JudgePassword2022', '1', 1)")
	common.GetDB().Exec("INSERT INTO user(uuid, user_name, nick_name, password, role_id, enabled) " +
		"VALUES (0, 'Alexander', 'Alex', 'alexanderPass2022', '3', 1)")
	common.GetDB().Exec("INSERT INTO user(uuid, user_name, nick_name, password, role_id, enabled) " +
		"VALUES (0, 'Benjamin', 'Benj', 'BenjaminPass2022', '3', 1)")
	common.GetDB().Exec("INSERT INTO user(uuid, user_name, nick_name, password, role_id, enabled) " +
		"VALUES (0, 'RockyViavia', 'Rocky', 'RockyPass2022', '3', 1)")
}

func InitDataForUserOther() {
	//add deleted data
	common.GetDB().Exec("INSERT INTO user(uuid, user_name, nick_name, password, role_id, enabled) " +
		"VALUES (0, 'JudgeAdminQ', 'JudgeAdminQ', 'JudgePassword2022Q', '3', 0)")
	common.GetDB().Exec("INSERT INTO user(uuid, user_name, nick_name, password, role_id, enabled) " +
		"VALUES (0, 'AlexanderQ', 'AlexQ', 'alexanderPass2022Q', '3', 0)")
	common.GetDB().Exec("INSERT INTO user(uuid, user_name, nick_name, password, role_id, enabled) " +
		"VALUES (0, 'BenjaminQ', 'BenjQ', 'BenjaminPass2022Q', '2', 0)")
	common.GetDB().Exec("INSERT INTO user(uuid, user_name, nick_name, password, role_id, enabled) " +
		"VALUES (0, 'RockyViaviaQ', 'RockyQ', 'RockyPass2022Q', '1', 0)")
}
