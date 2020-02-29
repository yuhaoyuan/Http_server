package main

import (
	"fmt"
	"github.com/yuhaoyuan/Http_server/util"
	RpcProto "github.com/yuhaoyuan/RPC_server/proto"
	"io"
	"log"
	"net/http"
	"os"
	"time"
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
	user :=	r.Form["user_name"][0]
	passwd := r.Form["user_pwd"][0]

	//log.Println("user = ", user)
	//log.Println("passwd = ", passwd)

	// 调用 RPC server
	var loginRequest = RpcProto.LoginRequest
	rpcClient.Call("userLogin", &loginRequest)
	rsp, err := loginRequest(user, passwd)  // 发送请求
	if err != nil{
		log.Println(err)
	}
	log.Println(rsp)

	// 返回数据给h5
	_, _ = fmt.Fprintf(w, "%s", fmt.Sprintf(string(HtmlInfoMp["login_success"]), rsp.Name, rsp.Pwd, rsp.NickName, rsp.Picture))

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
	var registerRequest = RpcProto.RegisterRequest
	rpcClient.Call("userLogin", &registerRequest)
	rsp, err := registerRequest(user, passwd)  // 发送请求
	if err != nil{
		log.Println(err)
	}
	log.Println(rsp)

	// 返回数据给h5
	_, _ = fmt.Fprintf(w, "%s", fmt.Sprintf(string(HtmlInfoMp["login_success"]), rsp.Name, rsp.Pwd, rsp.NickName, rsp.Picture))
}

func HandModify(w http.ResponseWriter, r *http.Request){
	if r.Method != "POST" {
		http.Error(w, "the method is not allowed！", http.StatusMethodNotAllowed)
		return
	}
	// ------------获得用户修改的资料信息--------------
	err := r.ParseForm()
	imgFile, _, imgErr := r.FormFile("imgfile")
	if imgErr != nil && imgErr != http.ErrMissingFile{
		log.Println(imgErr.Error())
		return
	}
	userName :=	r.Form["user_name"][0]
	passwd := r.Form["user_pwd"][0]
	nickName := r.Form["nick_name"][0]

	pictureCdnUrl := ""  // 为空则代表没有更换头像
	if imgFile!= nil {
		// ------------上传图片至cdn------------
		fileName := fmt.Sprintf("%s_%d.jpg", userName, time.Now().UnixNano())
		filePath := fmt.Sprintf("/tmp/%s", fileName)
		out, err := os.Create(filePath)
		if err != nil {
			log.Println(err)
			_, _ = fmt.Fprintf(w, "Failed to open the file for writing")
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
	var modifyRequest = RpcProto.ModifyInfoRequest
	rpcClient.Call("UserModifyInfo", &modifyRequest)
	_, err = modifyRequest(userName, passwd, nickName, pictureCdnUrl)   // 发送请求
	if err != nil{
		log.Println(err)
		// 返回数据给h5
		_, _ = fmt.Fprintf(w, "%s", fmt.Sprintf(string(HtmlInfoMp["modify_error"]), err))
	} else{
		// 返回数据给h5
		_, _ = fmt.Fprintf(w, "%s", fmt.Sprintf(string(HtmlInfoMp["modify_success"]), userName, passwd))
	}
	log.Println("handle Login done! ")
}


