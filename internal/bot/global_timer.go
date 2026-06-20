package bot

import (
	"log/slog"
	"time"

	tele "gopkg.in/telebot.v4"
)

func (h *Handler) ResetIdleTimer(c tele.Context) {
	if h.IdleTimer != nil {
		h.IdleTimer.Stop()
	}

	h.IdleTimer = time.AfterFunc(120*time.Second, func() {
		slog.Info("exit game by global timeout")

		h.Exit(c)
	})
}

func (h *Handler) StopIdleTimer() {
	if h.IdleTimer != nil {
		h.IdleTimer.Stop()
		h.IdleTimer = nil
	}
}
