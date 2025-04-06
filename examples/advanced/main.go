package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"time"

	"github.com/AlonMell/grovelog/grovelog"
	"github.com/AlonMell/grovelog/grovelog/helper"
)

func main() {
	// Создаем логгер с пользовательскими настройками
	opts := grovelog.DefaultOptions()
	opts.Level = slog.LevelDebug
	opts.TimeFormat = "2006-01-02 15:04:05.000"
	opts.Format = grovelog.ColorFormat

	log := grovelog.New(opts)
	log.Info("Запуск приложения с пользовательскими настройками")

	// Логгер, который пишет и в консоль, и в файл
	fileLogger, closer, err := grovelog.NewWithFile("app.log", grovelog.ProductionOptions())
	if err != nil {
		log.Error("Не удалось создать файловый логгер", helper.Err(err))
		return
	}
	defer closer.Close()

	// Логируем и в консоль, и в файл
	fileLogger.Info("Приложение запущено", "pid", os.Getpid())

	// Демонстрация использования контекста
	ctx := context.Background()
	ctx = helper.ContextWithLogger(ctx, log.Logger)

	// Запускаем сервис с контекстом
	if err := runService(ctx); err != nil {
		log.Error("Ошибка в сервисе", helper.Err(err))
	}
}

func runService(ctx context.Context) error {
	// Получаем логгер из контекста
	log := helper.WithContext(ctx)

	// Используем логгер из контекста
	log.Info("Сервис запущен")

	// Имитируем работу
	for i := range 3 {
		log.Debug("Обработка итерации", "iteration", i)

		// Имитируем периодические ошибки
		if i == 1 {
			err := errors.New("временная ошибка")
			log.Warn("Обнаружена проблема", helper.Err(err), helper.Caller(0))
		}

		time.Sleep(100 * time.Millisecond)
	}

	// Демонстрация логирования с детальной информацией
	userLog := log.With(
		"user_id", 42,
		"ip", "192.168.1.100",
		"user_agent", "Mozilla/5.0...",
	)

	userLog.Info("Действие пользователя",
		"action", "purchase",
		"item_id", "product-123",
		"amount", 29.99,
	)

	return nil
}
