package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-mcp-demo/internal/config"
	"go-mcp-demo/internal/mcp"
	"go-mcp-demo/internal/model"
	"go-mcp-demo/internal/prompt"
	"io"
	"log"
	"net/http"
	"strings"
)

// TalkHandler 根据用户输入多轮次模型对话
func TalkHandler(c *gin.Context) {
	var req model.TalkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sendError(c, http.StatusBadRequest, "Invalid JSON")
		return
	}

	promptText := req.Prompt

	// First classification
	if category, err := classify(prompt.FirstPrompt(promptText)); err != nil || category != "b" {
		streamOllama(promptText, c)
		return
	}

	// Second classification
	if category, err := classify(prompt.SecondPrompt(promptText)); err != nil || category != "b" {
		streamOllama(promptText, c)
		return
	}

	sqlFileContent := mcp.CallMCPTool("read_file", promptText)
	if sqlFileContent == "" {
		streamOllama(promptText, c)
		return
	}

	// Third classification
	sqlOrIndicator, err := classify(prompt.ThirdPrompt(promptText, sqlFileContent))
	if err != nil {
		sendError(c, http.StatusInternalServerError, "Failed classification")
		return
	}

	// Handle inability to answer with SQL
	if strings.HasPrefix(sqlOrIndicator, "a") {
		response := formatFallbackResponse(promptText, sqlFileContent, "", "")
		streamOllama(response, c)
		return
	}

	// Execute SQL query
	queryResult := mcp.CallMCPTool("query_db", sqlOrIndicator)
	if strings.TrimSpace(queryResult) == "" {
		response := formatFallbackResponse(promptText, sqlFileContent, sqlOrIndicator, "")
		streamOllama(response, c)
		return
	}

	// Success response
	response := formatFallbackResponse(promptText, sqlFileContent, sqlOrIndicator, queryResult)
	streamOllama(response, c)
}

// classify sends a prompt to Ollama and returns the classification result.
func classify(prompt string) (string, error) {
	req := map[string]interface{}{
		"prompt": prompt,
		"model":  config.AppConfig.Ollama.Model,
		"stream": true,
	}

	body, _ := json.Marshal(req)
	resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	var full strings.Builder

	for {
		var chunk map[string]interface{}
		if err := decoder.Decode(&chunk); err == io.EOF {
			break
		} else if err != nil {
			return "", err
		}
		if val, ok := chunk["response"].(string); ok {
			full.WriteString(val)
		}
	}

	return strings.TrimSpace(full.String()), nil
}

// streamOllama streams response from Ollama to client.
func streamOllama(prompt string, c *gin.Context) error {
	req := map[string]interface{}{
		"prompt": prompt,
		"model":  config.AppConfig.Ollama.Model,
		"stream": true,
	}
	body, _ := json.Marshal(req)
	resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(body))
	if err != nil {
		sendError(c, http.StatusInternalServerError, "Failed to connect to Ollama")
		return err
	}
	defer resp.Body.Close()

	c.Writer.Header().Set("Content-Type", "application/json")
	c.Status(http.StatusOK)

	decoder := json.NewDecoder(resp.Body)
	encoder := json.NewEncoder(c.Writer)

	for {
		var chunk map[string]interface{}
		if err := decoder.Decode(&chunk); err == io.EOF {
			break
		} else if err != nil {
			log.Println("Stream decode error:", err)
			break
		}
		if val, ok := chunk["response"].(string); ok {
			if err := encoder.Encode(model.TalkResponse{Response: val}); err != nil {
				log.Println("Stream encode error:", err)
				break
			}
			if flusher, ok := c.Writer.(http.Flusher); ok {
				flusher.Flush()
			}
		}
	}

	return nil
}

// sendError sends an error response with standard format.
func sendError(c *gin.Context, status int, msg string) {
	c.JSON(status, gin.H{"error": msg})
}

// formatFallbackResponse formats fallback or success response.
func formatFallbackResponse(prompt, sqlFile, sql, result string) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("以下是用户的问题：\n%s\n\n", prompt))
	b.WriteString(fmt.Sprintf("下面是 SQL 文件的内容：\n%s\n\n", sqlFile))
	if sql != "" {
		b.WriteString(fmt.Sprintf("下面是 SQL 语句：\n%s\n\n", sql))
	}
	if result != "" {
		b.WriteString(fmt.Sprintf("下面是 SQL 查询结果：\n%s\n\n", result))
	}
	return b.String()
}
