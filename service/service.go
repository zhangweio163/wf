package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"time"
	"userservice/user"
)

var (
	port = flag.Int("port", 8080, "port")
	addr *string
	name *string
)

func InitConf() {
	addr = flag.String("reg", "127.0.0.1:8500", "reg")
	name = flag.String("name", "default", "name")
	flag.Parse()
}
func RunServer() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatal(err)
	}
	//先注册
	err = Reg(*port)
	if err != nil {
		log.Fatal("无法连接注册中心")
	}
	//注册完毕后每15S激活一下状态
	go func() {
		for {
			err = Reg(*port)
			if err != nil {
				log.Println("无法连接注册中心")
			} else {
				log.Println("健康检测,连接注册中心成功")
			}
			time.Sleep(15 * time.Second)
		}
	}()
	s := grpc.NewServer()
	user.RegisterUserServiceServer(s, &user.MyUser{})
	log.Printf("server started on port %d. Can be used Ctrl+C to stop server", *port)
	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
func Reg(port int) error {
	con, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println(err)
		return err
	}
	c := user.NewRegisterClient(con)
	reg := &user.RegisterReq{
		Address:      fmt.Sprintf("127.0.0.1:%d", port),
		Heathaddress: "",
		Token:        "",
		Servicename:  *name,
	}
	_, err = c.GoRegister(context.Background(), reg)
	if err != nil {
		return err
	}
	return nil
}
func main() {
	InitConf()
	RunServer()
}
