package cases

import (
	"course_system/test"
	"course_system/vo"
	"math/rand"
	"net/http"
	"strconv"
)

//二选一

func Choose2To1(i int, a, b int) int {
	return []int{a, b}[i]
}

//四选一

func Choose4To1(i int, a, b, c, d int) int {
	return []int{a, b, c, d}[i]
}

var GetMemberListCases = []test.GetMemberListTest{
	{
		Req: vo.GetMemberListRequest{
			Offset: 1,
			Limit:  0,
		},
		ExpCode: http.StatusOK,
		ExpResp: vo.GetMemberListResponse{
			Code: vo.ParamInvalid,
			Data: struct{ MemberList []vo.TMember }{MemberList: []vo.TMember{}},
		},
	},
	{
		Req: vo.GetMemberListRequest{
			Offset: 5,
			Limit:  1,
		},
		ExpCode: http.StatusOK,
		ExpResp: vo.GetMemberListResponse{
			Code: vo.OK,
			Data: struct{ MemberList []vo.TMember }{MemberList: []vo.TMember{
				{
					UserID:   "6",
					Nickname: "AlexQ",
					Username: "AlexanderQ",
					UserType: 3,
				},
			}},
		},
	},
	{
		Req: vo.GetMemberListRequest{
			Offset: 3,
			Limit:  2,
		},
		ExpCode: http.StatusOK,
		ExpResp: vo.GetMemberListResponse{
			Code: vo.OK,
			Data: struct{ MemberList []vo.TMember }{MemberList: []vo.TMember{
				{
					UserID:   "4",
					Nickname: "Rocky",
					Username: "RockyViavia",
					UserType: 3,
				},
				{
					UserID:   "5",
					Nickname: "JudgeAdminQ",
					Username: "JudgeAdminQ",
					UserType: 3,
				},
			}},
		},
	},
	{
		Req: vo.GetMemberListRequest{
			Offset: 7,
			Limit:  100,
		},
		ExpCode: http.StatusOK,
		ExpResp: vo.GetMemberListResponse{
			Code: vo.OK,
			Data: struct{ MemberList []vo.TMember }{MemberList: []vo.TMember{
				{
					UserID:   "8",
					Nickname: "RockyQ",
					Username: "RockyViaviaQ",
					UserType: 1,
				},
			}},
		},
	},
	{
		Req: vo.GetMemberListRequest{
			Offset: 0,
			Limit:  1,
		},
		ExpCode: http.StatusOK,
		ExpResp: vo.GetMemberListResponse{
			Code: vo.OK,
			Data: struct{ MemberList []vo.TMember }{MemberList: []vo.TMember{
				{
					UserID:   "1",
					Nickname: "JudgeAdmin",
					Username: "JudgeAdmin",
					UserType: 1,
				},
			}},
		},
	},
	{
		Req: vo.GetMemberListRequest{
			Offset: 0,
			Limit:  0,
		},
		ExpCode: http.StatusOK,
		ExpResp: vo.GetMemberListResponse{
			Code: vo.OK,
			Data: struct{ MemberList []vo.TMember }{MemberList: []vo.TMember{
				{
					UserID:   "1",
					Nickname: "JudgeAdmin",
					Username: "JudgeAdmin",
					UserType: 1,
				},
				{
					UserID:   "2",
					Nickname: "Alex",
					Username: "Alexander",
					UserType: 3,
				},
				{
					UserID:   "3",
					Nickname: "Benj",
					Username: "Benjamin",
					UserType: 3,
				},
				{
					UserID:   "4",
					Nickname: "Rocky",
					Username: "RockyViavia",
					UserType: 3,
				},
				{
					UserID:   "5",
					Nickname: "JudgeAdminQ",
					Username: "JudgeAdminQ",
					UserType: 3,
				},
				{
					UserID:   "6",
					Nickname: "AlexQ",
					Username: "AlexanderQ",
					UserType: 3,
				},
				{
					UserID:   "7",
					Nickname: "BenjQ",
					Username: "BenjaminQ",
					UserType: 2,
				},
				{
					UserID:   "8",
					Nickname: "RockyQ",
					Username: "RockyViaviaQ",
					UserType: 1,
				},
			}},
		},
	},
}

func GenerateGetMemberListCase(i int) (tc test.GetMemberListTest) {
	j := rand.Intn(2)
	tc = test.GetMemberListTest{
		Req: vo.GetMemberListRequest{
			Offset: rand.Intn(10000) + 10,
			Limit:  Choose2To1(j, 0, rand.Intn(10000)+1),
		},
		ExpCode: http.StatusOK,
		ExpResp: vo.GetMemberListResponse{
			Code: vo.ErrNo(Choose2To1(j, int(vo.ParamInvalid), int(vo.OK))),
			Data: struct{ MemberList []vo.TMember }{MemberList: []vo.TMember{}},
		},
	}
	return tc
}

var UpdateMemberCases = []test.UpdateMemberTest{
	{
		Req: vo.UpdateMemberRequest{
			UserID:   strconv.FormatInt(1, 10),
			Nickname: RandString(18),
		},
		ExpCode: http.StatusOK,
		ExpResp: vo.UpdateMemberResponse{
			Code: vo.OK,
		},
	},
	{
		Req: vo.UpdateMemberRequest{
			UserID:   strconv.FormatInt(5, 10),
			Nickname: RandString(18),
		},
		ExpCode: http.StatusOK,
		ExpResp: vo.UpdateMemberResponse{
			Code: vo.UserHasDeleted,
		},
	},
	{
		Req: vo.UpdateMemberRequest{
			UserID:   strconv.FormatInt(1, 10),
			Nickname: RandString(30),
		},
		ExpCode: http.StatusOK,
		ExpResp: vo.UpdateMemberResponse{
			Code: vo.ParamInvalid,
		},
	},
	{
		Req: vo.UpdateMemberRequest{
			UserID:   strconv.FormatInt(100, 10),
			Nickname: RandString(18),
		},
		ExpCode: http.StatusOK,
		ExpResp: vo.UpdateMemberResponse{
			Code: vo.UserNotExisted,
		},
	},
}

func GenerateUpdateMemberCase(i int) (tc test.UpdateMemberTest) {
	j := rand.Intn(4)
	tc = test.UpdateMemberTest{
		Req: vo.UpdateMemberRequest{
			UserID:   strconv.FormatInt(int64(Choose4To1(j, 1, 5, 1, rand.Intn(1000)+10)), 10),
			Nickname: RandString(18) + strconv.FormatInt(int64(Choose4To1(j, 0, 0, 1000000, 0)), 10),
		},
		ExpCode: http.StatusOK,
		ExpResp: vo.UpdateMemberResponse{
			Code: vo.ErrNo(Choose4To1(j, int(vo.OK), int(vo.UserHasDeleted), int(vo.ParamInvalid), int(vo.UserNotExisted))),
		},
	}
	return tc
}

var DeleteMemberCases = []test.DeleteMemberTest{
	{
		Req: vo.DeleteMemberRequest{
			UserID: strconv.FormatInt(1, 10),
		},
		ExpCode: http.StatusOK,
		ExpResp: vo.DeleteMemberResponse{
			Code: vo.OK,
		},
	},
	{
		Req: vo.DeleteMemberRequest{
			UserID: strconv.FormatInt(1, 10),
		},
		ExpCode: http.StatusOK,
		ExpResp: vo.DeleteMemberResponse{
			Code: vo.UserHasDeleted,
		},
	},
	{
		Req: vo.DeleteMemberRequest{
			UserID: strconv.FormatInt(5, 10),
		},
		ExpCode: http.StatusOK,
		ExpResp: vo.DeleteMemberResponse{
			Code: vo.UserHasDeleted,
		},
	},
	{
		Req: vo.DeleteMemberRequest{
			UserID: strconv.FormatInt(100, 10),
		},
		ExpCode: http.StatusOK,
		ExpResp: vo.DeleteMemberResponse{
			Code: vo.UserNotExisted,
		},
	},
}

func GenerateDeleteMemberCase(i int) (tc test.DeleteMemberTest) {
	j := rand.Intn(2)
	tc = test.DeleteMemberTest{
		Req: vo.DeleteMemberRequest{
			UserID: strconv.FormatInt(int64(Choose2To1(j, 5, rand.Intn(1000)+10)), 10),
		},
		ExpCode: http.StatusOK,
		ExpResp: vo.DeleteMemberResponse{
			Code: vo.ErrNo(Choose2To1(j, int(vo.UserHasDeleted), int(vo.UserNotExisted))),
		},
	}
	return tc
}
