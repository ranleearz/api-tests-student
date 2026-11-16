package favorites

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/client/http/favorites"
	"api-tests-template/internal/helpers/response"
	commonModels "api-tests-template/internal/managers/common/models"
	"api-tests-template/internal/managers/favorites/models"
)

// AddFavorite добавляет объявление в избранное пользователя, ожидает 201
func AddFavorite(t *testing.T, token, advertisementId string) models.Favorite {
	favoriteResponse := models.Favorite{}
	response.ReadJSON(t, favorites.HttpPostFavorite(t, token, models.AddFavoriteRequest{AdvertisementId: advertisementId}),
		http.StatusCreated, &favoriteResponse)
	require.Equal(t, advertisementId, favoriteResponse.Advertisement.Id, "В избранное добавлено другое объявление")

	return favoriteResponse
}

// AddFavoriteWithError отправляет запрос на добавление в избранное и возвращает тело ошибки, ожидает expectedStatusCode
func AddFavoriteWithError(t *testing.T, token, advertisementId string, expectedStatusCode int) commonModels.ErrorResponse {
	errorResponse := commonModels.ErrorResponse{}
	response.ReadJSON(t, favorites.HttpPostFavorite(t, token, models.AddFavoriteRequest{AdvertisementId: advertisementId}),
		expectedStatusCode, &errorResponse)

	return errorResponse
}

// GetFavorites возвращает список избранного пользователя, ожидает 200
func GetFavorites(t *testing.T, token string) models.FavoritesList {
	favoritesResponse := models.FavoritesList{}
	response.ReadJSON(t, favorites.HttpGetFavorites(t, token), http.StatusOK, &favoritesResponse)

	return favoritesResponse
}

// GetFavoritesWithError отправляет запрос на получение избранного и возвращает тело ошибки, ожидает expectedStatusCode
func GetFavoritesWithError(t *testing.T, token string, expectedStatusCode int) commonModels.ErrorResponse {
	errorResponse := commonModels.ErrorResponse{}
	response.ReadJSON(t, favorites.HttpGetFavorites(t, token), expectedStatusCode, &errorResponse)

	return errorResponse
}

// RemoveFavorite удаляет объявление из избранного пользователя, ожидает 204
func RemoveFavorite(t *testing.T, token, advertisementId string) {
	response.ReadJSON(t, favorites.HttpDeleteFavorite(t, token, advertisementId), http.StatusNoContent, nil)
}

// RemoveFavoriteWithError отправляет запрос на удаление из избранного и возвращает тело ошибки, ожидает expectedStatusCode
func RemoveFavoriteWithError(t *testing.T, token, advertisementId string, expectedStatusCode int) commonModels.ErrorResponse {
	errorResponse := commonModels.ErrorResponse{}
	response.ReadJSON(t, favorites.HttpDeleteFavorite(t, token, advertisementId), expectedStatusCode, &errorResponse)

	return errorResponse
}
