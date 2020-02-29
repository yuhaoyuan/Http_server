package main

import (
	"fmt"
	"net/http"
)

type Router struct {
}

func (t *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	uPath := r.URL.Path
	fmt.Println("url call - ", uPath)
	switch uPath {
	case "/":
		t.Home(w, r)
	case "/login":
		t.Login(w, r)
	case "/register":
		t.Register(w, r)
	case "/register_upload":
		t.RegisterUpload(w, r)
	case "/modify":
		t.Modify(w, r)
	default:
		http.Error(w, "无效的访问", http.StatusBadRequest)
	}
}

func (t *Router) Home(w http.ResponseWriter, r *http.Request) {
	HandHome(w)
}

func (t *Router) Login(w http.ResponseWriter, r *http.Request){
	HandLogin(w, r)
}

func (t *Router) Register(w http.ResponseWriter, r *http.Request) {
	HandRegister(w, r)
}

func (t *Router) RegisterUpload(w http.ResponseWriter, r *http.Request) {
	HandRegisterUpload(w, r)
}

func (t *Router) Modify(w http.ResponseWriter, r *http.Request) {
	HandModify(w, r)
}
