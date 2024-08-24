package telegram

import (
	"bufio"
	"fmt"
	"gopkg.in/telebot.v3"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Bot struct {
	Bot *telebot.Bot
}

// Клавиатура с кнопками
var mainMenu = &telebot.ReplyMarkup{
	ResizeKeyboard: true,
}

var (
	btnDownload = mainMenu.Text("Скачать видео")
)

func (b *Bot) InitHandlers() {
	b.Bot.Handle("/start", b.handleStart)
	b.Bot.Handle(&btnDownload, b.handleDownloadRequest)
}

// handleStart обработчик команды /start
func (b *Bot) handleStart(c telebot.Context) error {
	mainMenu.Reply(
		mainMenu.Row(btnDownload),
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

	filePath := "video_" + time.Now().Format("20060102150405") + ".mp4"

	cmd := exec.Command("yt-dlp", "-o", filePath, videoURL)
	// Создаем pipe для чтения вывода yt-dlp
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Println("Ошибка создания pipe:", err)
		return c.Send("Произошла ошибка при подготовке к загрузке видео.")
	}

	if err := cmd.Start(); err != nil {
		log.Println("Ошибка запуска команды:", err)
		return c.Send("Произошла ошибка при запуске загрузки видео.")
	}

	// Читаем вывод yt-dlp в реальном времени
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line) // Логируем весь вывод yt-dlp

		// Парсим строку для получения процента загрузки
		if strings.Contains(line, "[download]") {
			parts := strings.Fields(line)
			if len(parts) > 1 && strings.HasSuffix(parts[1], "%") {
				progress := parts[1]
				c.Send(fmt.Sprintf("Прогресс загрузки: %s", progress))
			}
		}
	}

	// Проверяем ошибки чтения
	if err := scanner.Err(); err != nil {
		log.Println("Ошибка чтения вывода yt-dlp:", err)
		return c.Send("Произошла ошибка при чтении данных загрузки.")
	}

	// Ожидаем завершения команды
	if err := cmd.Wait(); err != nil {
		log.Println("Ошибка завершения команды yt-dlp:", err)
		return c.Send("Произошла ошибка при завершении загрузки видео.")
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
