package bot

import (
	"context"
	"fmt"

	u "go-darts-tgbot/internal/utils"
	tele "gopkg.in/telebot.v4"
)

func (h *Handler) Stats(c tele.Context) error {
	h.RemovePrevMsg(c)

	users, err := DatabaseInstance.GetAllUsers(context.Background())

	if err != nil {
		h.Exit(c)
		return err
	}

	caption := Msg.StatsDesc + "\n\n"

	for i, user := range users {
		lastName := ""
		if user.LastName != nil {
			lastName = *user.LastName
		}

		name := fmt.Sprintf("%s %s:", user.FirstName, lastName)
		plurScore := u.PluralizePoints(user.Score)

		caption += fmt.Sprintf("%d.) <b>%s</b> %d %s\n", i+1, name, user.Score, plurScore)
	}

	menu := &tele.ReplyMarkup{}
	menu.Inline(
		menu.Row(
			menu.Data(Msg.BtnExitGame, "btn_exit"),
		),
		menu.Row(
			menu.Data(Msg.BtnBackMain, "btn_back_main"),
		),
	)

	photo := &tele.Photo{
		File:    tele.FromDisk("assets/stats.png"),
		Caption: caption,
	}

	msg, err := c.Bot().Send(c.Chat(), photo, menu, tele.ModeHTML)

	if err == nil {
		h.LastMessage = msg
	}

	return err
}
