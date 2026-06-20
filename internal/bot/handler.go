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

type Handler struct {
	mu            sync.RWMutex
	GameStarted   bool
	Player1ID     int64
	Player2ID     int64
	IsMultiplayer bool
	Queue         StepState
	Score         Score
	GameChan      chan string
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

func InitHandler() *Handler {
	return &Handler{}
}

func (h *Handler) StartGame(Player1ID int64, Player2ID int64) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.CreateQueue(Player2ID, Player1ID)

	h.GameStarted = true

	h.Player1ID = Player1ID
	h.Player2ID = Player2ID

	h.Score.PlayerID = Player1ID
	h.Score.PlayerValue = 0
	h.Score.BotValue = 0
}

func (h *Handler) StartMultiplayerGame(Player1ID int64, Player2ID int64) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.GameStarted = true
	h.Player1ID = Player1ID
	h.Player2ID = Player2ID
}

func (h *Handler) EndGame() {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.GameStarted = false
	h.ChatID = 0
	h.Player1ID = 0
	h.Player2ID = 0
}

func (h *Handler) CreateQueue(player1 int64, player2 int64) {
	slog.Info("create queue")
	steps := make([]Step, 0, MaxSteps)

	for range MaxSteps {
		steps = append(steps, Step{
			Player1ID: player1,
			Player2ID: player2,
		})
	}

	h.Queue = StepState{
		Steps:         steps,
		CurrentPlayer: player1,
	}
}

func (h *Handler) ResetQueue() {
	slog.Info("reset queue")
	h.Queue = StepState{}
}

func (h *Handler) RemovePrevMsg(c tele.Context) {
	if h.LastMessage == nil {
		slog.Warn("no previous message for delete")
		return
	}

	if err := c.Bot().Delete(h.LastMessage); err != nil {
		slog.Warn("cannot to remove message", slog.Any("err", err))
	}
}

func (h *Handler) GetUserName(c tele.Context, playerID int64) string {
	player, err := c.Bot().ChatByID(playerID)

	if err != nil {
		slog.Error("failed to get player struct", slog.Any("err", err))
		return "N/A"
	}

	return player.FirstName + " " + player.LastName
}
