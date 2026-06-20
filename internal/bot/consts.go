package bot

const (
	MaxSteps     = 3
	AssetsScore1 = "assets/score_1.mp4"
	AssetsScore2 = "assets/score_2.mp4"
	AssetsScore3 = "assets/score_3.mp4"
	AssetsScore4 = "assets/score_4.mp4"
	AssetsScore5 = "assets/score_5.mp4"
	AssetsScore6 = "assets/score_6.mp4"
)

type Messages struct {
	BtnExitCommon       string
	ExitDesc            string
	BtnPlayWithBot      string
	BtnPlayWithFriend   string
	BtnShowStats        string
	GameDesc            string
	BtnStartGameMsg     string
	BtnBackMain         string
	StartGameBotDesc    string
	StartGameFriendDesc string
	BtnExitGame         string
	StatsDesc           string
	GameBotStarted      string
	ThrowResultDesc     string
	BtnMakeThrow        string
	BtnMakeThrow1       string
	BtnMakeThrow2       string
	BtnMakeThrow3       string
	BtnMakeThrow4       string
	BtnContinue         string
	GameResult          string
	FriendshipWon       string
	BtnMainMenu         string
}

var Msg = Messages{
	BtnExitCommon:       `🚪 Я передумал!`,
	ExitDesc:            `Смена не закончится никогда... Никогда не закончится смена! Будут меняться люди, лица, поколения... А смена будет продолжаться!`,
	BtnPlayWithBot:      `🤖 Сыграть с ботом`,
	BtnPlayWithFriend:   `👫 Сыграть с другом`,
	BtnShowStats:        `🏆 Посмотреть статистику`,
	GameDesc:            `Выбери режим игры или <s>иди нахуй</s> посмотри таблицу чемпионов.`,
	BtnStartGameMsg:     `🏁 Начать игру`,
	BtnBackMain:         `🔙 Назад`,
	StartGameBotDesc:    `Способен ли ты противостоять бездушной машине?`,
	StartGameFriendDesc: `Вижу ты решил поиграть со своим <em>дружком</em>?`,
	BtnExitGame:         `🚪 Выйти из игры`,
	StatsDesc:           "Таблица наших <b>чемпионов</b>!",
	GameBotStarted:      `Готовь своё очко! <b>Я начинаю...</b>`,
	ThrowResultDesc:     `<b>%s</b> выбил %d %s!`,
	BtnMakeThrow:        `🎯 Сделать бросок`,
	BtnMakeThrow1:       `💪 Стараться изо всех сил`,
	BtnMakeThrow2:       `🙏 Помолиться перед броском`,
	BtnMakeThrow3:       `🥃 Выпить перед броском`,
	BtnMakeThrow4:       `😎 Выполнить бросок на расслабоне`,
	BtnContinue:         `➡️ Продолжить`,
	GameResult:          "🏆 <b>%s ПОБЕДИЛ!</b>\n%s: %d %s\n%s: %d %s",
	FriendshipWon:       "🏆 <b>ПОБЕДИЛА ДРУЖБА!</b>\n%s: %d %s\n%s: %d %s",
	BtnMainMenu:         `🏠 Перейти в главное меню`,
}
