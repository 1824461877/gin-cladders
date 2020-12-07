package controller

import (
	"gin-cladder/middleware"
	"gin-cladder/server_container/mysql"
	"github.com/gin-gonic/gin"
)

type AdminLoginController struct {}

func AdminLoginRegister(group *gin.RouterGroup) {
	adminlogin := &AdminLoginController{}
	group.GET("/login",adminlogin.AdminLogin)
}


var s = make(chan int,50)

func (adminlogin *AdminLoginController) AdminLogin(c * gin.Context){
	s <- 1
	go func() {
		<- s
		mysql.InitMysqlDB()
	}()
	middleware.ResponseSuccess(c,"this is good response success")
}
