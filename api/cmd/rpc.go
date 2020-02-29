package main

import (
	"github.com/yuhaoyuan/Http_server/config"
	"github.com/yuhaoyuan/RPC_server/corn"
	"log"
	"net"
)

var rpcClient *corn.Client

func RpcInit(){
	// 调用rpc服务
	conn, err := net.Dial("tcp", config.BaseConf.Addr)
	if err != nil {
		log.Printf("client-dial failed!, err = ", err)
	}
	rpcClient = corn.NewClient(conn)

	// todo: 一直连接会有问题吗？
	//if conn != nil{
	//	_ = conn.Close()
	//}
}