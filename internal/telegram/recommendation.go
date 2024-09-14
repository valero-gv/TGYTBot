package telegram

import "gopkg.in/telebot.v3"

// Функция для создания Inline-кнопки "Получить видео из рекомендаций"
func (b *Bot) ShowRecommendationButton(c telebot.Context) error {
	inlineBtn := telebot.InlineButton{
		Unique: "get_recommendations",
		Text:   "Получить видео из рекомендаций",
		Data:   "get_recommendations",
	}

	inlineKeyboard := [][]telebot.InlineButton{
		{inlineBtn},
	}

	return c.Send("Авторизация прошла успешно! Нажмите кнопку ниже, чтобы получить видео из рекомендаций.", &telebot.ReplyMarkup{
		InlineKeyboard: inlineKeyboard,
	})
}
