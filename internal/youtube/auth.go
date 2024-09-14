package youtube

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
	"net/http"
)

type Auth struct {
	Config *oauth2.Config
	Token  *oauth2.Token
	config *oauth2.Config
	UserID int64
}

// Инициализация OAuth 2.0 конфигурации
func NewAuth(clientID, clientSecret, redirectURL string) *Auth {
	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			youtube.YoutubeReadonlyScope,
		},
		Endpoint: google.Endpoint,
	}

	return &Auth{Config: config}
}

// Начало процесса авторизации
func (a *Auth) StartAuth() string {
	authURL := a.Config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	return authURL
}

// HandleCallback обрабатывает ответ Google с кодом авторизации
func (a *Auth) HandleCallback(w http.ResponseWriter, r *http.Request) {
	if a.Config == nil {
		http.Error(w, "OAuth конфигурация не инициализирована", http.StatusInternalServerError)
		return
	}
	// Получение кода авторизации из запроса
	code := r.URL.Query().Get("code")

	// Обмен кода на токен
	token, err := a.Config.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Ошибка обмена кода на токен", http.StatusInternalServerError)
		return
	}
	// Проверяем наличие обновляющего токена
	if token.RefreshToken == "" {
		http.Error(w, "Не удалось получить refresh token. Попробуйте еще раз.", http.StatusInternalServerError)
		return
	}

	SaveToken(token)
	a.Token = token

	fmt.Println(token)

	http.Redirect(w, r, "https://t.me/GVYT_bot", http.StatusFound)

	fmt.Println("Авторизация успешна! Теперь вы можете вернуться в Telegram.")
}
