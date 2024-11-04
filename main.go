package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"time"
)

func main() {
	// 定义命令行参数
	port := flag.String("p", "8080", "HTTP server port")
	dir := flag.String("dir", ".", "Directory to serve")

	// 解析命令行参数
	flag.Parse()

	// 打印参数值
	fmt.Printf("Starting server on port %s, serving directory %s\n", *port, *dir)

	go openBrowser("http://localhost:" + *port)

	// 启动静态文件服务器
	http.ListenAndServe(":"+*port, loggingMiddleware(http.FileServer(http.Dir(*dir))))
}

func openBrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "darwin": // macOS
		err = exec.Command("open", url).Start()
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	default:
		err = exec.Command("cmd", "/c", "start", url).Start()
	}
	if err != nil {
		panic(err)
	}
}

// 日志中间件
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// 创建一个响应记录器
		lr := &logResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(lr, r)

		duration := time.Since(start)
		log.Printf(
			"[%s] %s %d %s",
			r.Method,
			r.URL.Path,
			lr.statusCode,
			duration,
		)
	})
}

// 自定义的响应记录器
type logResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lr *logResponseWriter) WriteHeader(code int) {
	lr.statusCode = code
	lr.ResponseWriter.WriteHeader(code)
}
