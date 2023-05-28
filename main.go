package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"hotel_project/repository"
	"hotel_project/service"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app := fiber.New()

	// Open a connection to the MySQL database
	// db, err := sql.Open("mysql", "root@tcp(localhost:3306)/hotel_project")
	db, err := sql.Open("mysql", os.Getenv("DB_USER")+":"+os.Getenv("DB_PASSWORD")+"@tcp("+os.Getenv("DB_HOST")+":"+os.Getenv("DB_PORT")+")/"+os.Getenv("DB_NAME"))

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Initialize the repository and service
	hotelRepo := repository.NewHotelRepository(db)
	hotelService := service.NewHotelService(hotelRepo)

	// Endpoint for scrape and insert hotels
	app.Get("/scraping", func(c *fiber.Ctx) error {
		hotels, err := hotelService.ScrapeHotels()
		if err != nil {
			log.Println("Scraping error:", err)
			return c.Status(fiber.StatusInternalServerError).SendString("Scraping error")
		}
		//insert hotels to db
		err = hotelService.InsertHotels(hotels)
		if err != nil {
			log.Println("Error inserting hotels:", err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error inserting hotels")
		}

		// Return the scraped data as JSON response
		return c.JSON(hotels)
	})

	// Endpoint to get list hotels
	app.Get("/hotels", func(c *fiber.Ctx) error {
		hotels, err := hotelService.GetAllHotels()
		if err != nil {
			log.Println("Error getting hotels:", err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error getting hotels")
		}

		jsonData, err := json.Marshal(hotels)
		if err != nil {
			log.Println(err)
			return c.SendStatus(http.StatusInternalServerError)
		}

		return c.Send(jsonData)
	})

	app.Listen(":3000")
}
