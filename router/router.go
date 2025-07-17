package router

import (
	"ReadBook/controller"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

func InitGin() *gin.Engine {
	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}                                                 // 允许前端域名
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}           // 允许的HTTP方法
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"} // 允许的请求头
	config.AllowCredentials = true                                                      // 允许携带凭证（如cookie）
	config.MaxAge = 12 * time.Hour
	router.Use(cors.New(config))
	router.GET("/api/index", controller.GetAllBooks)
	router.GET("/api/detail", controller.GetBookDetails)
	return router
}
