package main

import (
	"flag"
	"fmt"
	"google.golang.org/grpc/keepalive"
	"net"
	_ "os"
	_ "os/signal"
	_ "syscall"
	"time"

	_ "github.com/satori/go.uuid"
	_ "go.uber.org/zap"
	"google.golang.org/grpc"
	_ "google.golang.org/grpc/health"
	_ "google.golang.org/grpc/health/grpc_health_v1"

	//_ "srvs03/user_srv/global"
	"srvs03/user_srv/handler"
	//"srvs03/user_srv/initialize"
	"srvs03/user_srv/proto"
	//"srvs03/user_srv/utils"
	_ "github.com/hashicorp/consul/api"
)

func main() {
	IP := flag.String("ip", "127.0.0.1", "ip地址")
	Port := flag.Int("port", 50052, "端口号")

	flag.Parse()
	fmt.Println("ip:", *IP, "port:", *Port)

	if *IP == "" || *Port == 0 {
		fmt.Println("请输入正确的ip地址和端口号")
		return
	}
	//s := grpc.NewServer()
	s := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: 5 * time.Minute, // <--- This fixes it!
		}),
	)
	proto.RegisterUserServer(s, &handler.UserServer{})
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *Port))
	if err != nil {
		fmt.Println("监听失败")
		return
	}
	err = s.Serve(lis)
	if err != nil {
		fmt.Println("服务器启动失败")
		return
	}
}
