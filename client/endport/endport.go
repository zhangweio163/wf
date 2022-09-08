package endport

import (
	"Userclient/user"
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var serverList []user.UserServiceClient
var c user.RegisterClient
var addr *string
var name *string

func InitConf() {
	addr = flag.String("reg", "127.0.0.1:8500", "reg")
	name = flag.String("name", "default", "name")
	flag.Parse()
}
func ConnectGrpc() {
	serverList = []user.UserServiceClient{}
	//获取服务
	AddrList := GetServerlist()
	for _, addr := range AddrList {
		con, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Println(err)
		}
		uc := user.NewUserServiceClient(con)
		serverList = append(serverList, uc)
	}
}

func GetServerlist() []string {
	var addrlist []string
	con, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println(err)
	}
	c = user.NewRegisterClient(con)
	res, err := c.GetRegister(context.Background(), &user.GetRegisterReq{Servicename: *name})
	if res == nil {
		return addrlist
	}
	for i := 0; i < len(res.Serverlist); i++ {
		addrlist = append(addrlist, res.Serverlist[i].Address)
	}
	return addrlist
}

func RandomSelectService() user.UserServiceClient {
	ConnectGrpc()
	rand.Seed(time.Now().UnixNano())
	if len(serverList) == 0 {
		return nil
	}
	n := rand.Intn(len(serverList))
	return serverList[n]
}

func GetUserEndPort(req *user.UserReq) *user.UsersRes {
	uc := RandomSelectService()
	if uc == nil {
		return &user.UsersRes{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("请求失败:当前没有可用节点"),
			Data:    []*user.User{}}
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := uc.GetUser(ctx, req)
	if err != nil {
		log.Println(err)
		return &user.UsersRes{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("请求失败:%s", err),
			Data:    []*user.User{}}
	}
	return r
}
