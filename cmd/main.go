package main

import (
	"TGYTBot/internal/config"
	"TGYTBot/internal/telegram"
	"TGYTBot/internal/youtube"
	"fmt"
	"gopkg.in/telebot.v3"
	"log"
	"net/http"
	"time"
)

func main() {
	cfg, err := config.Load()

	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	pref := telebot.Settings{
		Token:  "7387581650:AAFdFc-Ezn8pNnMYp09_CkVfeTcXKm60n8M", // Замените на ваш токен
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := telebot.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	// Создание экземпляра Auth для YouTube
	auth := youtube.NewAuth(cfg.ClientID, cfg.ClientSecret, cfg.RedirectURL)
	if auth == nil {
		log.Fatalf("Ошибка создания Auth: проверьте ваши clientID, clientSecret и redirectURL")
	}

	myBot := &telegram.Bot{
		Bot:  bot,
		Auth: *auth,
	}

	myBot.InitHandlers()
	bot.Handle(telebot.OnText, myBot.HandleDownloadVideo)

	go func() {
		fmt.Println("Запуск Telegram бота...")
		bot.Start()
	}()

	http.HandleFunc("/callback", auth.HandleCallback)
	fmt.Println("Сервер запущен на http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
