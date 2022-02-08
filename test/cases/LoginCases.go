package cases

import (
	"course_system/test"
	"course_system/vo"
	"fmt"
	"net/http"
	"strconv"
)

var LoginCases = []test.LoginTest{
	{ //成功登录
		Req: vo.LoginRequest{
			Username: "JudgeAdmin",
			Password: "JudgePassword2022",
		},
		ExpCode: http.StatusOK,
		ExpResp: vo.LoginResponse{
			Code: vo.OK,
			Data: struct{ UserID string }{UserID: "1"},
		},
	},
	{ //参数不合法, 密码长度过短
		Req: vo.LoginRequest{
			Username: "Alexander",
			Password: "alexand",
		},
		ExpCode: http.StatusOK,
		ExpResp: vo.LoginResponse{
			Code: vo.ParamInvalid,
			Data: struct{ UserID string }{UserID: "0"},
		},
	},
	{ //参数不合法, 密码长度过长
		Req: vo.LoginRequest{
			Username: "Alexander",
			Password: "alexandfgdsesdfdsfswaweddsfdd",
		},
		ExpCode: http.StatusOK,
		ExpResp: vo.LoginResponse{
			Code: vo.ParamInvalid,
			Data: struct{ UserID string }{UserID: "0"},
		},
	},
	{ //参数不合法, 用户名含数字
		Req: vo.LoginRequest{
			Username: "Alexander123",
			Password: "alexanderPass2022",
		},
		ExpCode: http.StatusOK,
		ExpResp: vo.LoginResponse{
			Code: vo.ParamInvalid,
			Data: struct{ UserID string }{UserID: "0"},
		},
	},
	{ //错误用户名
		Req: vo.LoginRequest{
			Username: "alexander",
			Password: "alexanderPass2022",
		},
		ExpCode: http.StatusOK,
		ExpResp: vo.LoginResponse{
			Code: vo.WrongPassword,
			Data: struct{ UserID string }{UserID: "0"},
		},
	},
	{ //已删除用户
		Req: vo.LoginRequest{
			Username: "JudgeAdminQ",
			Password: "JudgePassword2022Q",
		},
		ExpCode: http.StatusOK,
		ExpResp: vo.LoginResponse{
			Code: vo.UserHasDeleted,
			Data: struct{ UserID string }{UserID: "5"},
		},
	},
	{ //错误密码
		Req: vo.LoginRequest{
			Username: "Alexander",
			Password: "AlexanderPass2021",
		},
		ExpCode: http.StatusOK,
		ExpResp: vo.LoginResponse{
			Code: vo.WrongPassword,
			Data: struct{ UserID string }{UserID: "2"},
		},
	},
}

//仿照UserCases.go的写法, 调用其生成随机字符串函数
func GenerateLoginCase(i int) (tc test.LoginTest) {
	tc = test.LoginTest{
		Req: vo.LoginRequest{
			Username: RandString(10),
			Password: fmt.Sprintf("passworD%d", i),
		},
		ExpCode: http.StatusOK,
		ExpResp: vo.LoginResponse{
			Code: vo.WrongPassword,
			Data: struct {
				UserID string
			}{"0"},
		},
	}
	return tc
}

var LogoutCases = []test.LogoutTest{
	{
		Req:     "1",
		ExpCode: http.StatusOK,
		ExpResp: vo.LogoutResponse{
			Code: vo.OK,
		},
	},
}

func GenerateLogoutCase(i int) (tc test.LogoutTest) {
	tc = test.LogoutTest{
		Req:     strconv.FormatInt(int64(i), 10),
		ExpCode: http.StatusOK,
		ExpResp: vo.LogoutResponse{
			Code: vo.LoginRequired,
		},
	}
	return tc
}

var WhoAmICases = []test.WhoAmITest{
	{ //正常登录
		Req:     "1",
		ExpCode: http.StatusOK,
		ExpResp: vo.WhoAmIResponse{
			Code: vo.OK,
			Data: vo.TMember{
				UserID:   strconv.FormatInt(1, 10),
				Nickname: "JudgeAdmin",
				Username: "JudgeAdmin",
				UserType: 1,
			},
		},
	},
	{ //不存在
		Req:     "10",
		ExpCode: http.StatusOK,
		ExpResp: vo.WhoAmIResponse{
			Code: vo.UserNotExisted,
			Data: vo.TMember{
				UserID:   strconv.FormatInt(0, 10),
				Nickname: "",
				Username: "",
				UserType: 0,
			},
		},
	},
}

func GenerateWhoAmICase(i int) (tc test.WhoAmITest) {
	tc = test.WhoAmITest{
		Req:     strconv.FormatInt(int64(i), 10),
		ExpCode: http.StatusOK,
		ExpResp: vo.WhoAmIResponse{
			Code: vo.LoginRequired,
			Data: vo.TMember{
				UserID:   strconv.FormatInt(0, 10),
				Nickname: "",
				Username: "",
				UserType: 0,
			},
		},
	}
	return tc
}
