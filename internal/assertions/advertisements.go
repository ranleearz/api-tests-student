package assertions

import (
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/advertisements/models"
)

func ContainsAdvertisementByTitle(t *testing.T, items []models.Advertisement, title string) {
	for _, item := range items {
		if item.Title == title {
			return
		}
	}

	require.Failf(t, "Объявление не найдено", "Среди %d объявлений нет объявления с названием %q", len(items), title)
}
