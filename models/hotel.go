package models

import "time"

type Hotel struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Address    string    `json:"address"`
	ImageURL   string    `json:"image_url"`
	StarRating int       `json:"star_rating"`
	Price      int       `json:"price"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
