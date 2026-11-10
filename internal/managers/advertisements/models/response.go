package models

type AdvertisementsList struct {
	Items []Advertisement `json:"items,omitempty"`
}

type Advertisement struct {
	Title string `json:"title,omitempty"`
}
