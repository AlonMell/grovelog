package grovelog

import (
	"io"
	"log/slog"
	"os"
	"time"
)

// LogFormat определяет формат вывода логов
type LogFormat int

const (
	// JSONFormat - формат JSON для анализа логов
	JSONFormat LogFormat = iota
	// TextFormat - текстовый формат для читаемости
	TextFormat
	// ColorFormat - цветной текстовый формат для разработки
	ColorFormat
)

// Options конфигурирует поведение логгера
type Options struct {
	// Level - минимальный уровень логирования
	Level slog.Level
	// TimeFormat - формат времени для логов (используя формат time.Format)
	TimeFormat string
	// AddSource - добавляет информацию о местоположении в коде
	AddSource bool
	// Format - формат вывода логов
	Format LogFormat
	// Output - место назначения для логов
	Output io.Writer
	// AddCaller - добавляет информацию о вызывающей функции
	AddCaller bool
}

// DefaultOptions возвращает опции логгера по умолчанию
func DefaultOptions() Options {
	return Options{
		Level:      slog.LevelInfo,
		TimeFormat: time.RFC3339,
		AddSource:  false,
		Format:     JSONFormat,
		Output:     os.Stdout,
		AddCaller:  false,
	}
}

// DevelopmentOptions возвращает опции для разработки
func DevelopmentOptions() Options {
	return Options{
		Level:      slog.LevelDebug,
		TimeFormat: "15:04:05.000",
		AddSource:  true,
		Format:     ColorFormat,
		Output:     os.Stdout,
		AddCaller:  true,
	}
}

// ProductionOptions возвращает опции для продакшена
func ProductionOptions() Options {
	return Options{
		Level:      slog.LevelInfo,
		TimeFormat: time.RFC3339,
		AddSource:  false,
		Format:     JSONFormat,
		Output:     os.Stdout,
		AddCaller:  false,
	}
}
