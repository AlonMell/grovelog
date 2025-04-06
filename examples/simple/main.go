package main

import (
	"errors"

	"github.com/AlonMell/grovelog/grovelog"
	"github.com/AlonMell/grovelog/grovelog/helper"
)

func main() {
	// Создаем логгер для разработки
	log := grovelog.Development()

	// Логируем некоторые сообщения
	log.Info("Запуск приложения")
	log.Debug("Отладочная информация", "count", 42)

	// Логирование с атрибутами
	log.Info("Пользователь вошел в систему",
		"user_id", 12345,
		"role", "admin",
	)

	// Логирование ошибки
	err := errors.New("что-то пошло не так")
	log.Error("Ошибка обработки запроса",
		helper.Err(err),
		"request_id", "abcd-1234",
	)

	// Создаем логгер с дополнительным контекстом
	userLogger := log.With(
		"user_id", 12345,
		"session_id", "xyz-789",
	)

	userLogger.Info("Пользователь выполнил действие", "action", "logout")

	// Пример с группировкой
	apiLogger := log.WithGroup("api")
	apiLogger.Info("API готов", "version", "1.0")
	apiLogger.Warn("Устаревший метод", "method", "GET /old-endpoint")
}
