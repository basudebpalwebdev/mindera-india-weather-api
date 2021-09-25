package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

type Weather struct {
	Temperature float64 `json:"temperature"`
	WindSpeed   float64 `json:"wind_speed"`
}

type WeatherStackResponse struct {
	Current Weather `json:"current"`
}

type Main struct {
	Temp float64 `json:"temp"`
}

type Wind struct {
	Speed float64 `json:"speed"`
}

type OpenWeatherResponse struct {
	Main Main `json:"main"`
	Wind Wind `json:"wind"`
}

type ENVConfig struct {
	OpenWeatherAPIKey  string `mapstructure:"OPEN_WEATHER_API_KEY"`
	WeatherStackAPIKey string `mapstructure:"WEATHER_STACK_API_KEY"`
}

// LoadConfig reads configuration from app.env file.
func LoadConfig(path string) (config ENVConfig, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}

func FetchWeather(c *fiber.Ctx) error {
	// city := c.Query("city", "")
	data, err := fetchWeatherStack("Melbourne")
	// data, err := fetchOpenWeather(city)
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Currently we are unable to fetch weather data, please try again later.")
	}
	return c.JSON(data)
}

func fetchOpenWeather(city string) (map[string]float64, error) {
	var resp *http.Response
	var err error
	city = strings.ToLower(city)
	config, err := LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	resp, err = http.Get("http://api.openweathermap.org/data/2.5/weather?appid=" + config.OpenWeatherAPIKey + "&q=" + city + ",AU")
	if err != nil {
		log.Default().Println("Openweather api could not be connected")
		return nil, err
	}
	defer resp.Body.Close()

	var apiResponse OpenWeatherResponse
	json.NewDecoder(resp.Body).Decode(&apiResponse)

	return map[string]float64{
		"wind_speed":          apiResponse.Wind.Speed,
		"temperature_degrees": apiResponse.Main.Temp - 273,
	}, nil
}

func fetchWeatherStack(city string) (map[string]float64, error) {
	var resp *http.Response
	var err error
	config, err := LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	resp, err = http.Get("http://api.weatherstack.com/current?access_key=" + config.WeatherStackAPIKey + "&query=" + city)
	if err != nil {
		log.Default().Println("Weatherstack api could not be connected")
		fetchOpenWeather(city)
	}
	defer resp.Body.Close()

	var apiResponse WeatherStackResponse
	json.NewDecoder(resp.Body).Decode(&apiResponse)

	return map[string]float64{
		"wind_speed":          apiResponse.Current.WindSpeed,
		"temperature_degrees": apiResponse.Current.Temperature,
	}, nil
}
