package app

import (
	"fmt"
	"github.com/TryHanger/digital_signage/internal/config"
	"github.com/TryHanger/digital_signage/internal/handler"
	"github.com/TryHanger/digital_signage/internal/model"
	"github.com/TryHanger/digital_signage/internal/repository"
	"github.com/TryHanger/digital_signage/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Run() {
	cfg := config.Load()
	db := repository.InitDB(cfg)
	//db.Migrator().DropTable(&model.Location{}, &model.Monitor{}, &model.Content{}, &model.Schedule{}, &model.ScheduleDay{})
	db.AutoMigrate(&model.Location{}, &model.Monitor{}, &model.Content{}, &model.Schedule{}, &model.ScheduleDay{})

	// --- Repos/Services/Handlers ---
	monitorRepo := repository.NewMonitorRepository(db)
	monitorService := service.NewMonitorService(monitorRepo)
	monitorHandler := handler.NewMonitorHandler(monitorService)

	contentRepo := repository.NewContentRepository(db)
	contentService := service.NewContentService(contentRepo)
	contentHandler := handler.NewContentHandler(contentService)

	scheduleRepo := repository.NewScheduleRepository(db)
	scheduleService := service.NewScheduleService(scheduleRepo)
	scheduleHandler := handler.NewScheduleHandler(scheduleService)

	locationRepo := repository.NewLocationRepository(db)
	locationService := service.NewLocationService(locationRepo)
	locationHandler := handler.NewLocationHandler(locationService)

	// --- Инициализация Socket.IO ---
	socketServer := InitSocketServer()
	go socketServer.Serve()
	defer socketServer.Close()

	// --- Gin маршруты ---
	r := gin.Default()

	// Подключаем REST-хендлеры
	monitorHandler.RegisterRoutes(r)
	contentHandler.RegisterRoutes(r)
	scheduleHandler.RegisterRoutes(r)
	locationHandler.RegisterRoutes(r)

	// Подключаем Socket.IO endpoint
	r.GET("/socket.io/*any", gin.WrapH(socketServer))
	r.POST("/socket.io/*any", gin.WrapH(socketServer))

	//sch := utils.NewScheduler(scheduleService, socketServer, 30*time.Second)
	//sch.Start()

	// Запуск
	fmt.Println("🚀 Сервер запущен на порту", cfg.ServerPort)
	http.ListenAndServe(":"+cfg.ServerPort, r)
}
