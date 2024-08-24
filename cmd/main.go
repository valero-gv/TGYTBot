package main

import (
	"TGYTBot/internal/telegram"
	"log"
	"time"

	"gopkg.in/telebot.v3"
)

func main() {
	pref := telebot.Settings{
		Token:  "7387581650:AAFdFc-Ezn8pNnMYp09_CkVfeTcXKm60n8M", // Замените на ваш токен
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := telebot.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	myBot := &telegram.Bot{Bot: bot}

	myBot.InitHandlers()

	bot.Handle(telebot.OnText, myBot.HandleDownloadVideo)

	bot.Start()
}
