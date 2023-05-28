// repository/hotel_repository.go
package repository

import (
	"database/sql"
	"log"
	"time"

	"hotel_project/models"
)

type HotelRepository struct {
	DB *sql.DB
}

func NewHotelRepository(db *sql.DB) *HotelRepository {
	return &HotelRepository{DB: db}
}

func (r *HotelRepository) InsertHotel(hotel *models.Hotel) error {
	stmt, err := r.DB.Prepare("INSERT INTO hotel (name, address, image_url, star_rating, price, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Get the current timestamp for created_at and updated_at fields
	currentTime := time.Now()

	// Execute the INSERT statement
	_, err = stmt.Exec(hotel.Name, hotel.Address, hotel.ImageURL, hotel.StarRating, hotel.Price, currentTime, currentTime)
	if err != nil {
		return err
	}

	log.Println("Hotel data inserted successfully.")
	return nil
}

func (r *HotelRepository) GetAllHotels() ([]models.Hotel, error) {
	rows, err := r.DB.Query("SELECT * FROM hotel")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var hotels []models.Hotel

	for rows.Next() {
		var hotel models.Hotel
		var createdAt, updatedAt string
		err := rows.Scan(&hotel.ID, &hotel.Name, &hotel.Address, &hotel.ImageURL, &hotel.StarRating, &hotel.Price, &createdAt, &updatedAt)
		if err != nil {
			log.Println(err)
			continue
		}
		hotel.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		hotel.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)
		hotels = append(hotels, hotel)
	}

	return hotels, nil
}
