package telegram

import (
	"TGYTBot/internal/youtube"
	"fmt"
	"gopkg.in/telebot.v3"
	"log"
	"os"
	"os/exec"
	"time"
)

type Bot struct {
	Bot  *telebot.Bot
	Auth youtube.Auth
}

// Клавиатура с кнопками
var mainMenu = &telebot.ReplyMarkup{
	ResizeKeyboard: true,
}

var authMenu = &telebot.ReplyMarkup{
	ResizeKeyboard: true,
}

var (
	btnDownload            = mainMenu.Text("Скачать видео")
	btnAuthorize           = telebot.Btn{Text: "Авторизация в YouTube"}
	btnViewRecommendations = telebot.Btn{Text: "Показать рекомендации YouTube"}
)

func (b *Bot) InitHandlers() {
	b.Bot.Handle("/start", b.handleStart)
	b.Bot.Handle(&btnDownload, b.handleDownloadRequest)
	b.Bot.Handle(&btnAuthorize, b.HandleDownloadVideo)
	b.Bot.Handle(&btnAuthorize, b.handleYouTubeAuth)
	b.Bot.Handle(&btnViewRecommendations, b.HandleGetRecommended)
}

// handleStart обработчик команды /start
func (b *Bot) handleStart(c telebot.Context) error {
	mainMenu.Reply(
		mainMenu.Row(btnDownload),
		mainMenu.Row(btnAuthorize),
		mainMenu.Row(btnViewRecommendations),
	)
	return c.Send("Выберите команду:", mainMenu)
}

func (b *Bot) handleDownloadRequest(c telebot.Context) error {
	return c.Send("Пожалуйста, отправьте ссылку на видео для скачивания.")
}

func (b *Bot) HandleDownloadVideo(c telebot.Context) error {
	videoURL := c.Text()

	if videoURL == "" {
		return c.Send("Пожалуйста, отправьте ссылку на видео.")
	}
	c.Send("Скачивание началось.")

	filePath := "video_" + time.Now().Format("20060102150405") + ".mp4"

	cmd := exec.Command("yt-dlp", "--newline", "-o", filePath, videoURL)
	err := cmd.Run()
	if err != nil {
		log.Println("Ошибка скачивания видео:", err)
		return c.Send("Не удалось скачать видео. Проверьте ссылку и попробуйте снова.")
	}
	video := &telebot.Video{File: telebot.FromDisk(filePath)}
	if err := c.Send(video); err != nil {
		log.Println("Ошибка отправки видео:", err)
		return c.Send("Не удалось отправить видео.")
	}

	if err := os.Remove(filePath); err != nil {
		log.Println("Ошибка удаления файла:", err)
	}

	return nil
}

func (b *Bot) handleYouTubeAuth(c telebot.Context) error {
	// Запуск процесса авторизации через OAuth 2.0
	authURL := b.Auth.StartAuth()
	msg := "Перейдите по следующей ссылке для авторизации в YouTube: " + authURL
	return c.Send(msg)
}

func (b *Bot) HandleGetRecommended(c telebot.Context) error {
	if err := b.Auth.CheckAndRefreshToken(); err != nil {
		return c.Send("Ошибка обновления токена: " + err.Error())
	}

	// Получаем рекомендованные видео
	videos, err := b.Auth.GetRecommendedVideos()
	if err != nil {
		_ = fmt.Errorf("ошибка получения видео: %v", err)
		return c.Send(fmt.Sprintf("Ошибка получения видео: %v", err))
	}

	// Если видео нет
	if len(videos) == 0 {
		return c.Send("Нет рекомендованных видео.")
	}
	_, err = b.Bot.Send(c.Chat(), &telebot.Photo{
		File: telebot.FromURL(videos[0].Snippet.Thumbnails.Default.Url),
	})
	if err != nil {
		return c.Send(fmt.Sprintf("Ошибка отправки изображения: %v", err))
	}
	// Отправляем информацию о каждом видео
	//for _, video := range videos {
	//	// Создаем сообщение с превью и названием
	//	//videoTitle := video.Snippet.Title
	//	//videoThumbnail := video.Snippet.Thumbnails.Default.Url
	//	//videoURL := fmt.Sprintf("https://www.youtube.com/watch?v=%s", video.Id)
	//
	//	_, err := b.Bot.Send(c.Chat(), "LOL")
	//	if err != nil {
	//		return c.Send(fmt.Sprintf("Ошибка отправки изображения: %v", err))
	//	}
	//}

	return nil
}
