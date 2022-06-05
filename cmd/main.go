package main

import (
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"main/cache_layer"
	"main/controller_layer"
	"main/helper"
	"main/logger"
	"main/repository_layer"
	"main/service_layer"
	"main/tracer"
	"strconv"
	"time"
)

func main() {

	logger.InitializeLogger(helper.GetEnv("LOG_LEVEL", "INFO"), helper.GetEnv("LOG_FILE", "logfile"))
	tracer.InitializeTracer()
	router := gin.New()
	router.Use(otelgin.Middleware("gin-docker-client-api"))
	initializeLayers().Setup(router)
	router.Run(":" + helper.GetEnv("APP_PORT", "8080"))
}

func initializeLayers() controller_layer.ContainerControllerInterface {
	dockerClient, _ := client.NewEnvClient()
	cacheDefaultExpirationTime, _ := strconv.Atoi(helper.GetEnv("CACHE_DEFAULT_EXPIRATION_TIME", "10"))
	cacheCleanUpIntervalTime, _ := strconv.Atoi(helper.GetEnv("CACHE_CLEANUP_INTERVAL_TIME", "10"))
	repository := repository_layer.NewContainerRepository(cache_layer.NewCache(time.Duration(cacheDefaultExpirationTime)*time.Second, time.Duration(cacheCleanUpIntervalTime)*time.Second), dockerClient)
	service := service_layer.NewContainerService(repository)
	return controller_layer.NewContainerController(service)
}