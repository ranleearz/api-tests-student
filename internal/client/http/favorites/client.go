package favorites

import (
	"fmt"
	"net/http"
	"testing"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
	favoritesModels "api-tests-template/internal/managers/favorites/models"
)

// HttpPostFavorite выполняет POST-запрос для добавления объявления в избранное
func HttpPostFavorite(t *testing.T, token string, request favoritesModels.AddFavoriteRequest) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Post(path.FavoritesPath).
		ContentType("application/json").
		JSON(request).
		Expect(t).
		End().Response
}

// HttpGetFavorites выполняет GET-запрос для получения списка избранного пользователя
func HttpGetFavorites(t *testing.T, token string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Get(path.FavoritesPath).
		Expect(t).
		End().Response
}

// HttpDeleteFavorite выполняет DELETE-запрос для удаления объявления из избранного
func HttpDeleteFavorite(t *testing.T, token, advertisementId string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Delete(fmt.Sprintf(path.FavoritePath, advertisementId)).
		Expect(t).
		End().Response
}
