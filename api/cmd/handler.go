package main

import (
	"fmt"
	"github.com/yuhaoyuan/Http_server/api/auth"
	"github.com/yuhaoyuan/Http_server/rpc"
	"github.com/yuhaoyuan/Http_server/util"
	"github.com/yuhaoyuan/RPC_server/dal"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func HandHome(w http.ResponseWriter, r *http.Request) {
	token, _ := auth.FromRequest(r)
	if token == "" {    // 没有token
		_, _ = fmt.Fprintf(w, "%s", HtmlInfoMp["home"])
	}
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
	userName :=	r.Form["user_name"][0]
	passwd := r.Form["user_pwd"][0]
	token := r.Form["token"][0]
	log.Println("in HandLogin--------------------userName=", userName)
	/*

	*/
	rpc.Mut.Lock()

	log.Println("in HandLogin--------------------getrpc-client-lock---userName=", userName)
	defer rpc.Mut.Unlock()
	rpcClient := rpc.GetSingleton()
	defer rpcClient.Close()

	log.Println("in HandLogin--------------------getrpc-client-done----userName=", userName)
	if token != "" {  // 如果有token,校验token
		//var checkTokenRequest = RpcProto.CheckTokenRequest  !!!!!!!!!这么用的话 = 并发的请求共用一个request = gg!!!!!!!
		var checkTokenRequest func(userName, token string) (dal.UserInfo, error)
		newReq := checkTokenRequest
		rpcClient.Call("CheckToken", &newReq)
		tokenInfo, _ := newReq(userName, token)   // 发送请求 ---------------------- 瓶颈之一， 有一个请求拿到锁之后卡在这里的话，后面的都会gg----------
		if tokenInfo == (dal.UserInfo{}){
			_, _ = fmt.Fprintf(w, "%s", fmt.Sprintf(string(HtmlInfoMp["home"]), "token 过期，请重新登陆"))
			return
		}

		// debug 一下
		log.Printf("--------check bug--------- req-name=%s, rsp-name=%s", userName, tokenInfo.Name)
		// 返回数据给h5
		_, _ = fmt.Fprintf(w, "%s", fmt.Sprintf(string(HtmlInfoMp["login_success"]), tokenInfo.Name, tokenInfo.Token, tokenInfo.NickName, tokenInfo.Picture))
	} else {   // 如果没有，校验密码、获得token
		// 调用 RPC server
		//var loginRequest= RpcProto.LoginRequest
		var loginRequest func(string, string) (dal.UserInfo, error)
			rpcClient.Call("userLogin", &loginRequest)
		rsp, err := loginRequest(userName, passwd) // 发送请求
		if err != nil {
			log.Println("HandLogin - loginRequest error = ", err)
		}
		// 返回数据给h5
		_, _ = fmt.Fprintf(w, "%s", fmt.Sprintf(string(HtmlInfoMp["login_success"]), rsp.Name, rsp.Token, rsp.NickName, rsp.Picture))
	}
	log.Println("handle Login done! ")
}

func HandRegister(w http.ResponseWriter, r *http.Request) {
	ret, _ := fmt.Fprintf(w, "%s", HtmlInfoMp["register"])
	fmt.Println("call register ", ret)
}

func HandRegisterUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "the method is not allowed！", http.StatusMethodNotAllowed)
		return
	}
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user :=	r.Form["user_name"][0]
	passwd := r.Form["user_pwd"][0]

	// 调用 RPC server
	//rpcClient:= rpc.RpcInit() // 记得在用完之后关闭连接
	rpcClient := rpc.GetSingleton() // 并发起来，为什么会直接跳过呢？

	//var registerRequest = RpcProto.RegisterRequest
	var registerRequest func(userName string, pwd string) (dal.UserInfo, error)
	rpcClient.Call("userRegister", &registerRequest)
	rsp, err := registerRequest(user, passwd)  // 发送请求
	if err != nil{
		log.Println("HandRegisterUpload - registerRequest error = ", err)
	}

	// 返回数据给h5
	_, _ = fmt.Fprintf(w, "%s", fmt.Sprintf(string(HtmlInfoMp["login_success"]), rsp.Name, rsp.Token, rsp.NickName, rsp.Picture))
	//rpcClient.Close()
}

func HandModify(w http.ResponseWriter, r *http.Request){
	if r.Method != "POST" {
		http.Error(w, "the method is not allowed！", http.StatusMethodNotAllowed)
		return
	}
	//token, _ := auth.FromRequest(r)

	// ------------获得用户修改的资料信息--------------
	err := r.ParseForm()
	imgFile, _, imgErr := r.FormFile("imgfile")
	if imgErr != nil && imgErr != http.ErrMissingFile{
		log.Println("HandModify - FormFile(imgfile) error = ", imgErr.Error())
		return
	}
	userName :=	r.Form["user_name"][0]
	token := r.Form["token"][0]
	nickName := r.Form["nick_name"][0]

	pictureCdnUrl := ""  // 为空则代表没有更换头像
	if imgFile!= nil {
		// ------------上传图片至cdn------------
		fileName := fmt.Sprintf("%s_%d.jpg", userName, time.Now().UnixNano())
		filePath := fmt.Sprintf("/tmp/%s", fileName)
		out, err := os.Create(filePath)
		if err != nil {
			log.Println("HandModify -Failed to open the file for writing")
			return
		}
		_, err = io.Copy(out, imgFile)
		if err != nil {
			_, _ = fmt.Fprintln(w, err)
		}
		pictureCdnUrl = util.UploadQiniu(filePath, fileName)
		//删除图片
		out.Close()
		_ = os.Remove(filePath)
	}

	// ------------call rpc------------
	//rpcClient:= rpc.RpcInit() // 记得在用完之后关闭连接
	rpcClient := rpc.GetSingleton()

	var checkTokenRequest func(userName, token string) (dal.UserInfo, error)
	rpcClient.Call("CheckToken", &checkTokenRequest)
	tokenInfo, err := checkTokenRequest(userName, token)   // 发送请求
	if tokenInfo == (dal.UserInfo{}){
		_, _ = fmt.Fprintf(w, "%s", fmt.Sprintf(string(HtmlInfoMp["home"]), "token 过期，请重新登陆"))
		//rpcClient.Close()
		return
	}
	//var modifyRequest = RpcProto.ModifyInfoRequest
	var modifyRequest func(userName, pwd, nickName, picture string) (dal.UserInfo, error)
	rpcClient.Call("UserModifyInfo", &modifyRequest)
	_, err = modifyRequest(userName, tokenInfo.Pwd, nickName, pictureCdnUrl)   // 发送请求
	if err != nil{
		log.Println("HandModify -modifyRequest err = ", err)
		// 返回数据给h5
		_, _ = fmt.Fprintf(w, "%s", fmt.Sprintf(string(HtmlInfoMp["modify_error"]), err))
	} else{
		// 返回数据给h5
		_, _ = fmt.Fprintf(w, "%s", fmt.Sprintf(string(HtmlInfoMp["modify_success"]), userName, tokenInfo.Token))
	}
	log.Println("handle Login done! ")
}


