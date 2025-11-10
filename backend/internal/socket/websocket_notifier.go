package socket

//
//import (
//	"encoding/json"
//	"github.com/TryHanger/digital_signage/backend/internal/cache"
//	"github.com/TryHanger/digital_signage/backend/internal/model"
//	"github.com/TryHanger/digital_signage/backend/internal/repository"
//	"log"
//	"sync"
//
//	"github.com/gorilla/websocket"
//)
//
//type WebSocketNotifier struct {
//	mu            sync.RWMutex
//	connections   map[uint]*websocket.Conn
//	monitorRepo   *repository.MonitorRepository
//	scheduleCache *cache.ScheduleCache
//	onConnect     func(monitorID uint)
//}
//
//func NewWebSocketNotifier(monitorRepo *repository.MonitorRepository, cache *cache.ScheduleCache) *WebSocketNotifier {
//	return &WebSocketNotifier{
//		connections:   make(map[uint]*websocket.Conn),
//		monitorRepo:   monitorRepo,
//		scheduleCache: cache,
//	}
//}
//
//func (n *WebSocketNotifier) HandleConnection(conn *websocket.Conn) {
//	defer conn.Close()
//
//	for {
//		_, msg, err := conn.ReadMessage()
//		if err != nil {
//			log.Printf("⚠️ Ошибка чтения: %v", err)
//			return
//		}
//
//		var payload map[string]interface{}
//		if err := json.Unmarshal(msg, &payload); err != nil {
//			log.Printf("❌ Ошибка парсинга JSON: %v", err)
//			continue
//		}
//
//		event, ok := payload["event"].(string)
//		if !ok {
//			log.Println("⚠️ Неверный формат события")
//			continue
//		}
//
//		data, ok := payload["data"].(map[string]interface{})
//		if !ok {
//			log.Println("⚠️ Неверный формат данных")
//			continue
//		}
//
//		switch event {
//		case "register_monitor":
//			n.handleRegisterByToken(conn, data)
//
//		default:
//			log.Printf("📩 Неизвестное событие: %s", event)
//		}
//	}
//}
//
//func (n *WebSocketNotifier) handleRegisterByToken(conn *websocket.Conn, data map[string]interface{}) {
//	token, ok := data["token"].(string)
//	if !ok || token == "" {
//		log.Println("❌ Токен монитора не передан")
//		conn.WriteJSON(map[string]interface{}{
//			"event": "error",
//			"data":  "token_required",
//		})
//		conn.Close()
//		return
//	}
//
//	monitor, err := n.monitorRepo.GetByToken(token)
//	if err != nil {
//		log.Printf("❌ Монитор с токеном %s не найден: %v", token, err)
//		conn.WriteJSON(map[string]interface{}{
//			"event": "error",
//			"data":  "invalid_token",
//		})
//		conn.Close()
//		return
//	}
//
//	// 🔒 Защита от одновременного доступа
//	n.mu.Lock()
//	if _, exists := n.connections[monitor.ID]; exists {
//		n.mu.Unlock()
//		log.Printf("⚠️ Монитор %s (ID: %d) уже подключён, закрываем новое соединение", monitor.Name, monitor.ID)
//		conn.WriteJSON(map[string]interface{}{
//			"event": "error",
//			"data":  "monitor_already_connected",
//		})
//		conn.Close()
//		return
//	}
//
//	n.connections[monitor.ID] = conn
//	n.mu.Unlock()
//
//	log.Printf("🖥️ Монитор подключён: %s (ID: %d, Token: %s)", monitor.Name, monitor.ID, token)
//
//	// 🧹 Гарантируем очистку соединения при выходе
//	defer func() {
//		n.mu.Lock()
//		delete(n.connections, monitor.ID)
//		n.mu.Unlock()
//		log.Printf("🔌 Монитор %s (ID: %d) отключён", monitor.Name, monitor.ID)
//	}()
//
//	if n.onConnect != nil {
//		n.onConnect(monitor.ID)
//	}
//
//	// Отправляем актуальные данные монитору
//	schedules := n.scheduleCache.GetByMonitorID(monitor.ID)
//	conn.WriteJSON(map[string]interface{}{
//		"event": "init_schedules",
//		"data":  schedules,
//	})
//
//	// ⏳ Читаем сообщения до закрытия (чтобы defer сработал)
//	for {
//		_, _, err := conn.ReadMessage()
//		if err != nil {
//			log.Printf("⚠️ Соединение с монитором %d закрыто: %v", monitor.ID, err)
//			break
//		}
//	}
//}
//
//func (n *WebSocketNotifier) NotifyScheduleUpdate(monitorID uint, schedule model.Schedule) {
//	n.mu.RLock()
//	conn, ok := n.connections[monitorID]
//	n.mu.RUnlock()
//
//	if !ok {
//		log.Printf("⚠️ Монитор %d не подключён", monitorID)
//		return
//	}
//
//	msg := map[string]interface{}{
//		"event": "schedule_update",
//		"data":  schedule,
//	}
//
//	if err := conn.WriteJSON(msg); err != nil {
//		log.Printf("❌ Ошибка отправки данных монитору %d: %v", monitorID, err)
//	}
//
//	log.Printf("📤 Отправка обновления расписания монитору %d", monitorID)
//
//}
//
//func (n *WebSocketNotifier) OnConnect(handler func(monitorID uint)) {
//	n.onConnect = handler
//}
