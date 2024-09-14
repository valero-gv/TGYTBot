package youtube

import (
	"context"
	"fmt"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func (a *Auth) GetRecommendedVideos() ([]*youtube.Video, error) {

	if &a.Token == nil {
		return nil, fmt.Errorf("пользователь не авторизован")
	}

	ctx := context.Background()
	service, err := youtube.NewService(ctx, option.WithTokenSource(a.Config.TokenSource(ctx, a.Token)))
	if err != nil {
		return nil, fmt.Errorf("ошибка создания YouTube сервиса: %v", err)
	}
	call := service.Videos.List([]string{"snippet"}).MyRating("like").MaxResults(4)
	response, err := call.Do()
	if err != nil {
		return nil, fmt.Errorf("ошибка получения рекомендованных видео: %v", err)
	}

	return response.Items, nil
}
