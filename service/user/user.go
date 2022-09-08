package user

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
)

type MyUser struct {
}

var (
	User1 = User{
		Userid: 1,
		Name:   "张三",
		Email:  "1030268572@qq.com",
		Phone:  "13260618533",
	}
	User2 = User{
		Userid: 2,
		Name:   "李四",
		Email:  "1030268572@qq.com",
		Phone:  "13260618533",
	}
	Users = []*User{&User1, &User2}
)

func (m MyUser) GetUser(ctx context.Context, req *UserReq) (*UsersRes, error) {
	log.Println("收到查询用户请求")
	if req.Userid == 0 {
		return &UsersRes{
			Code:    http.StatusOK,
			Message: "请求成功",
			Data:    Users,
		}, nil
	}
	u, err := FindUser(req.Userid, Users)
	if err != nil {
		return &UsersRes{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("请求失败%s", err),
			Data:    []*User{u},
		}, err
	}
	return &UsersRes{
		Code:    http.StatusOK,
		Message: "请求成功",
		Data:    []*User{u},
	}, nil
}

func FindUser(id int32, data []*User) (*User, error) {
	for i := 0; i < len(data); i++ {
		if data[i].Userid == id {
			return data[i], nil
		}
	}
	return nil, errors.New("not found")
}
