package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Структура для хранения конфигурационных данных
type Config struct {
	TelegramToken string `json:"telegram_token"`
	ClientID      string `json:"client_id"`
	ClientSecret  string `json:"client_secret"`
	RedirectURL   string `json:"redirect_url"`
}

// Load загружает конфигурацию из файла или переменных окружения
func Load() (*Config, error) {
	// Попробуем загрузить конфигурацию из файла config.json
	file, err := os.Open("C:\\Users\\Валеро\\Desktop\\TGYTBot\\internal\\config\\config.json")
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия файла конфигурации: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := &Config{}
	if err := decoder.Decode(config); err != nil {
		return nil, fmt.Errorf("ошибка декодирования файла конфигурации: %v", err)
	}

	// Проверим, загружены ли все необходимые параметры
	if config.TelegramToken == "" || config.ClientID == "" || config.ClientSecret == "" || config.RedirectURL == "" {
		return nil, fmt.Errorf("не все параметры конфигурации установлены")
	}

	return config, nil
}
