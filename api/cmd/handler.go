package main

import (
	"fmt"
	"github.com/yuhaoyuan/Http_server/rpc"
	"github.com/yuhaoyuan/Http_server/util"
	"github.com/yuhaoyuan/RPC_server/dal"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// HandHome home
func HandHome(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "%s", HTMLInfoMp["home"])
}

// HandLogin login
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
	userName := r.Form["user_name"][0]
	passwd := r.Form["user_pwd"][0]
	token := r.Form["token"][0]

	rpc.Mut.Lock()
	defer rpc.Mut.Unlock()
	rpcClient := rpc.GetSingleton()
	//rpcClient := rpc.GetNewClient()

	//log.Println("in HandLogin-getrpc-client-done----userName=", userName)
	if token != "" { // 如果有token,校验token
		var checkTokenRequest func(userName, token string) (dal.UserInfo, error)
		newReq := checkTokenRequest
		rpcClient.Call("CheckToken", &newReq)
		tokenInfo, err := newReq(userName, token) // 发送请求
		if err != nil {
			log.Println("newReq - err = ", err)
			_, _ = fmt.Fprintf(w, "%s", fmt.Sprintf(string(HTMLInfoMp["home"]), err))
			return
		}
		if tokenInfo == (dal.UserInfo{}) {
			_, _ = fmt.Fprintf(w, "%s", fmt.Sprintf(string(HTMLInfoMp["home"]), "token 过期，请重新登陆"))
			return
		}

		//log.Printf("-------- req-name=%s, rsp-name=%s", userName, tokenInfo.Name)
		// 返回数据给h5
		_, _ = fmt.Fprintf(w, "%s", fmt.Sprintf(string(HTMLInfoMp["login_success"]), tokenInfo.Name, tokenInfo.Token, tokenInfo.NickName, tokenInfo.Picture))
	} else { // 如果没有，校验密码、获得token
		// 调用 RPC server
		var loginRequest func(string, string) (dal.UserInfo, error)
		rpcClient.Call("userLogin", &loginRequest)
		rsp, err := loginRequest(userName, passwd) // 发送请求
		if err != nil {
			log.Println("HandLogin - loginRequest error = ", err)
		}
		// 返回数据给h5
		_, _ = fmt.Fprintf(w, "%s", fmt.Sprintf(string(HTMLInfoMp["login_success"]), rsp.Name, rsp.Token, rsp.NickName, rsp.Picture))
	}
}

// HandRegister register
func HandRegister(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "%s", HTMLInfoMp["register"])
	//fmt.Println("call register ", ret)
}

// HandRegisterUpload 提交注册
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
	user := r.Form["user_name"][0]
	passwd := r.Form["user_pwd"][0]

	// 调用 RPC server
	rpc.Mut.Lock()
	defer rpc.Mut.Unlock()
	rpcClient := rpc.GetSingleton() //

	//var registerRequest = RpcProto.RegisterRequest
	var registerRequest func(userName string, pwd string) (dal.UserInfo, error)
	rpcClient.Call("userRegister", &registerRequest)
	rsp, err := registerRequest(user, passwd) // 发送请求
	if err != nil {
		log.Println("HandRegisterUpload - registerRequest error = ", err)
	}

	// 返回数据给h5
	_, _ = fmt.Fprintf(w, "%s", fmt.Sprintf(string(HTMLInfoMp["login_success"]), rsp.Name, rsp.Token, rsp.NickName, rsp.Picture))
	//rpcClient.Close()
}

// HandModify 修改资料
func HandModify(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "the method is not allowed！", http.StatusMethodNotAllowed)
		return
	}
	//token, _ := auth.FromRequest(r)

	// ------------获得用户修改的资料信息--------------
	err := r.ParseForm()
	imgFile, _, imgErr := r.FormFile("imgfile")
	if imgErr != nil && imgErr != http.ErrMissingFile {
		log.Println("HandModify - FormFile(imgfile) error = ", imgErr.Error())
		return
	}
	userName := r.Form["user_name"][0]
	token := r.Form["token"][0]
	nickName := r.Form["nick_name"][0]

	pictureCDNURl := "" // 为空则代表没有更换头像
	if imgFile != nil {
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
		pictureCDNURl = util.UploadQiniu(filePath, fileName)
		//删除图片
		out.Close()
		_ = os.Remove(filePath)
	}

	// ------------call rpc-----------
	rpc.Mut.Lock()
	defer rpc.Mut.Unlock()
	rpcClient := rpc.GetSingleton()

	var checkTokenRequest func(userName, token string) (dal.UserInfo, error)
	rpcClient.Call("CheckToken", &checkTokenRequest)
	tokenInfo, err := checkTokenRequest(userName, token) // 发送请求
	if tokenInfo == (dal.UserInfo{}) {
		_, _ = fmt.Fprintf(w, "%s", fmt.Sprintf(string(HTMLInfoMp["home"]), "token 过期，请重新登陆"))
		return
	}
	var modifyRequest func(userName, pwd, nickName, picture string) (dal.UserInfo, error)
	rpcClient.Call("UserModifyInfo", &modifyRequest)
	_, err = modifyRequest(userName, tokenInfo.Pwd, nickName, pictureCDNURl) // 发送请求
	if err != nil {
		log.Println("HandModify -modifyRequest err = ", err)
		// 返回数据给h5
		_, _ = fmt.Fprintf(w, "%s", fmt.Sprintf(string(HTMLInfoMp["modify_error"]), err))
	} else {
		// 返回数据给h5
		_, _ = fmt.Fprintf(w, "%s", fmt.Sprintf(string(HTMLInfoMp["modify_success"]), userName, tokenInfo.Token))
	}
	//log.Println("handle Login done! ")
}
