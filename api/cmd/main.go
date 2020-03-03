package main

import (
	"encoding/gob"
	"fmt"
	"github.com/yuhaoyuan/Http_server/config"
	"github.com/yuhaoyuan/Http_server/yhylog"
	"github.com/yuhaoyuan/RPC_server/dal"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
)

func init() {
	//nowPath, err := os.Getwd()
	//if err != nil {
	//	panic("get nowPath failed!")
	//}
	loadHtml("login_error", "/Users/yuhaoyuan/work/Http_server/template/login_error.html")
	loadHtml("login_success", "/Users/yuhaoyuan/work/Http_server/template/login_success.html")
	loadHtml("home", "/Users/yuhaoyuan/work/Http_server/template/home.html")
	loadHtml("modify_error", "/Users/yuhaoyuan/work/Http_server/template/modify_error.html")
	loadHtml("modify_success", "/Users/yuhaoyuan/work/Http_server/template/modify_success.html")

	loadHtml("register", "/Users/yuhaoyuan/work/Http_server/template/register.html")

	config.BaseConfInit()
	yhylog.LogInit(config.BaseConf.LogName)
	//rpc.InitRpc()
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

func main() {
	gob.Register(dal.UserInfo{})

	router := &Router{}                   // todo 想一下这个如果不是指针呢?
	ln, err := net.Listen("tcp", ":8001") // todo：想一下这里支持的最大并发数是多少
	if err != nil {
		log.Println("Listen failed！")
	}
	err2 := http.Serve(ln, router)
	if err2 != nil {
		log.Printf("Serve failed")
	}
}
