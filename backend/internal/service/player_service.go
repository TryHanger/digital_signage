package service

import (
	"github.com/TryHanger/digital_signage/backend/internal/socket"
	"log"
)

// PlayerService отвечает за регистрацию и работу плееров (мониторов)
type PlayerService struct {
	notifier *socket.WebSocketNotifier
}

// NewPlayerService создаёт сервис для управления плеерами
func NewPlayerService(notifier *socket.WebSocketNotifier) *PlayerService {
	return &PlayerService{notifier: notifier}
}

// RegisterHandlers задаёт callback, который вызывается при подключении нового плеера
func (s *PlayerService) RegisterHandlers(handler func(monitorID string)) {
	//s.notifier.RegisterMonitor("register_monitor", handler)
	log.Println("🎧 PlayerService: обработчики зарегистрированы")
}
