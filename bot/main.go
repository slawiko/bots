package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	TriggerKeyword     = "як будзе "
	ErrorMessage       = "Нешта чамусці пайшло ня так. Стварыце калі ласка ішшу на гітхабе https://github.com/slawiko/ru-bel-bot/issues"
	EmptyResultMessage = "Нічога не знайшоў :("
	HelpMessage        = `Спосабы ўзаемадзеяння:
<b>У прываце</b>: наўпрост пішыце слова на рускай мове. Увага: я лагірую тэкст, што вы напішаце.
<b>У группе</b>: пачніце ваша паведамленне са словаў <code>як будзе</code> і далей слово на русском языке. Напрыклад: <code>як будзе письмо</code>.
Увага: тут я лагірую толькі факт карыстання гэтай функцыяй 

Таксама вы можаце не пераходзіць на рускую раскладку і пытацца, напрыклад, слова <code>ўавель</code> ці <code>олівка</code>.

У тым выпадку, калі вы баіцеся дадаць мяне ў вашыя чаты, вы можаце запусціць мяне самастойна. Інструкцыя тут: https://github.com/slawiko/ru-bel-bot/blob/master/README.md#run. Калі нешта незразумела - пішыце ў https://github.com/slawiko/ru-bel-bot/issues

<i>На дадзены момант я ня разумею памылкі, прабачце.</i>

© Усе пераклады я бяру з https://skarnik.by, дзякуй яму вялікі.`
	StartMessage = `Прывітаннечка. Мяне клічуць Жэўжык, я дапамагаю перайсці на родную мову. Вы можаце пытацца ў мяне слова на рускай, а я адкажу вам на беларускай.

Вы можаце дадаць мяне ў группу і пытацца не выходзячы з дыялогу з сябрамі. За дапамогай клацайце /help`
	DetailedButton = "Падрабязней"
	ShortButton    = "Карацей"
)

var BotApiKey = os.Args[1]
var Version = os.Getenv("VERSION")

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
			log.Println("callback") // do not log callback requests, since it could go from group
			handleCallback(bot, update.CallbackQuery)
			continue
		}
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			log.Println("command", update.Message.Command())
			handleCommand(bot, &update)
		} else if update.Message.Chat.IsGroup() || update.Message.Chat.IsSuperGroup() {
			log.Println("group") // do not log group message requests, since there could be sensitive data
			handleGroupMessage(bot, &update)
		} else if update.Message.Chat.IsPrivate() {
			log.Println("private", update.Message.Text)
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

func PrepareRequestText(searchTerm string) string {
	cleanSearchTerm := strings.ToLower(strings.TrimSpace(searchTerm))
	cleanSearchTerm = strings.ReplaceAll(cleanSearchTerm, "ў", "щ")
	cleanSearchTerm = strings.ReplaceAll(cleanSearchTerm, "і", "и")
	cleanSearchTerm = strings.ReplaceAll(cleanSearchTerm, "’", "ъ")
	cleanSearchTerm = strings.ReplaceAll(cleanSearchTerm, "'", "ъ")

	return cleanSearchTerm
}

func handleGroupMessage(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	requestText := PrepareRequestText(update.Message.Text)

	if strings.HasPrefix(requestText, TriggerKeyword) {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		msg.ReplyToMessageID = update.Message.MessageID
		msg.ParseMode = tgbotapi.ModeHTML

		requestText = strings.TrimPrefix(requestText, TriggerKeyword)
		translation, err := Translate(requestText, false)
		if err != nil {
			msg.Text = EmptyResultMessage
			log.Println(err)
		} else {
			msg.Text = *translation
			keyboard := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(DetailedButton, marshallCallbackData(requestText, true))),
			)
			msg.ReplyMarkup = keyboard
		}

		sendMsg(bot, msg)
	}
}

func handlePrivateMessage(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	msg.ReplyToMessageID = update.Message.MessageID
	msg.ParseMode = tgbotapi.ModeHTML

	requestText := PrepareRequestText(update.Message.Text)
	translation, err := Translate(requestText, false)
	if err != nil {
		msg.Text = EmptyResultMessage
		log.Println(err)
	} else {
		msg.Text = *translation
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(DetailedButton, marshallCallbackData(requestText, true))),
		)
		msg.ReplyMarkup = keyboard
	}

	sendMsg(bot, msg)
}

func marshallCallbackData(word string, shouldNextBeDetailed bool) string {
	return fmt.Sprintf("%s$%v", word, shouldNextBeDetailed)
}

func unmarshallCallbackData(data string) (string, bool) {
	parts := strings.Split(data, "$")
	isDetailed, _ := strconv.ParseBool(parts[1])
	return parts[0], isDetailed
}

func handleCallback(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery) {
	editMsg := tgbotapi.NewEditMessageText(callback.Message.Chat.ID, callback.Message.MessageID, "")
	editMsg.ParseMode = tgbotapi.ModeHTML
	word, isDetailed := unmarshallCallbackData(callback.Data)

	translation, err := Translate(word, isDetailed)
	var buttonText string
	if isDetailed {
		buttonText = ShortButton
	} else {
		buttonText = DetailedButton
	}
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(buttonText, marshallCallbackData(word, !isDetailed))),
	)
	editMsg.ReplyMarkup = &keyboard
	if err != nil {
		log.Println(err)
		editMsg.Text = EmptyResultMessage
	} else {
		editMsg.Text = *translation
	}

	_, err = bot.Send(editMsg)
	if err != nil {
		log.Println(err)
	}

	bot.Request(tgbotapi.NewCallback(callback.ID, "")) // for hiding alert. Looks wrong, but donno how else
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
	case "version":
		if len(Version) > 0 {
			msg.ParseMode = tgbotapi.ModeHTML
			msg.Text = fmt.Sprintf("<a href=\"https://github.com/slawiko/ru-bel-bot/releases/tag/%s\">%s</a>", Version, Version)
		} else {
			msg.Text = "unknown"
		}
	}

	sendMsg(bot, msg)
}
