package main

import (
	"math/rand"
	"time"
)

const PROBABILITY = 0.042

func getJokes() []string {
	return []string{"Сорамна такое ня ведаць.", "Я спадзяюся вы запытваеце гэта ў апошні раз.", "Адчуванне, быццам ужо запытвалі.", "А вы сапраўды з Беларусі?"}
}

func joke() bool {
	rand.Seed(time.Now().UnixNano())
	v := rand.Float64()

	return v <= PROBABILITY
}

func jokeMessage() string {
	rand.Seed(time.Now().UnixNano())
	jokes := getJokes()
	randomIndex := rand.Intn(len(jokes))

	return jokes[randomIndex]
}
