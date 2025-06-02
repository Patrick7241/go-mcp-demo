package mcp

import (
	"context"
	"fmt"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	"log"
	"time"
)

// CallMCPTool 调用任意 MCP 工具
// toolName 是工具名，如 "query_db" 或自定义工具名
// args 是传递给工具的参数（如 SQL 字符串），可为空
func CallMCPTool(toolName string, args any) string {
	fmt.Printf("正在调用工具 %s...\n", toolName)
	// 创建 MCP 客户端（SSE）
	mcpClient, err := client.NewSSEMCPClient("http://localhost:2002/sse")
	if err != nil {
		panic(fmt.Errorf("创建 MCP 客户端失败: %w", err))
	}
	defer mcpClient.Close()

	// 设置超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 启动客户端
	if err := mcpClient.Start(ctx); err != nil {
		panic(fmt.Errorf("启动 MCP 客户端失败: %w", err))
	}

	// 初始化请求
	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "Client Demo",
		Version: "1.0.0",
	}
	initResult, err := mcpClient.Initialize(ctx, initRequest)
	if err != nil {
		panic(fmt.Errorf("初始化失败: %w", err))
	}
	fmt.Printf("初始化成功，服务器信息: %s %s\n", initResult.ServerInfo.Name, initResult.ServerInfo.Version)

	// 获取工具列表
	tools, err := mcpClient.ListTools(ctx, mcp.ListToolsRequest{})
	if err != nil {
		panic(fmt.Errorf("获取工具列表失败: %w", err))
	}
	fmt.Println("可用工具列表:")
	for _, tool := range tools.Tools {
		fmt.Printf("- %s: %s\n", tool.Name, tool.Description)
	}

	// 构建工具请求
	toolRequest := mcp.CallToolRequest{
		Request: mcp.Request{Method: "tools/call"},
		Params: mcp.CallToolParams{
			Name:      toolName,
			Arguments: args,
		},
	}

	fmt.Printf("参数是：%v\n", args)

	// 调用工具
	result, err := mcpClient.CallTool(ctx, toolRequest)
	if err != nil {
		log.Printf("调用工具 %s 失败: %v", toolName, err)
		return "调用工具失败"
	}

	// 提取并返回结果
	if len(result.Content) == 0 {
		return "没有返回结果"
	}
	if textContent, ok := result.Content[0].(mcp.TextContent); ok {
		return textContent.Text
	}
	return "返回结果格式不正确"
}
