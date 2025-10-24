package service

import (
	"log"
	"time"
)

func (s *ScheduleService) StartScheduler() {
	go func() {
		for {
			now := time.Now()
			next := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 5, 0, now.Location())
			sleepTime := next.Sub(now)

			log.Println("⏳ Обновление кеша расписаний...")
			if err := s.LoadDailyCache(); err != nil {
				log.Printf("Ошибка обновления кеша: %v\n", err)
			} else {
				log.Println("✅ Кеш расписаний обновлён на текущий день.")
			}

			time.Sleep(sleepTime)
		}
	}()
}
