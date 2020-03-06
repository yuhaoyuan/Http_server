package main

import (
	"log"
	"net/http"
	"time"
)

// Router 路由
type Router struct {
}

// ServeHTTP .
func (t *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	uPath := r.URL.Path
	//log.Println("url call - \n", uPath)
	switch uPath {
	case "/":
		t.Home(w, r)
	case "/home":
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

// Home .
func (t *Router) Home(w http.ResponseWriter, r *http.Request) {
	bT := time.Now() // 开始时间
	HandHome(w, r)
	eT := time.Since(bT) // 从开始到当前所消耗的时间
	log.Println("func time log---------------api:home, time=\n\n", eT)
}

// Login .
func (t *Router) Login(w http.ResponseWriter, r *http.Request) {
	//bT := time.Now()
	HandLogin(w, r)
	//eT := time.Since(bT)
	//log.Println("func time log---------------api:login, time=\n\n", eT)
}

// Register .
func (t *Router) Register(w http.ResponseWriter, r *http.Request) {
	bT := time.Now()
	HandRegister(w, r)
	eT := time.Since(bT)
	log.Println("func time log---------------api:register, time=\n\n", eT)
}

// RegisterUpload .
func (t *Router) RegisterUpload(w http.ResponseWriter, r *http.Request) {
	bT := time.Now()
	HandRegisterUpload(w, r)
	eT := time.Since(bT)
	log.Println("func time log---------------api:register_upload, time=\n\n", eT)
}

// Modify .
func (t *Router) Modify(w http.ResponseWriter, r *http.Request) {
	bT := time.Now()
	HandModify(w, r)
	eT := time.Since(bT)
	log.Println("func time log---------------api:modify, time=\n\n", eT)
}
