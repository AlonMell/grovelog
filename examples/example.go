package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/fatih/color"
)

func main() {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}))

	s := color.BlueString("Hello, World!")
	fmt.Println(s)
	logger.Info(s)
	logger.Debug("Debug message")

}
