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
		if strings.HasPrefix(strings.ToLower(update.Message.Text), TRANLSATE_KEYWORD) {
			msg := handleTranslateRequest(&update)
			bot.Send(msg)
		}
	}
}

type Suggestion struct {
	Id    int    `json:"id"`
	Label string `json:"label"`
}

func handleTranslateRequest(update *tgbotapi.Update) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	msg.ReplyToMessageID = update.Message.MessageID
	messageText := strings.ToLower(update.Message.Text)
	searchTerm := strings.TrimPrefix(messageText, TRANLSATE_KEYWORD)
	requestUrl := fmt.Sprintf("https://www.skarnik.by/search_json?term=%s&lang=rus", searchTerm)

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
		msg.Text = "Адчапіся, дурны"

		return msg
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

	msg.Text = translation

	return msg
}
