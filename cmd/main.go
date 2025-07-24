package main

import (
	"eino-chatbot/config"
	"eino-chatbot/handlers"
	"eino-chatbot/services"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化配置
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Sprintf("加载配置失败: %v", err))
	}

	// 初始化 Volcengine 服务
	volcService := services.NewVolcengineService(cfg)

	// 初始化 Gin 路由
	r := gin.Default()

	// 服务静态文件
	r.Static("/static", "./static")
	r.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})
	r.GET("/stream", func(c *gin.Context) {
		c.File("./status/stream.html")
	})

	// 注册路由
	handlers.RegisterRoutes(r, volcService)

	// 启动服务
	addr := fmt.Sprintf(":%d", cfg.Port)
	if err = r.Run(addr); err != nil {
		panic(fmt.Sprintf("启动服务器失败: %v", err))
	}
}
