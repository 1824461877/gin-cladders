package router

import (
	"gin-cladder/controller"
	"gin-cladder/middleware"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	router := gin.Default()
	adminLogin := router.Group("/admin_login")
	adminLogin.Use(
		middleware.IPAuthMiddleware(),
		middleware.RecoveryMiddleware(),
		)
	{
		controller.AdminLoginRegister(adminLogin)
	}
	return router
}
