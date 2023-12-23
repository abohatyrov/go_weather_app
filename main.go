package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
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

var logger = logrus.New()

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		logger.WithFields(logrus.Fields{
			"status_code": c.Writer.Status(),
			"method":      c.Request.Method,
			"path":        c.Request.URL.Path,
			"ip":          c.ClientIP(),
		}).Info()
	}
}

func main() {
	logger.SetFormatter(&logrus.JSONFormatter{})

	apiKey := os.Getenv("OPENWEATHERMAP_API_KEY")
	if apiKey == "" {
		logger.Fatal("OPENWEATHERMAP_API_KEY environment variable is required")
		return
	}

	router := gin.Default()

	router.Use(Logger())
	router.Use(MetricsMiddleware())

	router.Static("/static", "./static")

	router.LoadHTMLGlob("templates/*")

	router.GET("/", func(c *gin.Context) {
		logger.Info("Root endpoint hit")
		c.HTML(http.StatusOK, "index.html", nil)
	})

	router.GET("/weather", func(c *gin.Context) {
		city := c.Query("city")

		if city == "" {
			c.JSON(400, gin.H{"error": "City parameter is required"})
			logger.Warn("Weather endpoint hit without city parameter")
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
			logger.Error("Error contacting OpenWeatherMap API: ", err)
			return
		}

		if resp.IsSuccess() {
			var weatherData WeatherData
			err := json.Unmarshal(resp.Body(), &weatherData)

			if err != nil {
				c.JSON(500, gin.H{"error": "Error parsing OpenWeatherMap API response"})
				logger.Error("Error parsing OpenWeatherMap API response: ", err)
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
			logger.Fatal("Error from OpenWeatherMap API: ", resp.Status())
		}
	})

	router.GET("/metrics", MetricsHandler)

	logger.Info("Server starting on port 8080")
	router.Run(":8080")
}
