package handlers

import (
	"net/http"
	"task-manager/internal/config"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
	"os"
)

type ConfigHandler struct{}

type APIKeyRequest struct {
	APIKey string `json:"api_key"`
	Model  string `json:"model"`
}

func NewConfigHandler() *ConfigHandler {
	return &ConfigHandler{}
}

func (h *ConfigHandler) GetOpenAIConfig(c *gin.Context) {
	cfg := config.GetConfig()
	
	c.JSON(http.StatusOK, gin.H{
		"api_key": cfg.OpenAI.APIKey,
		"model": cfg.OpenAI.Model,
		"use_llm": cfg.OpenAI.UseLLM,
	})
}

func (h *ConfigHandler) UpdateOpenAIConfig(c *gin.Context) {
	var req APIKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update config in memory
	cfg := config.GetConfig()
	cfg.OpenAI.APIKey = req.APIKey
	
	// Update model if provided
	if req.Model != "" {
		cfg.OpenAI.Model = req.Model
	}
	
	// Save to config file
	configPath := "configs/config.yaml"
	data, err := os.ReadFile(configPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read config file"})
		return
	}
	
	var fileConfig config.Config
	if err := yaml.Unmarshal(data, &fileConfig); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse config file"})
		return
	}
	
	fileConfig.OpenAI.APIKey = req.APIKey
	if req.Model != "" {
		fileConfig.OpenAI.Model = req.Model
	}
	
	updatedData, err := yaml.Marshal(fileConfig)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate updated config"})
		return
	}
	
	if err := os.WriteFile(configPath, updatedData, 0644); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write config file"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "OpenAI configuration updated successfully",
		"api_key": req.APIKey,
		"model": cfg.OpenAI.Model,
	})
}

func (h *ConfigHandler) RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api/v1")
	{
		config := api.Group("/config")
		{
			config.GET("/openai", h.GetOpenAIConfig)
			config.PUT("/openai", h.UpdateOpenAIConfig)
		}
	}
} 