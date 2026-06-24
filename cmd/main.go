package main

import (
	"go-darts-tgbot/internal/bot"
	"go-darts-tgbot/internal/config"
	database "go-darts-tgbot/internal/db"
	web "go-darts-tgbot/internal/http"
	"go-darts-tgbot/internal/logger"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	homedir, err := os.UserHomeDir()

	if err != nil {
		slog.Error("cannot to get homedir", slog.Any("err", err))
		panic("fatal error")
	}

	loggerConfig := logger.Config{
		LogDir:     homedir + "/.local/state/darts-tgbot",
		LogFile:    "log",
		MaxSize:    50,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   true,
	}

	cleanup := logger.Init(loggerConfig)
	defer cleanup()
	config.LoadConfig()

	d := database.InitDatabase()
	if err := database.Connect(d); err != nil {
		panic("fatal error")
	}

	g := bot.InitGame()
	bot.InitBot(g, d)

	web.HttpLocal()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	slog.Info("shutdown programm")
}
