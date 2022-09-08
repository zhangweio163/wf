package server

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

type ServerRegistration struct {
	ServerList map[string][]*Server `json:"servers"`
	mu         sync.Mutex
}
type MyRegister struct {
	ServerRegistration
}

var RegServer MyRegister

func (m *MyRegister) GoRegister(ctx context.Context, in *RegisterReq) (*RegisterRes, error) {
	var res *RegisterRes
	ms, err := m.InsertServer(in)
	if ms {
		res = &RegisterRes{
			Code:        fmt.Sprintf("%d", http.StatusOK),
			Address:     in.Address,
			Servicename: in.Servicename,
			Success:     true,
		}
		if err != nil {
			log.Println(err)
		} else {
			GoCheck(m, res)
		}
		return res, nil
	} else {
		res = &RegisterRes{Code: fmt.Sprintf("%d", http.StatusInternalServerError),
			Address:     in.Address,
			Servicename: in.Servicename,
			Success:     false,
		}
	}
	return res, fmt.Errorf("can't register'")
}
func (m *MyRegister) GetRegister(ctx context.Context, in *GetRegisterReq) (*GetRegisterRes, error) {
	serverlist := m.ServerList[in.Servicename]
	var res []*Server
	for i := 0; i < len(serverlist); i++ {
		if serverlist[i].Status {
			res = append(res, serverlist[i])
		}
	}
	if len(serverlist) == 0 {
		return nil, fmt.Errorf("no registered")
	}
	return &GetRegisterRes{
		Code:       http.StatusOK,
		Serverlist: res,
	}, nil
}
func (SReg *ServerRegistration) InitServer() {
	SReg.ServerList = make(map[string][]*Server)
}
func (s *ServerRegistration) InsertServer(req *RegisterReq) (bool, error) {
	_, err := s.HealthCheck(req.Address)
	if err != nil {
		log.Println(err)
		return false, nil
	}
	id, err := uuid.NewUUID()
	if err != nil {
		log.Println(err)
		return false, nil
	}
	uid := id.String()
	server := &Server{
		Name:    req.Servicename,
		Id:      uid,
		Address: req.Address,
		Status:  true,
	}
	if len(s.ServerList[req.Servicename]) == 0 {
		s.ServerList[req.Servicename] = append(s.ServerList[req.Servicename], server)
	} else {
		for i := 0; i < len(s.ServerList[req.Servicename]); i++ {
			if s.ServerList[req.Servicename][i].Address == req.Address {
				return true, fmt.Errorf("发现自主检测健康服务节点，地址:%s,服务名称：%s;", req.Address, req.Servicename)
			}
		}
		s.ServerList[req.Servicename] = append(s.ServerList[req.Servicename], server)
	}
	return true, nil
}
func (s *ServerRegistration) HealthCheck(addr string) (bool, error) {
	conn, err := net.DialTimeout("tcp", addr, 15*time.Second)
	if conn != nil {
		defer conn.Close()
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
func (s *ServerRegistration) SetHealthCheck(t bool, res *RegisterRes) {
	for i := 0; i < len(s.ServerList[res.Servicename]); i++ {
		if s.ServerList[res.Servicename][i].Address == res.Address {
			s.ServerList[res.Servicename][i].Status = t
		}
	}
}
func GoCheck(s *MyRegister, res *RegisterRes) {
	if res.Success == true {
		go func() {
			for {
				log.Printf("心跳检测开启,检测地址：%s\n", res.Address)
				r, err := s.HealthCheck(res.Address)
				if err != nil {
					log.Printf("心跳检测失败,检测地址：%s\n", res.Address)
					log.Println(err)
					s.SetHealthCheck(false, res)
				}
				if r == false {
					log.Printf("心跳检测失败,检测地址：%s\n", res.Address)
					s.SetHealthCheck(false, res)
				} else {
					log.Printf("心跳检测成功,响应地址：%s", res.Address)
					s.SetHealthCheck(true, res)
				}
				time.Sleep(15 * time.Second)
			}
		}()
	}
}
func NewRegServer() *MyRegister {
	return &RegServer
}
