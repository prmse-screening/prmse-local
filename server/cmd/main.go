package main

import (
	"github.com/gin-gonic/gin"
	"server/internal/config"
	"server/internal/handlers"
	"server/internal/logger"
	"server/internal/middlewares"
	"server/internal/schedule"
)

func NewServer(task *handlers.TasksHandler, ts *schedule.TasksScheduler) *gin.Engine {
	//engine := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Recovery(), middlewares.Logger())
	engine.GET("/ping", func(c *gin.Context) {
		c.JSON(200, "pong")
	})
	engine.GET("/test", func(c *gin.Context) {
		c.JSON(200, task.CreateTask)
	})
	go ts.Start()
	return engine
}

func main() {
	config.Init()
	logger.Init()
	engine, _ := wireApp()
	engine.Run()
}
