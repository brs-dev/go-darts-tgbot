package bot

import (
	"log/slog"
	"time"

	tele "gopkg.in/telebot.v4"
)

func (g *Game) ResetIdleTimer(c tele.Context) {
	if g.IdleTimer != nil {
		g.IdleTimer.Stop()
	}

	g.IdleTimer = time.AfterFunc(120*time.Second, func() {
		slog.Info("exit game by global timeout")

		g.Exit(c)
	})
}

func (g *Game) StopIdleTimer() {
	if g.IdleTimer != nil {
		g.IdleTimer.Stop()
		g.IdleTimer = nil
	}
}
