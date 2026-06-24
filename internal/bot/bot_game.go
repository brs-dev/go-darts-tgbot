package bot

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	database "go-darts-tgbot/internal/db"
	u "go-darts-tgbot/internal/utils"

	tele "gopkg.in/telebot.v4"
)

var (
	ContinueGame = make(chan int, 1)
)

func (g *Game) StartGameBot(c tele.Context) error {
	g.RemovePrevMsg(c)

	menu := &tele.ReplyMarkup{}
	menu.Inline(
		menu.Row(
			menu.Data(Msg.BtnStartGameMsg, "btn_start_game_bot"),
		),
		menu.Row(
			menu.Data(Msg.BtnExitCommon, "btn_exit"),
		),
		menu.Row(
			menu.Data(Msg.BtnBackMain, "btn_back_main"),
		),
	)

	photo := &tele.Photo{
		File:    tele.FromDisk("assets/start_game_bot.jpg"),
		Caption: Msg.StartGameBotDesc,
	}

	msg, err := c.Bot().Send(c.Chat(), photo, menu, tele.ModeHTML)

	if err == nil {
		g.LastMessage = msg
	}

	return err
}

func (g *Game) GameBotLoop(c tele.Context, stop <-chan struct{}) {
	slog.Info("game loop starts here")

	g.mu.RLock()
	steps := make([]Step, len(g.Queue.Steps))
	copy(steps, g.Queue.Steps)
	g.mu.RUnlock()

	g.mu.Lock()
	g.GameChan <- "throw"
	g.mu.Unlock()

	for i := range steps {
		select {
		case <-stop:
			slog.Info("game loop stopped")
			return
		default:
		}

		pair := steps[i]
		players := []int64{pair.Player1ID, pair.Player2ID}

		for _, playerID := range players {
			isBot, err := isBot(c, playerID)

			if err != nil {
				g.Exit(c)
				return
			}

			g.mu.Lock()
			g.Queue.CurrentPlayer = playerID
			g.mu.Unlock()

			throwMessage := g.betweenThrowMessage(c, isBot)

			select {
			case <-stop:
				return
			case continueGameStatus := <-ContinueGame:
				if continueGameStatus > 0 {
					if throwMessage == nil {
						g.Exit(c)
						return
					}

					g.RemovePrevMsg(c)
				}
			}
		}
	}

	slog.Info("game loop ends here")

	g.GameResults(c)
}

func (g *Game) GameBotStarted(c tele.Context) error {
	g.RemovePrevMsg(c)

	photo := &tele.Photo{
		File:    tele.FromDisk("assets/game_bot_started.jpg"),
		Caption: Msg.GameBotStarted,
	}

	msg, err := c.Bot().Send(c.Chat(), photo, tele.ModeHTML)
	g.LastMessage = msg

	if err != nil {
		slog.Error("cannot to send message", slog.Any("err", err))
		return g.Exit(c)
	}

	time.Sleep(5 * time.Second)

	g.RemovePrevMsg(c)

	g.GameChan = make(chan string, 1)
	slog.Info("game channel created")

	g.StopChan = make(chan struct{})
	go g.GameBotLoop(c, g.StopChan)

	return nil
}

func (g *Game) GameResults(c tele.Context) error {
	slog.Info("game results called")

	var message string
	var filepath string

	playerName1 := g.GetUserName(c, g.Player1ID)
	playerName2 := g.GetUserName(c, g.Player2ID)

	plurPlayerScore := u.PluralizePoints(g.Score.PlayerValue)
	plurBotScore := u.PluralizePoints(g.Score.BotValue)

	if g.Score.PlayerValue > g.Score.BotValue {
		filepath = "assets/player_won.jpg"

		message = fmt.Sprintf(
			Msg.GameResult,
			playerName1,
			playerName1,
			g.Score.PlayerValue,
			plurPlayerScore,
			playerName2,
			g.Score.BotValue,
			plurBotScore,
		)
	} else if g.Score.PlayerValue == g.Score.BotValue {
		filepath = "assets/friendship_won.jpg"

		message = fmt.Sprintf(
			Msg.FriendshipWon,
			playerName1,
			g.Score.PlayerValue,
			plurPlayerScore,
			playerName2,
			g.Score.BotValue,
			plurBotScore,
		)
	} else {
		filepath = "assets/bot_won.jpg"

		message = fmt.Sprintf(
			Msg.GameResult,
			playerName2,
			playerName1,
			g.Score.PlayerValue,
			plurPlayerScore,
			playerName2,
			g.Score.BotValue,
			plurBotScore,
		)
	}

	menu := &tele.ReplyMarkup{}
	menu.Inline(
		menu.Row(menu.Data(Msg.BtnMainMenu, "btn_back_main")),
		menu.Row(menu.Data(Msg.BtnShowStats, "btn_show_stats")),
		menu.Row(menu.Data(Msg.BtnExitGame, "btn_exit")),
	)

	photo := &tele.Photo{
		File:    tele.FromDisk(filepath),
		Caption: message,
	}

	g.ResetQueue()
	g.EndGame()

	msg, err := c.Bot().Send(c.Chat(), photo, menu, tele.ModeHTML)

	if err == nil {
		g.LastMessage = msg
	}

	return err
}

func isBot(c tele.Context, playerID int64) (bool, error) {
	member, err := c.Bot().ChatMemberOf(c.Chat(), &tele.User{ID: playerID})

	if err != nil {
		slog.Error("failed to get telegram user struct", slog.Any("err", err))
		return false, err
	}

	isBot := member.User.IsBot

	return isBot, nil
}

func (g *Game) getVideoStruct(score int, caption string) *tele.Video {
	var video *tele.Video
	var filepath string

	switch score {
	case 1:
		filepath = AssetsScore1
	case 2:
		filepath = AssetsScore2
	case 3:
		filepath = AssetsScore3
	case 4:
		filepath = AssetsScore4
	case 5:
		filepath = AssetsScore5
	case 6:
		filepath = AssetsScore6
	default:
		filepath = AssetsScore1
	}

	video = &tele.Video{
		File:    tele.FromDisk(filepath),
		Caption: caption,
	}

	return video
}

func (g *Game) betweenThrowMessage(c tele.Context, isBot bool) *tele.Message {
	if <-g.GameChan == "throw" {
		msg, _ := g.ThrowHandler(c)
		time.Sleep(7 * time.Second)

		playerName := g.GetUserName(c, g.Queue.CurrentPlayer)
		score := msg.Dice.Value
		pluralized := u.PluralizePoints(score)
		message := fmt.Sprintf(Msg.ThrowResultDesc, playerName, score, pluralized)

		if err := c.Bot().Delete(msg); err != nil {
			slog.Warn("cannot to remove message", slog.Any("err", err))
		}

		if isBot {
			g.Score.BotValue += score

			menu := &tele.ReplyMarkup{}
			menu.Inline(
				menu.Row(menu.Data(Msg.BtnMakeThrow, "btn_make_throw")),
				menu.Row(menu.Data(Msg.BtnMakeThrow1, "btn_make_throw")),
				menu.Row(menu.Data(Msg.BtnMakeThrow2, "btn_make_throw")),
				menu.Row(menu.Data(Msg.BtnMakeThrow3, "btn_make_throw")),
				menu.Row(menu.Data(Msg.BtnMakeThrow4, "btn_make_throw")),
				menu.Row(menu.Data(Msg.BtnExitGame, "btn_exit")),
			)

			msg, err := c.Bot().Send(
				c.Chat(),
				g.getVideoStruct(score, message),
				menu,
				tele.ModeHTML,
			)
			if err != nil {
				slog.Error("cannot to send message", slog.Any("err", err))
				return nil
			}

			slog.Info(
				"making a throw by bot",
				slog.Int("val", g.Score.BotValue),
				slog.Int64("id", g.Queue.CurrentPlayer),
			)

			g.LastMessage = msg

			return msg
		} else {
			g.Score.PlayerValue += score

			DatabaseInstance.IncrementUserScore(database.IncUserScoreParams{
				Ctx:    context.Background(),
				UserID: g.Score.PlayerID,
				Points: score,
			})

			menu := &tele.ReplyMarkup{}
			menu.Inline(
				menu.Row(menu.Data(Msg.BtnContinue, "btn_make_throw")),
				menu.Row(menu.Data(Msg.BtnExitGame, "btn_exit")),
			)

			msg, err := c.Bot().Send(
				c.Chat(),
				g.getVideoStruct(score, message),
				menu,
				tele.ModeHTML,
			)

			if err != nil {
				slog.Error("cannot to send message", slog.Any("err", err))
				return nil
			}

			slog.Info(
				"making a throw by player",
				slog.Int("val", g.Score.PlayerValue),
				slog.Int64("id", g.Queue.CurrentPlayer),
			)

			g.LastMessage = msg

			return msg
		}

	}

	return nil
}
