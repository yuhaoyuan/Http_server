package main

import (
	"encoding/gob"
	"fmt"
	"github.com/yuhaoyuan/Http_server/config"
	"github.com/yuhaoyuan/RPC_server/dal"
	"io/ioutil"
	"net"
	"net/http"
	"os"
)


func init() {
	nowPath, err := os.Getwd()
	if err != nil {
		panic("get nowPath failed!")
	}
	loadHtml("login_error", nowPath + "/api/cmd/template/login_error.html")
	loadHtml("login_success", nowPath + "/api/cmd/template/login_success.html")
	loadHtml("home", nowPath + "/api/cmd/template/home.html")
	loadHtml("modify_error", nowPath + "/api/cmd/template/modify_error.html")
	loadHtml("modify_success", nowPath + "/api/cmd/template/modify_success.html")

	loadHtml("register", nowPath + "/api/cmd/template/register.html")
	//loadHtml("err", "/home/guaniu/code/src/http/err.html")
	//loadHtml("reg", "/home/guaniu/code/src/http/reg.html")
	//loadHtml("errtwo", "/home/guaniu/code/src/http/errtwo.html")


	config.BaseConfInit()
	RpcInit()
}

var (
	HtmlInfoMp = make(map[string][]byte)
	userMap    = make(map[string]string)
)

func loadHtml(key, fileName string) {
	info, err := readFile(fileName)
	if err != nil {
		fmt.Print(err)
		return
	}
	HtmlInfoMp[key] = info
}

func readFile(fileName string) ([]byte, error) {
	f, err := os.Open(fileName)
	if f != nil {
		defer f.Close()
	}
	if err != nil {
		panic(err)
	}
	return ioutil.ReadAll(f)
}

func main(){
	gob.Register(dal.UserInfo{})

	router := &Router{}                    // todo 想一下这个如果不是指针呢?
	ln, err := net.Listen("tcp", ":8001") // todo：想一下这里支持的最大并发数是多少
	if err != nil{
		fmt.Println("Listen failed！")
	}
	err2 := http.Serve(ln, router)
	if err2 != nil{
		fmt.Printf("Serve failed")
	}
}