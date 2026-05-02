package main

import (
	"fmt"
	"log"
	"parking/config"
	"parking/models"
	"parking/router"
)

func main() {
	config.Init()

	if err := models.Init(); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}
	log.Println("数据库初始化成功")

	r := router.SetupRouter()

	port := config.AppConfig.Server.Port
	addr := fmt.Sprintf(":%d", port)

	log.Printf("服务器启动在 http://localhost:%d", port)
	log.Printf("默认账号: admin / 123456")

	if err := r.Run(addr); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
