package main

import (
	"flag"
	"log"
	"math/rand"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var r *rand.Rand

type FortuneBot struct {
	bot                  *tgbotapi.BotAPI
	fortuneTellerNames   []string
	fortuneTellerAnswers []string
}

func New() (*FortuneBot, error) {
	bot, err := tgbotapi.NewBotAPI(mustToken())
	if err != nil {
		return nil, err
	}

	return &FortuneBot{
		bot:                bot,
		fortuneTellerNames: []string{"пушистый нострадамус", "пушистик"},
		fortuneTellerAnswers: []string{
			"Да",
			"Нет",
			"Мурмяу! Звёзды мяукают — да, несомненно!",
			"Хвостик дрожит, усы шевелятся — нетушки!",
			"Пушистое пророчество гласит: да-да-да!",
			"Я свернулся калачиком сомнения... Это — нет.",
			"Вижу в миске отражение твоего успеха. Ответ — да!",
			"Сегодня мои лапки не чувствуют уверенности... нет.",
			"Прыжок веры с подоконника удался — значит да!",
			"Шёрстка встала дыбом... лучше не стоит. Нет.",
			"Я оставил предсказание в лотке. Там написано: нет.",
			"Мурлыканье души говорит — да, но осторожно!",
			"Призрачные мышки в тумане шепчут: нет!",
			"Судьба кувыркается, как клубок — да!",
			"Знак был в облаках, когда я смотрел сквозь занавеску: да.",
			"Мурррр... Путь освещён солнечным лучом. Это да.",
			"Наблюдал за пылью в солнечном луче — не время. Нет.",
			"Провёл ритуал с клубком — он развернулся в 'да'!",
			"Символы на обоях говорят — увы, нет.",
			"Покрутил ушками, повертел лапками… нет.",
			"Хвост уложился в спираль судьбы — да!",
			"Только что видел это в сне о рыбках — однозначно да!",
		},
	}, nil
}

func (f *FortuneBot) SendMessage(chatID int64, msg string) {
	msgConfig := tgbotapi.NewMessage(chatID, msg)
	f.bot.Send(msgConfig)
}

func (f *FortuneBot) IsMessageForFortuneTeller(update *tgbotapi.Update) bool {
	if update.Message == nil || update.Message.Text == "" {
		return false
	}

	msgInLowerCase := strings.ToLower(update.Message.Text)
	for _, name := range f.fortuneTellerNames {
		if strings.Contains(msgInLowerCase, name) {
			return true
		}
	}
	return false
}

func (f *FortuneBot) GetFortuneTellerAnswer() string {
	index := r.Intn(len(f.fortuneTellerAnswers))
	return f.fortuneTellerAnswers[index]
}

func (f *FortuneBot) SendAnswer(chatID int64, update *tgbotapi.Update) {
	msgConfig := tgbotapi.NewMessage(chatID, f.GetFortuneTellerAnswer())
	msgConfig.ReplyToMessageID = update.Message.MessageID
	f.bot.Send(msgConfig)
}

func (f *FortuneBot) HandleUpdate(update tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	chatID := update.Message.Chat.ID

	if update.Message.Text == "/start" {
		f.SendMessage(chatID, "Задай свой вопрос, назвав меня по имени. Ответом на вопрос должны быть либо \"Да\", либо \"Нет\". Например, \"Пушистый нострадамус, съесть ли мне еще одну пачку чипсов?\""+
			"или \"Пушистик, готовиться ли мне к завтрашнему экзамену?\"")
		return
	}

	if f.IsMessageForFortuneTeller(&update) {
		f.SendAnswer(chatID, &update)
	}
}

func (f *FortuneBot) Run() {
	updateConfig := tgbotapi.NewUpdate(0)
	updates := f.bot.GetUpdatesChan(updateConfig)

	const workerCount = 5
	jobs := make(chan tgbotapi.Update, 100)

	for i := 0; i < workerCount; i++ {
		go func() {
			for update := range jobs {
				f.HandleUpdate(update)
			}
		}()
	}

	for update := range updates {
		jobs <- update
	}
}

func mustToken() string {
	token := flag.String(
		"tg-bot-token",
		"",
		"token for access to telegram bot",
	)

	flag.Parse()

	if *token == "" {
		log.Fatal("token is not specified")
	}

	return *token
}

func main() {
	f, err := New()
	if err != nil {
		log.Fatalf("error connecting to Telegram: %s", err.Error())
	}
	log.Print("service is running")
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
	f.Run()
}
