package main

import (
	"go-mcp-demo/internal/config"
	"go-mcp-demo/internal/mcp"
	"log"
)

func main() {
	err := config.LoadConfig("../../pkg/config/config.yml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}
	router, err := mcp.InitMCPServer()
	if err != nil {
		panic(err)
	}

	if err := router.Run(":2002"); err != nil {
		panic(err)
	}
}
