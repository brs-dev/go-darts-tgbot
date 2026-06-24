package bot

import (
	"log/slog"

	tele "gopkg.in/telebot.v4"
)

func (g *Game) NextTurn(c tele.Context) error {
	g.mu.Lock()
	g.GameChan <- "throw"
	g.mu.Unlock()

	ContinueGame <- 1

	return nil
}

func (g *Game) ThrowHandler(c tele.Context) (*tele.Message, error) {
	msg, err := c.Bot().Send(c.Chat(), &tele.Dice{
		Type: tele.Dart.Type,
	})

	if err != nil {
		slog.Error("failed to send dart", slog.Any("err", err))
		return nil, g.Exit(c)
	}

	return msg, nil
}
