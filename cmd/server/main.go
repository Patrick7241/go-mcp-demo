package main

import (
	"go-mcp-demo/internal/api"
	"go-mcp-demo/internal/config"
	"log"
)

// 主服务入口
func main() {
	err := config.LoadConfig("../../pkg/config/config.yml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}
	r := api.NewRouter()
	log.Println("Listening on :2001...")
	if err := r.Run(":2001"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
