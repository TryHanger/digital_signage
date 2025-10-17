package utils

import (
	"fmt"
	"log"
	"time"

	"github.com/TryHanger/digital_signage/internal/service"
)

type SocketServer interface {
	BroadcastToRoom(namespace, room, event string, message ...interface{}) bool
}

type Scheduler struct {
	ScheduleService *service.ScheduleService
	SocketServer    SocketServer
	Interval        time.Duration
}

func NewScheduler(svc *service.ScheduleService, socketServer SocketServer, interval time.Duration) *Scheduler {
	return &Scheduler{
		ScheduleService: svc,
		SocketServer:    socketServer,
		Interval:        interval,
	}
}

func (s *Scheduler) Start() {
	// карта внутри Scheduler
	lastContentSent := make(map[uint]uint) // monitorID -> contentID

	ticker := time.NewTicker(s.Interval)
	go func() {
		for range ticker.C {
			active, err := s.ScheduleService.GetAllActive()
			if err != nil {
				log.Println("scheduler error:", err)
				continue
			}

			now := time.Now()
			for _, sch := range active {
				// если контент закончился, убираем из lastContentSent
				if now.After(sch.EndTime) {
					if lastContentSent[sch.MonitorID] == sch.Content.ID {
						delete(lastContentSent, sch.MonitorID)
					}
					continue
				}

				// если контент уже отправлен — пропускаем
				if lastContentSent[sch.MonitorID] == sch.Content.ID {
					continue
				}

				// отправляем контент на монитор
				log.Printf("🎬 Отправка контента на монитор %d: %s", sch.MonitorID, sch.Content.Title)
				s.SocketServer.BroadcastToRoom(
					"/",
					fmt.Sprintf("monitor_%d", sch.MonitorID),
					"show",
					map[string]interface{}{
						"title":    sch.Content.Title,
						"type":     sch.Content.Type,
						"url":      sch.Content.URL,
						"duration": sch.Content.Duration,
					},
				)

				// отмечаем, что контент отправлен
				lastContentSent[sch.MonitorID] = sch.Content.ID
			}
		}
	}()
}
