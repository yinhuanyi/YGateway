/**
 * @Author：Robby
 * @Date：2022/1/18 18:02
 * @Function：
 **/

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (b *BackendServer) HelloHandler(w http.ResponseWriter, req *http.Request) {
	userPath := fmt.Sprintf("http://%s%s\n", b.Addr, req.URL.Path)
	realIP := fmt.Sprintf("RemoteAddr = %s, X-Forwarded-For = %+v, X-Real-Ip = %v\n", req.RemoteAddr, req.Header.Get("X-Forwarded-For"), req.Header.Get("X-Real-Ip"))
	header:=fmt.Sprintf("headers = %+v",req.Header)

	_, err := fmt.Fprintln(w, userPath, realIP, header)
	if err != nil {
		log.Println(err)
	}
}

func (b *BackendServer) ErrorHandler(w http.ResponseWriter, req *http.Request) {
	userPath := "error handler"
	w.WriteHeader(500)
	_, err := fmt.Fprint(w, userPath)
	if err != nil {
		log.Println(err)
	}
}

func (b *BackendServer) TimeoutHandler(w http.ResponseWriter, req *http.Request) {
	time.Sleep(6*time.Second)
	userPath := "timeout handler"
	w.WriteHeader(200)
	_, err := fmt.Fprint(w, userPath)
	if err != nil {
		log.Println(err)
	}
}

func (b *BackendServer) Run() (srv *http.Server) {
	log.Println("Starting httpserver at " + b.Addr)

	mux := http.NewServeMux()
	mux.HandleFunc("/", b.HelloHandler)
	mux.HandleFunc("/base/error", b.ErrorHandler)
	mux.HandleFunc("/test_http_timeout/timeout", b.TimeoutHandler)

	srv = &http.Server{
		Addr:         b.Addr,
		WriteTimeout: time.Second * 3,
		Handler:      mux,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen err: %s\n", err)
		}
	}()

	return
}

type BackendServer struct {
	Addr string
}


func main() {
	backend1 := &BackendServer{Addr: "127.0.0.1:8081"}
	srv1 := backend1.Run()

	backend2 := &BackendServer{Addr: "127.0.0.1:8082"}
	srv2 := backend2.Run()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("准备关闭HTTP服务器...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 关闭HTTP服务，且清理连接
	if err := srv1.Shutdown(ctx); err != nil {
		log.Fatal("HTTP服务器退出失败")
	}

	if err := srv2.Shutdown(ctx); err != nil {
		log.Fatal("HTTP服务器退出失败")
	}

	log.Println("HTTP服务器退出完毕")
}