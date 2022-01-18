/**
 * @Author：Robby
 * @Date：2022/1/18 18:15
 * @Function：
 **/

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

var addr = "127.0.0.1:7070"

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

func joinURLPath(a, b *url.URL) (path, rawpath string) {
	if a.RawPath == "" && b.RawPath == "" {
		return singleJoiningSlash(a.Path, b.Path), ""
	}
	apath := a.EscapedPath()
	bpath := b.EscapedPath()

	aslash := strings.HasSuffix(apath, "/")
	bslash := strings.HasPrefix(bpath, "/")

	switch {
	case aslash && bslash:
		return a.Path + b.Path[1:], apath + bpath[1:]
	case !aslash && !bslash:
		return a.Path + "/" + b.Path, apath + "/" + bpath
	}
	return a.Path + b.Path, apath + bpath
}

func main() {
	backendSrv := "http://127.0.0.1:8081/base"
	parseUrl, err := url.Parse(backendSrv)
	if err != nil {
		log.Println(err)
	}

	// 创建director请求修改函数
	targetQuery := parseUrl.RawQuery
	director := func(req *http.Request) {
		req.URL.Scheme = parseUrl.Scheme
		req.URL.Host = parseUrl.Host
		req.URL.Path, req.URL.RawPath = joinURLPath(parseUrl, req.URL)
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			req.Header.Set("User-Agent", "")
		}
	}

	// 创建modifyResponse响应修改函数
	modifyResponse := func(resp *http.Response) error {
		// 读取下游服务器响应数据
		respData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		// 在原有的响应数据之上，添加其他的数据
		newRespData := []byte("Proxy has handled: " + string(respData))
		// resp.Body是io.ReadCloser接口类型，赋值的对象必须实现Reader和Closer两个方法，那么就需要通过ioutil.NopCloser()方法返回实现了io.ReadCloser接口类型接口类型的变量
		resp.Body = ioutil.NopCloser(bytes.NewBuffer(newRespData))
		// 重新计算响应数据长度
		resp.ContentLength = int64(len(newRespData))
		// 在header头部添加包长度
		resp.Header.Set("Content-Length", fmt.Sprintf("%s", resp.Body))
		return nil
	}

	// 创建reverseProxy实例
	proxy := &httputil.ReverseProxy{
		Director: director,
		ModifyResponse: modifyResponse,
	}


	log.Println("Starting ProxyServer at " + addr)
	// 启动HTTP服务器
	log.Fatal(http.ListenAndServe(addr, proxy))
}