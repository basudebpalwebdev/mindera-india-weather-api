package handlers_test

import (
	"log"
	"net/http/httptest"
	"testing"

	"github.com/basudebpalwebdev/mindera-india-weather-api/api/v1/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestFetchWeather(t *testing.T) {
	testApi := fiber.New()
	resp, err := testApi.Test(httptest.NewRequest(fiber.MethodGet, "/v1/weather?city=melbourne", nil))
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestFetchOpenWeather(t *testing.T) {
	config, err := handlers.LoadConfig("./../../../")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	weatherData, err := handlers.FetchOpenWeather("http://api.openweathermap.org/data/2.5/weather?appid=" + config.OpenWeatherAPIKey + "&q=melbourne,AU")
	assert.NoError(t, err)
	assert.NotEmpty(t, weatherData)
	assert.Len(t, weatherData, 2)
	assert.NotNil(t, weatherData["wind_speed"])
	assert.NotNil(t, weatherData["temperature_degrees"])
}

func TestFetchFailureOpenWeather(t *testing.T) {
	weatherData, err := handlers.FetchOpenWeather("http://not.a.valid.url")
	assert.Error(t, err)
	assert.Empty(t, weatherData)
}

func TestFetchWeatherStack(t *testing.T) {
	config, err := handlers.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	weatherData, err := handlers.FetchWeatherStack("http://api.weatherstack.com/current?access_key=" + config.WeatherStackAPIKey + "&query=Melbourne")
	assert.NoError(t, err)
	assert.NotEmpty(t, weatherData)
	assert.Len(t, weatherData, 2)
	assert.NotNil(t, weatherData["wind_speed"])
	assert.NotNil(t, weatherData["temperature_degrees"])
}

func TestFetchFailureWeatherStack(t *testing.T) {
	weatherData, err := handlers.FetchWeatherStack("http://not.a.valid.url")
	assert.Error(t, err)
	assert.Empty(t, weatherData)
}
