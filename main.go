package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const TRANLSATE_KEYWORD = "як будзе "

var BOT_API_KEY = os.Args[1]

func main() {
	bot, err := tgbotapi.NewBotAPI(BOT_API_KEY)

	if err != nil {
		log.Println("Error occurred")
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
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
		msg.Text = "Ня змог чамусьці. Стварыце калі ласка ішшу на гітхабе https://github.com/slawiko/ru-bel-tg-bot/issues"

		_, err := bot.Send(msg)
		if err != nil {
			log.Println(err)
		}
	}
}

func handleGroupMessage(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	if strings.HasPrefix(strings.ToLower(update.Message.Text), TRANLSATE_KEYWORD) {
		requestText := strings.TrimPrefix(strings.ToLower(update.Message.Text), TRANLSATE_KEYWORD)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		msg.ReplyToMessageID = update.Message.MessageID
		msg.Text = translate(requestText)
		sendMsg(bot, msg)
	}
}

func handlePrivateMessage(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	// msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyToMessageID = update.Message.MessageID
	msg.Text = translate(strings.ToLower(update.Message.Text))
	sendMsg(bot, msg)
}

type Suggestion struct {
	Id    int    `json:"id"`
	Label string `json:"label"`
}

func cleanTerm(searchTerm string) string {
	cleanSearchTerm := strings.ReplaceAll(searchTerm, "ў", "щ")
	cleanSearchTerm = strings.ReplaceAll(cleanSearchTerm, "і", "и")
	cleanSearchTerm = strings.ReplaceAll(cleanSearchTerm, "’", "ъ")
	cleanSearchTerm = strings.ReplaceAll(cleanSearchTerm, "'", "ъ")

	return cleanSearchTerm
}

func translate(searchTerm string) string {
	cleanSearchTerm := cleanTerm(searchTerm)

	requestUrl := fmt.Sprintf("https://www.skarnik.by/search_json?term=%s&lang=rus", cleanSearchTerm)

	resp, err := http.Get(requestUrl)
	if err != nil {
		log.Println(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	var suggestions []Suggestion
	json.Unmarshal(body, &suggestions)

	if len(suggestions) == 0 {
		return "Адчапіся, дурны"
	}

	requestUrl = fmt.Sprintf("https://www.skarnik.by/rusbel/%d", suggestions[0].Id)

	resp, err = http.Get(requestUrl)
	if err != nil {
		log.Println(err)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Println(err)
	}

	section := doc.Find("#trn")

	translation := ""

	section.Find("font[color=\"831b03\"]").Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			translation += s.Text()
		} else {
			translation += fmt.Sprintf(", %s", s.Text())
		}
	})

	log.Println(translation)

	return translation

	// translation, err := section.Html()
	// if err != nil {
	// 	log.Println(err)
	// }

	// formatted := strings.ReplaceAll(translation, "<br/>", "\n")

	// section.Find("font[color=\"008000\"]").Each(func(i int, s *goquery.Selection) {
	// 	greenHtml, err := goquery.OuterHtml(s)
	// 	if err != nil {
	// 		log.Println(err)
	// 	}
	// 	greenText := s.Text()

	// 	formatted = strings.ReplaceAll(formatted, greenHtml, greenText)
	// })

	// section.Find("font[color=\"831b03\"]").Each(func(i int, s *goquery.Selection) {
	// 	bigRedHtml, err := goquery.OuterHtml(s)
	// 	if err != nil {
	// 		log.Println(err)
	// 	}
	// 	bigRedText := s.Text()

	// 	formatted = strings.ReplaceAll(formatted, bigRedHtml, "<code>"+bigRedText+"</code>")
	// })

	// section.Find("font[color=\"5f5f5f\"]").Each(func(i int, s *goquery.Selection) {
	// 	strongHtml, err := goquery.OuterHtml(s)
	// 	if err != nil {
	// 		log.Println(err)
	// 	}

	// 	strongText := s.Text()

	// 	fmt.Println(strongHtml)

	// 	formatted = strings.ReplaceAll(formatted, strongHtml, "<i>"+strongText+"</i>")
	// })

	// fmt.Println(formatted)

	// return formatted
}
