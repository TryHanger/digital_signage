package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/TryHanger/digital_signage/internal/cache"
	"github.com/TryHanger/digital_signage/internal/config"
	"github.com/TryHanger/digital_signage/internal/handler"
	"github.com/TryHanger/digital_signage/internal/model"
	"github.com/TryHanger/digital_signage/internal/repository"
	"github.com/TryHanger/digital_signage/internal/service"
	"github.com/TryHanger/digital_signage/internal/socket"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func Run() {
	cfg := config.Load()
	db := repository.InitDB(cfg)

	//db.Migrator().DropTable(&model.Location{}, &model.Monitor{}, &model.Content{}, &model.Schedule{}, &model.ScheduleDay{})
	db.AutoMigrate(&model.Location{}, &model.Monitor{}, &model.Content{}, &model.Schedule{}, &model.ScheduleDay{})

	// --- Repositories ---
	monitorRepo := repository.NewMonitorRepository(db)
	contentRepo := repository.NewContentRepository(db)
	scheduleRepo := repository.NewScheduleRepository(db)
	locationRepo := repository.NewLocationRepository(db)

	// --- Cache ---
	scheduleCache := cache.NewScheduleCache()

	// --- Notifier ---
	notifier := socket.NewWebSocketNotifier(monitorRepo, scheduleCache)

	// --- Services ---
	monitorService := service.NewMonitorService(monitorRepo)
	contentService := service.NewContentService(contentRepo)
	scheduleService := service.NewScheduleService(scheduleRepo, scheduleCache, notifier)
	locationService := service.NewLocationService(locationRepo)

	// --- Handlers ---
	monitorHandler := handler.NewMonitorHandler(monitorService)
	contentHandler := handler.NewContentHandler(contentService)
	scheduleHandler := handler.NewScheduleHandler(scheduleService)
	locationHandler := handler.NewLocationHandler(locationService)
	cacheHandler := handler.NewCacheHandler(scheduleCache)

	// --- Gin ---
	r := gin.Default()

	// 🔌 WebSocket endpoint
	r.GET("/ws", func(c *gin.Context) {
		conn, err := websocket.Upgrade(c.Writer, c.Request, nil, 1024, 1024)
		if err != nil {
			log.Println("❌ Ошибка апгрейда соединения:", err)
			return
		}

		// Просто передаём управление сокет-обработчику
		notifier.HandleConnection(conn)
	})

	// REST endpoints
	r.GET("/cache/schedules", cacheHandler.GetCache)
	monitorHandler.RegisterRoutes(r)
	contentHandler.RegisterRoutes(r)
	scheduleHandler.RegisterRoutes(r)
	locationHandler.RegisterRoutes(r)

	scheduleService.StartScheduler()

	// 🔔 Событие при подключении нового монитора
	notifier.OnConnect(func(monitorID uint) {
		scheduleService.SendSchedulesToMonitor(monitorID)
	})

	// ⏰ Запуск планировщика
	// scheduleService.StartScheduler()

	// 🚀 Старт сервера
	fmt.Println("🚀 Сервер запущен на порту", cfg.ServerPort)
	http.ListenAndServe(":"+cfg.ServerPort, r)
}
