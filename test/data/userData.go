package data

import (
	"course_system/common"
)

func InitAdmin() {
	common.GetDB().Exec("INSERT INTO user(uuid, user_name, nick_name, password, role_id, enabled) " +
		"VALUES (1, 'JudgeAdmin', 'JudgeAdmin', 'JudgePassword2022', '1', 1)")
}

func InitDataForUser() {
	//insert students
	InitAdmin()
	common.GetDB().Exec("INSERT INTO user(uuid, user_name, nick_name, password, role_id, enabled) " +
		"VALUES (0, 'Alexander', 'Alex', 'alexanderPass2022', '3', 1)")
	common.GetDB().Exec("INSERT INTO user(uuid, user_name, nick_name, password, role_id, enabled) " +
		"VALUES (0, 'Benjamin', 'Benj', 'BenjaminPass2022', '3', 1)")
	common.GetDB().Exec("INSERT INTO user(uuid, user_name, nick_name, password, role_id, enabled) " +
		"VALUES (0, 'RockyViavia', 'Rocky', 'RockyPass2022', '3', 1)")
	common.InitRedisData()
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
	common.InitRedisData()
}
