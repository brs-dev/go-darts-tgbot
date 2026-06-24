package bot

import (
	"log/slog"
	"sync"
	"time"

	tele "gopkg.in/telebot.v4"
)

var startGameVideoID string

type Score struct {
	PlayerID    int64
	PlayerValue int
	BotValue    int
}

type Game struct {
	mu            sync.RWMutex
	GameStarted   bool
	Player1ID     int64
	Player2ID     int64
	IsMultiplayer bool
	Queue         StepState
	Score         Score
	GameChan      chan string
	StopChan      chan struct{}
	IdleTimer     *time.Timer
	LastMessage   *tele.Message
	ChatID        int64
}

type Step struct {
	Player1ID int64
	Player2ID int64
}

type StepState struct {
	Steps         []Step
	CurrentPlayer int64
}

func InitGame() *Game {
	return &Game{}
}

func (g *Game) StartGame(Player1ID int64, Player2ID int64) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.CreateQueue(Player2ID, Player1ID)

	g.GameStarted = true

	g.Player1ID = Player1ID
	g.Player2ID = Player2ID

	g.Score.PlayerID = Player1ID
	g.Score.PlayerValue = 0
	g.Score.BotValue = 0
}

func (g *Game) EndGame() {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.GameStarted = false
	g.ChatID = 0
	g.Player1ID = 0
	g.Player2ID = 0
}

func (g *Game) CreateQueue(player1 int64, player2 int64) {
	slog.Info("create queue")
	steps := make([]Step, 0, MaxSteps)

	for range MaxSteps {
		steps = append(steps, Step{
			Player1ID: player1,
			Player2ID: player2,
		})
	}

	g.Queue = StepState{
		Steps:         steps,
		CurrentPlayer: player1,
	}
}

func (g *Game) ResetQueue() {
	slog.Info("reset queue")
	g.Queue = StepState{}
}

func (g *Game) RemovePrevMsg(c tele.Context) {
	if g.LastMessage == nil {
		slog.Warn("no previous message for delete")
		return
	}

	if err := c.Bot().Delete(g.LastMessage); err != nil {
		slog.Warn("cannot to remove message", slog.Any("err", err))
	}
}

func (g *Game) GetUserName(c tele.Context, playerID int64) string {
	player, err := c.Bot().ChatByID(playerID)

	if err != nil {
		slog.Error("failed to get player struct", slog.Any("err", err))
		return "N/A"
	}

	return player.FirstName + " " + player.LastName
}
