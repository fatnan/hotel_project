package service

import (
	"hotel_project/models"
	"hotel_project/repository"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

type HotelService struct {
	Repo *repository.HotelRepository
}

func NewHotelService(repo *repository.HotelRepository) *HotelService {
	return &HotelService{Repo: repo}
}

func (s *HotelService) ScrapeHotels() ([]map[string]string, error) {
	// Create a new Colly collector
	collector := colly.NewCollector()

	// Variables to store the scraped data
	var hotels []map[string]string

	// Callback function to handle the scraped data
	collector.OnHTML(".hotel", func(e *colly.HTMLElement) {
		hotel := make(map[string]string)
		hotel["name"] = e.ChildText("h3")
		hotel["address"] = e.ChildText(".loct")
		hotel["image_url"] = e.ChildAttr(".img-hotel", "src")
		hotel["star_rating"] = strconv.Itoa(countStarRating(e.ChildAttrs("li i", "class"))) // Convert int to string

		// Remove "IDR" prefix and dot (.) from price
		price := e.ChildText(".price-hotel")
		price = strings.ReplaceAll(price, "IDR", "")
		price = strings.ReplaceAll(price, ".", "")
		price = strings.ReplaceAll(price, " ", "")
		hotel["price"] = price

		hotels = append(hotels, hotel)
	})

	// Start the scraping process
	err := collector.Visit("http://115.85.80.33/test-scrapping/avail.html")
	if err != nil {
		return nil, err
	}

	return hotels, nil
}

func (s *HotelService) InsertHotels(hotels []map[string]string) error {
	for _, hotelData := range hotels {
		starRating, err_star := strconv.Atoi(hotelData["star_rating"])
		if err_star != nil {
			// Handle the error if the conversion fails
			return err_star
		}

		hotel := &models.Hotel{
			Name:       hotelData["name"],
			Address:    hotelData["address"],
			ImageURL:   hotelData["image_url"],
			StarRating: starRating,
			Price:      parsePrice(hotelData["price"]),
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		err := s.Repo.InsertHotel(hotel)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *HotelService) GetAllHotels() ([]models.Hotel, error) {
	return s.Repo.GetAllHotels()
}

func countStarRating(classAttrs []string) int {
	// Define the regular expression pattern to match the star rating class
	starClassRegex := regexp.MustCompile(`\bstar-hotel\b`)

	count := 0

	// Iterate over each class attribute and count the star ratings
	for _, classAttr := range classAttrs {
		matches := starClassRegex.FindAllString(classAttr, -1)
		count += len(matches)
	}

	return count
}

func parsePrice(price string) int {
	price = strings.ReplaceAll(price, "IDR", "")
	price = strings.ReplaceAll(price, ".", "")
	price = strings.ReplaceAll(price, " ", "")
	parsedPrice, err := strconv.Atoi(price)
	if err != nil {
		log.Println("parse error")
		return 0
	}
	return parsedPrice
}
