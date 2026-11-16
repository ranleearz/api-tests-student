package models

import advertisementsModels "api-tests-template/internal/managers/advertisements/models"

type FavoritesList struct {
	Items []Favorite `json:"items,omitempty"`
	Total int        `json:"total,omitempty"`
}

type Favorite struct {
	Advertisement advertisementsModels.Advertisement `json:"advertisement,omitempty"`
}
