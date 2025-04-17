package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"server/internal/config"
	"server/internal/handlers"
	"server/internal/logger"
	"server/internal/middlewares"
	"server/internal/schedule"
	"syscall"
	"time"
)

func NewServer(t *handlers.TasksHandler, d *handlers.DicomHandler, ts *schedule.TasksScheduler) *http.Server {
	//gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Recovery(), middlewares.Logger(), cors.Default())
	engine.GET("/ping", func(c *gin.Context) {
		c.JSON(200, "pong")
	})

	tasks := engine.Group("/tasks")
	{
		tasks.POST("/create", t.CreateTask)
		tasks.POST("/update", t.UpdateTask)
		tasks.POST("/prioritize", t.PrioritizeTask)
		tasks.POST("/delete", t.DeleteTask)
		tasks.POST("/device", t.SetWorkerDevice)
		tasks.GET("/upload", t.GetUploadUrl)
		tasks.GET("/list", t.GetListPagination)
	}

	dicom := engine.Group("/dicom")
	{
		dicom.GET("/:series/:file", d.Redirect)
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Cfg.App.Port),
		Handler: engine,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Errorf("Shutdown Server with error: %v", err)
		}
	}()

	//go ts.Start()
	return srv
}

func main() {
	config.Init()
	logger.Init()
	srv, _ := wireApp()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Infof("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Infof("Server forced to shutdown: %v", err)
	}
}
