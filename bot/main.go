package main

import (
	"log"
	"os"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	TriggerKeyword = "як будзе"
	ErrorMessage   = "Нешта чамусьці пайшло ня так. Стварыце калі ласка ішшу на гітхабе https://github.com/slawiko/ru-bel-bot/issues"
	EmptyResultMessage = "Нічога не знайшоў :("
	HelpMessage    = `Спосабы ўзаемадзеяння:
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
		if update.InlineQuery != nil {
			handleInlineQuery(bot, &update)
			continue
		}

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

func handleInlineQuery(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	if len(update.InlineQuery.Query) == 0 {
		inlineConf := tgbotapi.InlineConfig{
			InlineQueryID: update.InlineQuery.ID,
		}
		_, err := bot.Request(inlineConf)
		if err != nil {
			log.Println(err)
		}
		return
	}
	log.Println(update.InlineQuery.Query)
	suggestions, err := getScarnikSuggestions(update.InlineQuery.Query)
	if err != nil {
		inlineConf := tgbotapi.InlineConfig{
			InlineQueryID: update.InlineQuery.ID,
		}
		_, err = bot.Request(inlineConf)
		if err != nil {
			log.Println(err)
		}
		return
	}

	if len(suggestions) == 0 {
		inlineConf := tgbotapi.InlineConfig{
			InlineQueryID: update.InlineQuery.ID,
		}
		_, err := bot.Request(inlineConf)
		if err != nil {
			log.Println(err)
		}
		return
	}

	articles := []tgbotapi.InlineQueryResultArticle{}

	for i := 0; i < len(suggestions); i++ {
		if i > 2 {
			break
		}

		resp, err := requestSkarnik(suggestions[i])
		if err != nil {
			log.Println(err)
			continue
		}

		suggestionTranslation, err := parseSkarnikResponse(resp)
		log.Println(*suggestionTranslation)
		if err != nil {
			log.Println(err)
			continue
		}

		article := tgbotapi.NewInlineQueryResultArticle(strconv.Itoa(suggestions[i].ID), suggestions[i].Label, *suggestionTranslation)
		article.Description = *suggestionTranslation
		articles = append(articles, article)
	}

//	for _, suggestion := range suggestions[0:4] {
//		resp, err := requestSkarnik(suggestion)
//		if err != nil {
//			log.Println(err)
//			continue
//		}
//
//		suggestionTranslation, err := parseSkarnikResponse(resp)
//		if err != nil {
//			log.Println(err)
//			continue
//		}
//
//		article := tgbotapi.NewInlineQueryResultArticle(strconv.Itoa(suggestion.ID), suggestion.Label, *suggestionTranslation)
//		article.Description = *suggestionTranslation
//		articles = append(articles, article)
//	}

	results := make([]interface{}, len(articles))
	for i, v := range articles {
		results[i] = v
	}

	inlineConf := tgbotapi.InlineConfig{
		InlineQueryID: update.InlineQuery.ID,
		Results: results,
		IsPersonal: true,
	}
	_, err = bot.Request(inlineConf)
	if err != nil {
		log.Println("Request fail", len(results), err)
		for _, e := range articles {
			log.Println(e)
			if e.InputMessageContent == nil {
				log.Println(e.Description)
			}
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
		translation, err := translate(requestText)
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
	translation, err := translate(prepareRequestText(update.Message.Text))
	if err != nil {
		log.Println(err)
		msg.Text = EmptyResultMessage
	} else {
		msg.Text = *translation
		log.Println(*translation)
	}
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
