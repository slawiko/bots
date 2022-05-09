package main

import (
	"log"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	TriggerKeyword = "як будзе"
	ErrorMessage   = "Нешта чамусьці пайшло ня так. Стварыце калі ласка ішшу на гітхабе https://github.com/slawiko/ru-bel-tg-bot/issues"
	HelpMessage    = `Спосабы ўзаемадзеяння:
<b>У прываце</b>: наўпрост пішыце слова на рускай мове.
<b>У группе</b>: пачніце ваша паведамленне са словаў <code>як будзе</code> і далей слово на русском языке. Напрыклад: <code>як будзе письмо</code>.

Таксама вы можаце не пераходзіць на рускую раскладку і пытацца, напрыклад, слова <code>ўавель</code> ці <code>олівка</code>.

<i>На дадзены момант я ўмею перакладаць толькі па адным слове за раз і не разумею памылкі, прабачце.</i>

© Усе пераклады я бяру з https://skarnik.by, дзякуй яму вялікі.`
	StartMessage = `Прывітаннечка. Мяне клічуць Жэўжык, я дапамагаю перайсьці на родную мову. Вы можаце пытацца ў мяне слова на рускай, а я адкажу вам на беларускай.

Вы можаце дадаць мяне ў группу і пытацца не выходзячы з дыялогу з сябрамі. За дапамогай клацайце /help`
)

var BOT_API_KEY = os.Args[1]

func main() {
	bot, err := tgbotapi.NewBotAPI(BOT_API_KEY)
	if err != nil {
		log.Println(err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			handleCommand(bot, &update)
			continue
		}

		if update.Message.Chat.IsGroup() || update.Message.Chat.IsSuperGroup() {
			handleGroupMessage(bot, &update)
		}

		if update.Message.Chat.IsPrivate() {
			handlePrivateMessage(bot, &update)
		}
	}
}

func sendMsg(bot *tgbotapi.BotAPI, msg tgbotapi.MessageConfig) {
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
		msg.Text = ErrorMessage
		msg.DisableWebPagePreview = true

		_, err := bot.Send(msg)
		if err != nil {
			log.Println(err)
		}
	}
}

func prepareRequestText(dirtyRequestText string) string {
	return strings.ToLower(strings.TrimSpace(dirtyRequestText))
}

func handleGroupMessage(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	requestText := prepareRequestText(update.Message.Text)

	if strings.HasPrefix(requestText, TriggerKeyword) {
		requestText = prepareRequestText(strings.TrimPrefix(requestText, TriggerKeyword))
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		msg.ReplyToMessageID = update.Message.MessageID
		msg.Text = translate(requestText)
		sendMsg(bot, msg)
	}
}

func handlePrivateMessage(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	msg.ReplyToMessageID = update.Message.MessageID
	msg.Text = translate(prepareRequestText(update.Message.Text))
	sendMsg(bot, msg)
}

func handleCommand(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	switch update.Message.Command() {
	case "start":
		msg.Text = StartMessage
	case "help":
		msg.ParseMode = tgbotapi.ModeHTML
		msg.DisableWebPagePreview = true
		msg.Text = HelpMessage
	case "ping":
		msg.Text = "понг"
	}

	sendMsg(bot, msg)
}
