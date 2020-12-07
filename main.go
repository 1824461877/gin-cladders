package main

import (
	"gin-cladder/conf/elite/control"
	"gin-cladder/conf/elite/extend"
	"gin-cladder/router"
	"os"
	"os/signal"
	"syscall"
)


func main() {
	if err := control.InitModule(); err != nil {
		panic(err)
	}
	// 路由配置
	rs := extend.InitServerConf{router.Router}
	rs.HttpServerRun()
	defer control.Destroy()
	quit := make(chan os.Signal)

	signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	rs.HttpServerStop()
}