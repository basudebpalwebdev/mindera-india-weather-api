package handlers

import (
	"encoding/json"
	"log"
	"net/http"

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

// var Cache = new(map[string]float64)

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

// It fetches the data initially from WeatherStack api, if there is an error, it tries the same from OpenWeather api
// returns an error if the api is unavailable or the response doesn't have the required fields
// *** The requirement of returning last cached data in case of failure from both APIs has not been implemented ***
func FetchWeather(c *fiber.Ctx) error {
	// city := c.Query("city", "")
	config, err := LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	data, err := FetchWeatherStack("http://api.weatherstack.com/current?access_key=" + config.WeatherStackAPIKey + "&query=Melbourne")
	if err != nil {
		data, err = FetchOpenWeather("http://api.openweathermap.org/data/2.5/weather?appid=" + config.OpenWeatherAPIKey + "&q=melbourne,AU")
	}
	if err != nil {
		return c.JSON(nil)
	}
	// *Cache = data
	return c.JSON(data)
}

// It fetches the data from the OpenWeather api and extracts wind speed and temperature form it
// returns an error if the api is unavailable or the response doesn't have the required fields
func FetchOpenWeather(weatherApi string) (map[string]float64, error) {
	var resp *http.Response
	var err error
	resp, err = http.Get(weatherApi)
	if err != nil {
		log.Default().Println("Openweather api could not be connected")
		return nil, err
	}
	defer resp.Body.Close()

	var apiResponse OpenWeatherResponse
	err = json.NewDecoder(resp.Body).Decode(&apiResponse)

	if err != nil {
		return nil, err
	}

	return map[string]float64{
		"wind_speed":          apiResponse.Wind.Speed,
		"temperature_degrees": apiResponse.Main.Temp - 273,
	}, nil
}

// It fetches the data from the WeatherStack api and extracts wind speed and temperature form it
// returns an error if the api is unavailable or the response doesn't have the required fields
func FetchWeatherStack(weatherApi string) (map[string]float64, error) {
	var resp *http.Response
	var err error
	resp, err = http.Get(weatherApi)
	if err != nil {
		log.Default().Println("Weatherstack api could not be connected")
		return nil, err
	}
	defer resp.Body.Close()

	var apiResponse WeatherStackResponse
	err = json.NewDecoder(resp.Body).Decode(&apiResponse)

	if err != nil {
		return nil, err
	}

	return map[string]float64{
		"wind_speed":          apiResponse.Current.WindSpeed,
		"temperature_degrees": apiResponse.Current.Temperature,
	}, nil
}
