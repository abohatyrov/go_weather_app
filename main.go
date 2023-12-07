package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
	"net/http"
	"os"
)

type WeatherData struct {
	Main struct {
		Temp     float64 `json:"temp"`
		Humidity int     `json:"humidity"`
	} `json:"main"`
	Weather []struct {
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	apiKey := os.Getenv("OPENWEATHERMAP_API_KEY")

	router := gin.Default()

	router.Use(MetricsMiddleware())

	router.Static("/static", "./static")

	router.LoadHTMLGlob("templates/*")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	router.GET("/weather", func(c *gin.Context) {
		city := c.Query("city")

		if city == "" {
			c.JSON(400, gin.H{"error": "City parameter is required"})
			return
		}

		client := resty.New()
		resp, err := client.R().
			SetQueryParams(map[string]string{
				"q":     city,
				"appid": apiKey,
			}).
			Get("https://api.openweathermap.org/data/2.5/weather")

		if err != nil {
			c.JSON(500, gin.H{"error": "Error contacting OpenWeatherMap API"})
			return
		}

		if resp.IsSuccess() {
			var weatherData WeatherData
			err := json.Unmarshal(resp.Body(), &weatherData)

			if err != nil {
				c.JSON(500, gin.H{"error": "Error parsing OpenWeatherMap API response"})
				return
			}

			temperatureCelsius := weatherData.Main.Temp - 273.15
			weatherIconURL := fmt.Sprintf("http://openweathermap.org/img/w/%s.png", weatherData.Weather[0].Icon)

			c.HTML(http.StatusOK, "weather.html", gin.H{
				"cityName":           city,
				"temperature":        temperatureCelsius,
				"weatherDescription": weatherData.Weather[0].Description,
				"humidity":           weatherData.Main.Humidity,
				"weatherIconURL":     weatherIconURL,
			})
		} else {
			c.JSON(resp.StatusCode(), gin.H{"error": resp.Status()})
		}
	})

	router.GET("/metrics", MetricsHandler)

	router.Run(":8080")
}
