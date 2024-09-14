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
		Token:  cfg.TelegramToken, // Замените на ваш токен
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

	savedToken, err := youtube.LoadToken()
	if err == nil && savedToken != nil {
		myBot.Auth.Token = savedToken
		fmt.Println("Токен загружен из файла.")
	} else {
		fmt.Println("Не удалось загрузить токен. Пользователю потребуется авторизация.")
	}

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
