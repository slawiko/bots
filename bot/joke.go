package main

import (
	"math/rand"
	"time"
)

const PROBABILITY = 0.15

func getJokes() []string {
	return []string{"Сорамна такое ня ведаць.", "Я спадзяюся вы запытваеце гэта ў апошні раз.", "Адчуванне, быццам ужо запытвалі.", "Я думаў вы з Беларусі.", "Вы ведаеце такі фільм 'Памятай'?", "Я спадзяюся вы запісваеце, таму што мне не хочацца кожны раз адное і тое ж усім адказваць."}
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
