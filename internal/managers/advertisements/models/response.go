package models

type AdvertisementsList struct {
	Items []Advertisement `json:"items,omitempty"`
}

type Advertisement struct {
	Id          string  `json:"id,omitempty"`
	Title       string  `json:"title,omitempty"`
	Description string  `json:"description,omitempty"`
	Price       float64 `json:"price,omitempty"`
	Quantity    int     `json:"quantity,omitempty"`
	UserId      string  `json:"user_id,omitempty"`
	Photos      []Photo `json:"photos,omitempty"`
}

type Photo struct {
	Id  string `json:"id,omitempty"`
	Url string `json:"url,omitempty"`
}
