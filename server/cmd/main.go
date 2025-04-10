package cmd

import (
	"github.com/gin-gonic/gin"
	"server/internal/config"
	"server/internal/handlers"
	"server/internal/logger"
	"server/internal/schedule"
)

func NewServer(task *handlers.TasksHandler, ts *schedule.TasksScheduler) *gin.Engine {
	engine := gin.Default()
	engine.GET("/ping", func(c *gin.Context) {
		c.JSON(200, task.CreateTask)
	})
	go ts.Start()
	return engine
}

func main() {
	config.Init()
	logger.Init()
	//engine := wireApp()
	//engine.Run(":8000")
}
