package bot

import (
	"context"
	"fmt"
	database "go-darts-tgbot/internal/db"
	"log/slog"

	tele "gopkg.in/telebot.v4"
)

type CreateUsersParams struct {
	UserId    int64
	FirstName string
	LastName  *string
	Username  *string
}

func handlerCreateUser(db *database.Database, params CreateUsersParams) {
	ctx := context.Background()

	user := database.User{
		UserID:    params.UserId,
		FirstName: params.FirstName,
		LastName:  params.LastName,
		Username:  params.Username,
	}

	result := db.DB.WithContext(ctx).
		Omit("ID", "CreatedAt", "UpdatedAt").
		Where("user_id = ?", params.UserId).
		FirstOrCreate(&user)

	if result.Error != nil {
		slog.Error("error to create new user", slog.Any("err", result.Error))
		return
	}

	if result.RowsAffected > 0 {
		slog.Info(fmt.Sprintf("new user with ID %d added", params.UserId))
	} else {
		slog.Info(fmt.Sprintf("user with ID %d already exists", params.UserId))
	}
}

func strToPtr(s string) *string {
	if s == "" {
		return nil
	}

	return &s
}

func getUserStruct(c tele.Context) CreateUsersParams {
	sender := c.Sender()

	return CreateUsersParams{
		UserId:    sender.ID,
		FirstName: sender.FirstName,
		LastName:  strToPtr(sender.LastName),
		Username:  strToPtr(sender.Username),
	}
}

func (g *Game) Game(c tele.Context) error {
	slog.Info("main menu called")
	slog.Info(fmt.Sprintf("main menu started by %d", c.Sender().ID))

	botId := c.Bot().(*tele.Bot).Me.ID
	g.StartGame(c.Sender().ID, botId)
	g.RemovePrevMsg(c)

	userParams := getUserStruct(c)
	handlerCreateUser(DatabaseInstance, userParams)

	menu := &tele.ReplyMarkup{}
	menu.Inline(
		menu.Row(
			menu.Data(Msg.BtnPlayWithBot, "btn_play_with_bot"),
		),
		// menu.Row(
		// 	menu.Data(Msg.BtnPlayWithFriend, "btn_play_with_friend"),
		// ),
		menu.Row(
			menu.Data(Msg.BtnShowStats, "btn_show_stats"),
		),
		menu.Row(
			menu.Data(Msg.BtnExitCommon, "btn_exit"),
		),
	)

	if startGameVideoID != "" {
		video := &tele.Video{
			File:    tele.File{FileID: startGameVideoID},
			Caption: Msg.GameDesc,
		}

		msg, err := c.Bot().Send(c.Chat(), video, menu, tele.ModeHTML)

		if err == nil {
			g.LastMessage = msg
		}

		return err
	}

	video := &tele.Video{
		File:    tele.FromDisk("assets/start_game.mp4"),
		Caption: Msg.GameDesc,
	}

	msg, err := c.Bot().Send(c.Chat(), video, menu, tele.ModeHTML)

	if err != nil {
		slog.Warn("cannot to send video", slog.Any("err", err))
		g.Exit(c)
		return err
	} else {
		g.LastMessage = msg
	}

	if msg.Video != nil {
		startGameVideoID = msg.Video.FileID
	}

	return nil
}

func (g *Game) Exit(c tele.Context) error {
	if err := c.Respond(&tele.CallbackResponse{
		Text: Msg.ExitDesc,
	}); err != nil {
		slog.Warn("cannot to send response", slog.Any("err", err))
	}

	if g.StopChan != nil {
		close(g.StopChan)
		g.StopChan = nil
	}

	if g.GameChan != nil {
		close(g.GameChan)
		g.GameChan = nil
	}

	g.RemovePrevMsg(c)
	g.ResetQueue()
	g.EndGame()
	g.StopIdleTimer()

	slog.Info("exit game")
	return nil
}
