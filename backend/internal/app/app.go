package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/TryHanger/digital_signage/backend/internal/model"

	"github.com/TryHanger/digital_signage/backend/internal/config"
	handler2 "github.com/TryHanger/digital_signage/backend/internal/handler"
	repository2 "github.com/TryHanger/digital_signage/backend/internal/repository"
	service2 "github.com/TryHanger/digital_signage/backend/internal/service"
	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
)

func Run() {
	cfg := config.Load()
	db := repository2.InitDB(cfg)

	//db.Migrator().DropTable(&model.User{}, &model.RefreshSession{}, &model.Location{}, &model.Monitor{}, &model.MonitorGroup{}, &model.Content{}, &model.Schedule{}, &model.ScheduleBlock{}, &model.Schedule{}, &model.ScheduleBlockItem{}, &model.ScheduleException{}, &model.Template{}, &model.TemplateBlock{}, &model.TemplateContent{})
	//db.AutoMigrate(&model.User{}, &model.RefreshSession{}, &model.Location{}, &model.Monitor{}, &model.MonitorGroup{}, &model.Content{}, &model.Schedule{}, &model.ScheduleBlock{}, &model.ScheduleBlockItem{}, &model.ScheduleException{}, &model.Template{}, &model.TemplateBlock{}, &model.TemplateContent{})

	//db.Migrator().DropTable(
	//	// 1. Независимые таблицы (ни от кого не зависят)
	//	&model.User{},
	//	&model.Location{},
	//	&model.Content{},
	//	&model.Template{},
	//
	//	// 2. Таблицы, зависящие от User/Location/Content
	//	&model.Monitor{},         // зависит от Location
	//	&model.MonitorGroup{},    // зависит от Monitor
	//	&model.TemplateBlock{},   // зависит от Template
	//	&model.TemplateContent{}, // зависит от Template, Content
	//
	//	// 3. Основные таблицы расписаний
	//	&model.Schedule{}, // зависит от User, Content, Template
	//
	//	// 4. Компоненты расписаний (зависят от Schedule)
	//	&model.ScheduleBlock{},     // зависит от Schedule
	//	&model.ScheduleBlockItem{}, // зависит от ScheduleBlock, Content
	//	&model.ScheduleException{}, // зависит от Schedule
	//
	//	// 5. Сессии (зависят от User)
	//	&model.RefreshSession{}, // зависит от User
	//)
	//
	//db.Exec("DROP TABLE IF EXISTS schedule_monitors CASCADE")

	db.AutoMigrate(
		// 1. Независимые таблицы (ни от кого не зависят)
		&model.User{},
		&model.Location{},
		&model.Content{},
		&model.Template{},

		// 2. Таблицы, зависящие от User/Location/Content
		&model.Monitor{},         // зависит от Location
		&model.MonitorGroup{},    // зависит от Monitor
		&model.TemplateBlock{},   // зависит от Template
		&model.TemplateContent{}, // зависит от Template, Content

		// 3. Основные таблицы расписаний
		&model.Schedule{}, // зависит от User, Content, Template

		// 4. Компоненты расписаний (зависят от Schedule)
		&model.ScheduleBlock{},     // зависит от Schedule
		&model.ScheduleBlockItem{}, // зависит от ScheduleBlock, Content
		&model.ScheduleException{}, // зависит от Schedule

		// 5. Сессии (зависят от User)
		&model.RefreshSession{}, // зависит от User
	)

	// --- Repositories ---
	userRepo := repository2.NewUserRepository(db)
	authRepo := repository2.NewAuthRepository(db)
	monitorRepo := repository2.NewMonitorRepository(db)
	contentRepo := repository2.NewContentRepository(db)
	scheduleRepo := repository2.NewScheduleRepository(db)
	locationRepo := repository2.NewLocationRepository(db)
	templateRepo := repository2.NewTemplateRepository(db)

	// --- Cache ---
	//scheduleCache := cache.NewScheduleCache()

	// --- Notifier ---
	//notifier := socket.NewWebSocketNotifier(monitorRepo, scheduleCache)

	// --- Services ---
	userService := service2.NewUserService(userRepo)
	tokenService := service2.NewTokenService(authRepo, cfg.JWTSecret)
	monitorService := service2.NewMonitorService(monitorRepo)
	contentService := service2.NewContentService(contentRepo)
	scheduleService := service2.NewScheduleService(scheduleRepo)
	locationService := service2.NewLocationService(locationRepo)
	templateService := service2.NewTemplateService(templateRepo)

	// --- Handlers ---
	userHandler := handler2.NewUserHandler(userService, tokenService)
	monitorHandler := handler2.NewMonitorHandler(monitorService)
	contentHandler := handler2.NewContentHandler(contentService)
	scheduleHandler := handler2.NewScheduleHandler(scheduleService)
	locationHandler := handler2.NewLocationHandler(locationService)
	//cacheHandler := handler2.NewCacheHandler(scheduleCache)
	templateHandler := handler2.NewTemplateHandler(templateService)

	// --- Gin ---
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
			"http://127.0.0.1:3000",
			"http://localhost:5173",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true, // разрешаем куки и авторизацию
	}))

	r.OPTIONS("/*path", func(c *gin.Context) {
		c.Status(200)
	})

	r.RedirectTrailingSlash = false

	// 🔌 WebSocket endpoint
	//r.GET("/ws", func(c *gin.Context) {
	//	conn, err := websocket.Upgrade(c.Writer, c.Request, nil, 1024, 1024)
	//	if err != nil {
	//		log.Println("❌ Ошибка апгрейда соединения:", err)
	//		return
	//	}
	//
	//	// Просто передаём управление сокет-обработчику
	//	notifier.HandleConnection(conn)
	//})

	// REST endpoints under /api/v1
	api := r.Group("/api/v1")
	//api.GET("/cache/schedules", cacheHandler.GetCache)
	userHandler.RegisterRoutes(api)
	monitorHandler.RegisterRoutes(api)
	contentHandler.RegisterRoutes(api)
	scheduleHandler.RegisterRoutes(api)
	locationHandler.RegisterRoutes(api)
	templateHandler.RegisterRoutes(api)

	//scheduleService.StartScheduler()
	//
	//// 🔔 Событие при подключении нового монитора
	//notifier.OnConnect(func(monitorID uint) {
	//	scheduleService.SendSchedulesToMonitor(monitorID)
	//})

	// ⏰ Запуск планировщика
	// scheduleService.StartScheduler()

	// 🚀 Старт сервера
	fmt.Println("🚀 Сервер запущен на порту", cfg.ServerPort)
	err := http.ListenAndServe(":"+cfg.ServerPort, r)
	if err != nil {
		log.Fatal(err)
	}
}
