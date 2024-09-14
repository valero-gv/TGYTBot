package youtube

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"log"
	"os"
)

func SaveToken(token *oauth2.Token) {
	// Сохраняем токен в файл или базу данных
	tokenJSON, err := json.Marshal(token)
	if err != nil {
		log.Fatalf("Ошибка сериализации токена: %v", err)
	}

	err = os.WriteFile("token.json", tokenJSON, 0600)
	if err != nil {
		log.Fatalf("Ошибка сохранения токена: %v", err)
	}
}

func LoadToken() (*oauth2.Token, error) {
	// Загружаем токен из файла или базы данных
	file, err := os.Open("C:\\Users\\Валеро\\Desktop\\TGYTBot\\token.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	token := &oauth2.Token{}
	err = json.NewDecoder(file).Decode(token)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (a *Auth) CheckAndRefreshToken() error {
	if a.Token == nil || !a.Token.Valid() {
		return a.RefreshToken()
	}
	return nil
}

func (a *Auth) RefreshToken() error {
	if a.Token == nil || !a.Token.Valid() {
		if a.Token.RefreshToken == "" {
			return fmt.Errorf("Refresh token не установлен")
		}
		// Используем Refresh Token для получения нового Access Token
		tokenSource := a.Config.TokenSource(context.Background(), a.Token)
		newToken, err := tokenSource.Token()
		if err != nil {
			return fmt.Errorf("Ошибка обновления токена: %v", err)
		}

		a.Token = newToken
		SaveToken(newToken) // Обновляем сохраненный токен
	}

	return nil
}
