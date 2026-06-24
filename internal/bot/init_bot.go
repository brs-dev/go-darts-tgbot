package bot

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"gopkg.in/telebot.v4/middleware"

	cfg "go-darts-tgbot/internal/config"
	database "go-darts-tgbot/internal/db"

	tele "gopkg.in/telebot.v4"
)

var (
	Bot              *tele.Bot
	DatabaseInstance *database.Database
)

func registerHandler(g *Game) {
	gameGroup := Bot.Group()

	gameGroup.Use(func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			g.ResetIdleTimer(c)

			return next(c)
		}
	})

	gameGroup.Use(func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			if g.GameStarted && g.Player1ID != c.Sender().ID {
				slog.Debug(fmt.Sprintf(
					"interrupting the user's action: %d; in chat %d",
					c.Sender().ID,
					c.Chat().ID,
				))
				return nil
			}

			return next(c)
		}
	})

	gameGroup.Handle("/game", g.Game)
	gameGroup.Handle("\fbtn_exit", g.Exit)
	gameGroup.Handle("\fbtn_back_main", g.Game)
	gameGroup.Handle("\fbtn_play_with_bot", g.StartGameBot)
	gameGroup.Handle("\fbtn_show_stats", g.Stats)
	gameGroup.Handle("\fbtn_start_game_bot", g.GameBotStarted)
	gameGroup.Handle("\fbtn_make_throw", g.NextTurn)
}

func InitBot(g *Game, databaseInstance *database.Database) {
	go func() {
		DatabaseInstance = databaseInstance

		botConfig := tele.Settings{
			Token:  cfg.GlobalConfig.BotToken,
			Poller: &tele.LongPoller{Timeout: 10 * time.Second},
			OnError: func(err error, c tele.Context) {
				slog.Warn("bot error", slog.Any("err", err))
			},
		}

		for {
			var err error

			Bot, err = tele.NewBot(botConfig)

			if err != nil {
				slog.Warn("cannot to build new bot", slog.Any("err", err))
				time.Sleep(5 * time.Second)
				continue
			}

			Bot.Use(middleware.AutoRespond())
			Bot.Use(func(next tele.HandlerFunc) tele.HandlerFunc {
				return func(c tele.Context) error {
					if c.Chat().Type == tele.ChatPrivate {
						return nil
					}

					if strings.HasPrefix(c.Text(), "/game") && g.ChatID == 0 {
						g.ChatID = c.Chat().ID
					} else if strings.HasPrefix(c.Text(), "/game") && g.ChatID != 0 {
						return nil
					}

					return next(c)
				}
			})

			registerHandler(g)

			slog.Info("bot is ready")
			Bot.Start()

			slog.Warn("lost connection with Telegram; trying to restart polling in 5 seconds")
			time.Sleep(5 * time.Second)
		}
	}()
}
