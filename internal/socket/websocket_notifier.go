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

		event, ok := payload["event"].(string)
		if !ok {
			log.Println("⚠️ Неверный формат события")
			continue
		}

		data, ok := payload["data"].(map[string]interface{})
		if !ok {
			log.Println("⚠️ Неверный формат данных")
			continue
		}

		switch event {
		case "register_monitor":
			n.handleRegisterByToken(conn, data)

		default:
			log.Printf("📩 Неизвестное событие: %s", event)
		}
	}
}

func (n *WebSocketNotifier) handleRegisterByToken(conn *websocket.Conn, data map[string]interface{}) {
	token, ok := data["token"].(string)
	if !ok || token == "" {
		log.Println("❌ Токен монитора не передан")
		return
	}

	monitor, err := n.monitorRepo.GetByToken(token)
	if err != nil {
		log.Printf("❌ Монитор с токеном %s не найден: %v", token, err)
		return
	}

	n.mu.Lock()
	n.connections[monitor.ID] = conn
	n.mu.Unlock()

	log.Printf("🖥️ Монитор подключён: %s (ID: %d, Token: %s)", monitor.Name, monitor.ID, token)

	if n.onConnect != nil {
		n.onConnect(monitor.ID)
	}

	// Отправляем актуальные данные монитору
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

	log.Printf("📤 Отправка обновления расписания монитору %d", monitorID)

}

func (n *WebSocketNotifier) OnConnect(handler func(monitorID uint)) {
	n.onConnect = handler
}
