package main

import (
	"log"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	TriggerKeyword     = "як будзе"
	ErrorMessage       = "Нешта чамусьці пайшло ня так. Стварыце калі ласка ішшу на гітхабе https://github.com/slawiko/ru-bel-bot/issues"
	EmptyResultMessage = "Нічога не знайшоў :("
	HelpMessage        = `Спосабы ўзаемадзеяння:
<b>У прываце</b>: наўпрост пішыце слова на рускай мове.
<b>У группе</b>: пачніце ваша паведамленне са словаў <code>як будзе</code> і далей слово на русском языке. Напрыклад: <code>як будзе письмо</code>.

Таксама вы можаце не пераходзіць на рускую раскладку і пытацца, напрыклад, слова <code>ўавель</code> ці <code>олівка</code>.

У тым выпадку, калі вы баіцеся дадаць мяне ў вашыя чаты, вы можаце запусьціць мяне самастойна. Інструкцыя тут: https://github.com/slawiko/ru-bel-bot/blob/master/README.md#run. Калі нешта незразумела - пішыце ў https://github.com/slawiko/ru-bel-bot/issues

<i>На дадзены момант я ня разумею памылкі, прабачце.</i>

© Усе пераклады я бяру з https://skarnik.by, дзякуй яму вялікі.`
	StartMessage = `Прывітаннечка. Мяне клічуць Жэўжык, я дапамагаю перайсьці на родную мову. Вы можаце пытацца ў мяне слова на рускай, а я адкажу вам на беларускай.

Вы можаце дадаць мяне ў группу і пытацца не выходзячы з дыялогу з сябрамі. За дапамогай клацайце /help`
)

var BotApiKey = os.Args[1]

func main() {
	bot, err := tgbotapi.NewBotAPI(BotApiKey)
	if err != nil {
		log.Println(err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.CallbackQuery != nil {
			handleFullTranslationRequest(bot, update.CallbackQuery)
			continue
		}
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			handleCommand(bot, &update)
		} else if update.Message.Chat.IsGroup() || update.Message.Chat.IsSuperGroup() {
			handleGroupMessage(bot, &update)
		} else if update.Message.Chat.IsPrivate() {
			handlePrivateMessage(bot, &update)
		}
	}
}

func sendMsg(bot *tgbotapi.BotAPI, msg tgbotapi.MessageConfig) {
	msg.DisableNotification = true

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
		msg.ParseMode = tgbotapi.ModeHTML
		// TODO: do not show button for multi words request
		translation, err := translate(requestText, false)
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Падрабязней", prepareRequestText(requestText))),
		)
		msg.ReplyMarkup = keyboard
		if err != nil {
			log.Println(err)
			msg.Text = EmptyResultMessage
		} else {
			msg.Text = *translation
			log.Println(*translation)
		}
		sendMsg(bot, msg)
	}
}

func handlePrivateMessage(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	msg.ReplyToMessageID = update.Message.MessageID
	msg.ParseMode = tgbotapi.ModeHTML
	translation, err := translate(prepareRequestText(update.Message.Text), false)
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Падрабязней", prepareRequestText(update.Message.Text))),
	)
	msg.ReplyMarkup = keyboard
	if err != nil {
		log.Println(err)
		msg.Text = EmptyResultMessage
	} else {
		msg.Text = *translation
		log.Println(*translation)
	}
	sendMsg(bot, msg)
}

func handleFullTranslationRequest(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery) {
	editMsg := tgbotapi.NewEditMessageText(callback.Message.Chat.ID, callback.Message.MessageID, "")
	editMsg.ParseMode = tgbotapi.ModeHTML

	fullTranslation, err := translate(prepareRequestText(callback.Data), true)
	if err != nil {
		log.Println(err)
		editMsg.Text = EmptyResultMessage
	}

	editMsg.Text = *fullTranslation

	_, err = bot.Send(editMsg)
	if err != nil {
		log.Println(err)
	}
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
