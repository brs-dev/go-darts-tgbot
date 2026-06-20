package bot

import (
	tele "gopkg.in/telebot.v4"
)

func (h *Handler) StartGameFriend(c tele.Context) error {
	h.RemovePrevMsg(c)

	menu := &tele.ReplyMarkup{}
	menu.Inline(
		menu.Row(
			menu.Data(Msg.BtnStartGameMsg, "btn_start_game_friend"),
		),
		menu.Row(
			menu.Data(Msg.BtnExitCommon, "btn_exit"),
		),
		menu.Row(
			menu.Data(Msg.BtnBackMain, "btn_back_main"),
		),
	)

	photo := &tele.Photo{
		File:    tele.FromDisk("assets/start_friend_game.jpg"),
		Caption: Msg.StartGameFriendDesc,
	}

	return c.Send(photo, menu, tele.ModeHTML)
}
