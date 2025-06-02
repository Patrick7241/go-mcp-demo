package mcp

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"go-mcp-demo/internal/config"
	"os"
	"time"
)

var db *sql.DB // 全局数据库连接变量

// InitMCPServer 初始化 MCP 服务（供主服务调用）
func InitMCPServer() (*gin.Engine, error) {
	// 初始化数据库连接
	if err := initDB(); err != nil {
		return nil, err
	}

	// 创建 MCP 服务
	mcpServer := server.NewMCPServer("my-mcp-server", "1.0.0", server.WithLogging())

	// 注册工具
	registerMCPTools(mcpServer)

	// 创建 SSE 服务
	sseServer := server.NewSSEServer(mcpServer)
	expectedPath := sseServer.CompleteSsePath()
	//fmt.Printf("SSE路径: %s\n", expectedPath)

	// 创建 Gin 路由
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"POST", "GET", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 注册路由处理
	r.GET(expectedPath, func(c *gin.Context) {
		sseServer.ServeHTTP(c.Writer, c.Request)
	})
	r.POST("/message", func(c *gin.Context) {
		sseServer.ServeHTTP(c.Writer, c.Request)
	})

	return r, nil
}

// initDB 初始化数据库连接
func initDB() error {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		config.AppConfig.Mysql.User,
		config.AppConfig.Mysql.Password,
		config.AppConfig.Mysql.Host,
		config.AppConfig.Mysql.Port,
		config.AppConfig.Mysql.Database,
	)
	//fmt.Printf("DSN: %s\n", dsn)
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("数据库连接失败: %v", err)
	}

	// 检查连接
	if err = db.Ping(); err != nil {
		return fmt.Errorf("数据库Ping失败: %v", err)
	}
	return nil
}

// registerMCPTools 注册 MCP 工具
func registerMCPTools(s *server.MCPServer) {
	s.AddTool(newReadFileTool(), handleReadFile)
	s.AddTool(newQueryDBTool(), handleQueryDB)
}

// newReadFileTool 创建读取文件工具
func newReadFileTool() mcp.Tool {
	return mcp.NewTool("read_file",
		mcp.WithDescription("读取指定路径下的sql文件"),
	)
}

// handleReadFile 处理读取文件逻辑
func handleReadFile(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	filePath := config.AppConfig.Sql.SqlFilePath

	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %v", err)
	}

	return mcp.NewToolResultText(string(content)), nil
}

// newQueryDBTool 创建数据库查询工具
func newQueryDBTool() mcp.Tool {
	return mcp.NewTool("query_db",
		mcp.WithDescription("执行 SQL 查询并返回结果"),
	)
}

// handleQueryDB 处理数据库查询逻辑
func handleQueryDB(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	sqlQuery := request.Params.Arguments.(string)

	fmt.Printf("执行SQL: %s\n", sqlQuery)

	rows, err := db.QueryContext(ctx, sqlQuery)
	if err != nil {
		return nil, fmt.Errorf("SQL执行失败: %v", err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("获取列名失败: %v", err)
	}

	var results []map[string]interface{}

	for rows.Next() {
		values := make([]interface{}, len(columns))
		pointers := make([]interface{}, len(columns))
		for i := range values {
			pointers[i] = &values[i]
		}

		if err := rows.Scan(pointers...); err != nil {
			return nil, fmt.Errorf("扫描数据失败: %v", err)
		}

		row := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			if b, ok := val.([]byte); ok {
				row[col] = string(b)
			} else {
				row[col] = val
			}
		}
		results = append(results, row)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("遍历结果失败: %v", err)
	}

	jsonBytes, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("JSON编码失败: %v", err)
	}

	return mcp.NewToolResultText(string(jsonBytes)), nil
}
