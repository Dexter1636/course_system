package cases

import (
	"course_system/test"
	"course_system/vo"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

var CreateMemberCases = []test.CreateMemberTest{
	{
		Req: vo.CreateMemberRequest{
			Nickname: "JudgeAdmin",
			Username: "JudgeAdmin",
			Password: "JudgePassword2022",
			UserType: 1,
		},
		ExpCode: http.StatusOK,
		ExpResp: vo.CreateMemberResponse{
			Code: 0,
			Data: struct{ UserID string }{UserID: "1"},
		},
	},
	{
		Req: vo.CreateMemberRequest{
			Nickname: "Alex",
			Username: "Alexander",
			Password: "alexanderPass2022",
			UserType: 3,
		},
		ExpCode: http.StatusOK,
		ExpResp: vo.CreateMemberResponse{
			Code: 0,
			Data: struct{ UserID string }{UserID: "2"},
		},
	},
	{
		Req: vo.CreateMemberRequest{
			Nickname: "Benj",
			Username: "Benjamin",
			Password: "BenjaminPass2022",
			UserType: 3,
		},
		ExpCode: http.StatusOK,
		ExpResp: vo.CreateMemberResponse{
			Code: 0,
			Data: struct{ UserID string }{UserID: "3"},
		},
	},
	{
		Req: vo.CreateMemberRequest{
			Nickname: "Rocky",
			Username: "RockyViavia",
			Password: "RockyPass2022",
			UserType: 3,
		},
		ExpCode: http.StatusOK,
		ExpResp: vo.CreateMemberResponse{
			Code: 0,
			Data: struct{ UserID string }{UserID: "4"},
		},
	},
}

//生成随机字符串

func RandString(len int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var b int
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		if r.Intn(2) == 1 {
			b = r.Intn(26) + 65
		} else {
			b = r.Intn(26) + 97
		}

		bytes[i] = byte(b)
	}
	return string(bytes)
}

func GenerateCreateMemberCase(i int) (tc test.CreateMemberTest) {
	tc = test.CreateMemberTest{
		Req: vo.CreateMemberRequest{
			Nickname: RandString(5),
			Username: RandString(10),
			Password: fmt.Sprintf("passworD%d", i),
			UserType: vo.UserType(rand.Int()%2 + 2),
		},
		ExpCode: http.StatusOK,
		ExpResp: vo.CreateMemberResponse{
			Code: 0,
			Data: struct {
				UserID string
			}{strconv.FormatInt(int64(i+1), 10)},
		},
	}
	return tc
}

var GetMemberCases = []test.GetMemberTest{
	{
		Req:     vo.GetMemberRequest{UserID: "1"},
		ExpCode: http.StatusOK,
		ExpResp: vo.GetMemberResponse{
			Code: 0,
			Data: vo.TMember{
				UserID:   "1",
				Nickname: "JudgeAdmin",
				Username: "JudgeAdmin",
				UserType: 1,
			},
		},
	},
	{
		Req:     vo.GetMemberRequest{UserID: "2"},
		ExpCode: http.StatusOK,
		ExpResp: vo.GetMemberResponse{
			Code: 0,
			Data: vo.TMember{
				UserID:   "2",
				Nickname: "Alex",
				Username: "Alexander",
				UserType: 3,
			},
		},
	},
	{
		Req:     vo.GetMemberRequest{UserID: "3"},
		ExpCode: http.StatusOK,
		ExpResp: vo.GetMemberResponse{
			Code: 0,
			Data: vo.TMember{
				UserID:   "3",
				Nickname: "Benj",
				Username: "Benjamin",
				UserType: 3,
			},
		},
	},
	{
		Req:     vo.GetMemberRequest{UserID: "4"},
		ExpCode: http.StatusOK,
		ExpResp: vo.GetMemberResponse{
			Code: 0,
			Data: vo.TMember{
				UserID:   "4",
				Nickname: "Rocky",
				Username: "RockyViavia",
				UserType: 3,
			},
		},
	},
}

func GenerateGetMemberCase(i int) (tc test.GetMemberTest) {
	tc = test.GetMemberTest{
		Req:     vo.GetMemberRequest{UserID: strconv.FormatInt(rand.Int63n(1000)+1000, 10)},
		ExpCode: http.StatusOK,
		ExpResp: vo.GetMemberResponse{
			Code: vo.UserNotExisted,
			Data: vo.TMember{},
		},
	}
	return tc
}
