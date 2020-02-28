package main

import (
	"fmt"
	"github.com/yuhaoyuan/Http_server/config"
	"github.com/yuhaoyuan/RPC_server/corn"
	"github.com/yuhaoyuan/RPC_server/proto"
	"log"
	"net"
	"net/http"
)

func HandHome(w http.ResponseWriter) {
	ret, _ := fmt.Fprintf(w, "%s", HtmlInfoMp["home"])
	fmt.Println("call home ", ret)
}

func HandLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "the method is not allowed！", http.StatusMethodNotAllowed)
		return
	}
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// 设置Content-Type
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	user :=	r.Form["user_name"][0]
	passwd := r.Form["user_pwd"][0]

	// 调用rpc服务
	fmt.Println(user, passwd)
	conn, err := net.Dial("tcp", config.BaseConf.Addr)
	if err != nil {
		log.Printf("client-dial failed!")
	}
	cli := corn.NewClient(conn)

	var loginRequest func(string, string) (proto.User, error)
	cli.Call("userLogin", &loginRequest)
	rsp, err := loginRequest("u0001", "12345u")  // 发送请求
	if err != nil{
		log.Println(err)
	} else {
		log.Println(rsp)
	}
	_ = conn.Close()

}

func HandRegister(w http.ResponseWriter, r *http.Request) {

}

func HandModify(w http.ResponseWriter, r *http.Request){

}