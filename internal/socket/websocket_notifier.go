package socket

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/TryHanger/digital_signage/internal/cache"
	"github.com/TryHanger/digital_signage/internal/model"
	"github.com/TryHanger/digital_signage/internal/repository"
	"github.com/gorilla/websocket"
)

type WebSocketNotifier struct {
	mu            sync.RWMutex
	connections   map[uint]*websocket.Conn
	monitorRepo   *repository.MonitorRepository
	scheduleCache *cache.ScheduleCache
	onConnect     func(monitorID uint)
}

func NewWebSocketNotifier(monitorRepo *repository.MonitorRepository, cache *cache.ScheduleCache) *WebSocketNotifier {
	return &WebSocketNotifier{
		connections:   make(map[uint]*websocket.Conn),
		monitorRepo:   monitorRepo,
		scheduleCache: cache,
	}
}

func (n *WebSocketNotifier) HandleConnection(conn *websocket.Conn) {
	defer conn.Close()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("⚠️ Ошибка чтения: %v", err)
			return
		}

		var payload map[string]interface{}
		if err := json.Unmarshal(msg, &payload); err != nil {
			log.Printf("❌ Ошибка парсинга JSON: %v", err)
			continue
		}

		event := payload["event"].(string)
		data := payload["data"].(map[string]interface{})

		switch event {
		case "register_monitor":
			n.HandleRegister(conn, data)

		default:
			log.Printf("📩 Неизвестное событие: %s", event)
		}
	}
}

func (n *WebSocketNotifier) HandleRegister(conn *websocket.Conn, data map[string]interface{}) {
	id := uint(data["id"].(float64))

	monitor, err := n.monitorRepo.GetByID(id)
	if err != nil {
		log.Printf("❌ Не найден монитор %d: %v", id, err)
		return
	}

	n.mu.Lock()
	n.connections[id] = conn
	n.mu.Unlock()

	log.Printf("🖥️ Монитор подключён: %s (ID: %d)", monitor.Name, monitor.ID)

	if n.onConnect != nil {
		n.onConnect(monitor.ID)
	}

	// Отправляем начальные данные
	schedules := n.scheduleCache.GetByMonitorID(monitor.ID)
	conn.WriteJSON(map[string]interface{}{
		"event": "init_schedules",
		"data":  schedules,
	})
}

func (n *WebSocketNotifier) NotifyScheduleUpdate(monitorID uint, schedule model.Schedule) {
	n.mu.RLock()
	conn, ok := n.connections[monitorID]
	n.mu.RUnlock()

	if !ok {
		log.Printf("⚠️ Монитор %d не подключён", monitorID)
		return
	}

	msg := map[string]interface{}{
		"event": "schedule_update",
		"data":  schedule,
	}

	if err := conn.WriteJSON(msg); err != nil {
		log.Printf("❌ Ошибка отправки данных монитору %d: %v", monitorID, err)
	}
}

func (n *WebSocketNotifier) OnConnect(handler func(monitorID uint)) {
	n.onConnect = handler
}
