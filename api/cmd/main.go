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
	loadHTML("login_error", "/Users/yuhaoyuan/work/Http_server/template/login_error.html")
	loadHTML("login_success", "/Users/yuhaoyuan/work/Http_server/template/login_success.html")
	loadHTML("home", "/Users/yuhaoyuan/work/Http_server/template/home.html")
	loadHTML("modify_error", "/Users/yuhaoyuan/work/Http_server/template/modify_error.html")
	loadHTML("modify_success", "/Users/yuhaoyuan/work/Http_server/template/modify_success.html")

	loadHTML("register", "/Users/yuhaoyuan/work/Http_server/template/register.html")

	config.BaseConfInit()
	yhylog.LogInit(config.BaseConf.LogName)
	//rpc.SpecialRPClientInit()
}

var (
	// HTMLInfoMp htmlmap
	HTMLInfoMp = make(map[string][]byte)
)

// loadHtml .
func loadHTML(key, fileName string) {
	info, err := readFile(fileName)
	if err != nil {
		fmt.Print(err)
		return
	}
	HTMLInfoMp[key] = info
}
// readFile .
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

	router := &Router{}
	ln, err := net.Listen("tcp", ":8001") // todo：想一下这里支持的最大并发数是多少
	if err != nil {
		log.Println("Listen failed！")
	}
	err2 := http.Serve(ln, router)
	if err2 != nil {
		log.Printf("Serve failed")
	}
}
