package bot

import (
	"log/slog"

	tele "gopkg.in/telebot.v4"
)

func (h *Handler) NextTurn(c tele.Context) error {
	h.mu.Lock()
	h.GameChan <- "throw"
	h.mu.Unlock()

	ContinueGame <- 1

	return nil
}

func (h *Handler) ThrowHandler(c tele.Context) (*tele.Message, error) {
	msg, err := c.Bot().Send(c.Chat(), &tele.Dice{
		Type: tele.Dart.Type,
	})

	if err != nil {
		slog.Error("failed to send dart", slog.Any("err", err))
		return nil, h.Exit(c)
	}

	return msg, nil
}
