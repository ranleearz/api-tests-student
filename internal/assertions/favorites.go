package assertions

import (
	"testing"

	"github.com/stretchr/testify/require"

	favoritesModels "api-tests-template/internal/managers/favorites/models"
)

// ContainsFavoriteByAdvertisementId проверяет, что объявление с указанным id есть в списке избранного
func ContainsFavoriteByAdvertisementId(t *testing.T, items []favoritesModels.Favorite, advertisementId string) {
	require.Truef(t, hasFavoriteWithAdvertisementId(items, advertisementId),
		"Среди %d избранных нет объявления с id %q", len(items), advertisementId)
}

// NotContainsFavoriteByAdvertisementId проверяет, что объявления с указанным id нет в списке избранного
func NotContainsFavoriteByAdvertisementId(t *testing.T, items []favoritesModels.Favorite, advertisementId string) {
	require.Falsef(t, hasFavoriteWithAdvertisementId(items, advertisementId),
		"Объявление с id %q всё ещё присутствует в избранном", advertisementId)
}

// hasFavoriteWithAdvertisementId возвращает true, если в списке избранного есть объявление с указанным id
func hasFavoriteWithAdvertisementId(items []favoritesModels.Favorite, advertisementId string) bool {
	for _, item := range items {
		if item.Advertisement.Id == advertisementId {
			return true
		}
	}

	return false
}
