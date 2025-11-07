package utils

import (
	"github.com/TryHanger/digital_signage/internal/model"
	"time"
)

func IsActiveToday(sched model.Schedule, today time.Time) bool {
	// 1️⃣ Проверяем, что расписание вообще в диапазоне дат
	if today.Before(sched.DateStart.Truncate(24*time.Hour)) ||
		today.After(sched.DateEnd.Truncate(24*time.Hour)) {
		return false
	}

	// 2️⃣ Проверяем исключения (если есть)
	for _, ex := range sched.Exceptions {
		if ex.Date.Truncate(24 * time.Hour).Equal(today) {
			return false
		}
	}

	// 3️⃣ Проверяем тип повторения
	switch sched.RepeatPattern {
	case "none":
		// Просто однократное расписание
		return today.Equal(sched.DateStart.Truncate(24 * time.Hour))

	case "daily":
		return true // каждый день в диапазоне

	case "weekly":
		weekday := int(today.Weekday()) // 0 = Sunday, 1 = Monday...
		for _, d := range sched.DaysOfWeek {
			if int(d) == weekday {
				return true
			}
		}
		return false

	case "monthly":
		return sched.DateStart.Day() == today.Day()

	default:
		return false
	}
}
