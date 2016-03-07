package plugins

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"gopkg.in/yaml.v2"

	"github.com/ArjenSchwarz/igor/config"
	"github.com/ArjenSchwarz/igor/helpers"
	"github.com/ArjenSchwarz/igor/slack"
)

// Weather instantiates a WeatherPlugin
func Weather() WeatherPlugin {
	pluginName := "weather"
	pluginConfig := ParseWeatherConfig()
	description := fmt.Sprintf("Igor provides weather information for the city you specify. If no city is specified, the default city (%s) is used.", pluginConfig.DefaultCity)
	plugin := WeatherPlugin{
		name:        pluginName,
		Source:      "http://api.openweathermap.org/data/2.5/",
		description: description,
		Config:      pluginConfig,
	}
	return plugin
}

// Describe describes the functionalities offered by the WeatherPlugin
func (WeatherPlugin) Describe() map[string]string {
	descriptions := make(map[string]string)
	descriptions["weather [city]"] = "Show the current weather in the city provided as argument"
	descriptions["forecast [city]"] = "Shows a 7 day forecast for the city provided as argument"
	return descriptions
}

// Work makes the WeatherPlugin run its commands
func (w WeatherPlugin) Work(request slack.SlackRequest) (slack.SlackResponse, error) {
	response := slack.SlackResponse{}
	if len(request.Text) >= 7 && request.Text[:7] == "weather" {
		response, err := w.handleWeather(request)
		return response, err
	} else if len(request.Text) >= 8 && request.Text[:8] == "forecast" {
		response, err := w.handleForecast(request)
		return response, err
	}

	return response, errors.New("No Match")
}

// handleWeather handles a request for the current Weather
func (w *WeatherPlugin) handleWeather(request slack.SlackRequest) (slack.SlackResponse, error) {
	var city string
	if len(request.Text) > 8 {
		city = request.Text[8:]
	} else {
		city = w.Config.DefaultCity
	}
	city = url.QueryEscape(city)
	response := slack.SlackResponse{}
	url := fmt.Sprintf("%sfind?APPID=%s&q=%s&units=%s", w.Source, w.Config.ApiToken, city, w.Config.Units)
	resp, err := http.Get(url)
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()
	parsedResult := WeatherResponse{}

	if err := json.NewDecoder(resp.Body).Decode(&parsedResult); err != nil {
		return response, err
	}
	response.Text = "Your weather request"
	for _, record := range parsedResult.List {
		attach := slack.Attachment{}
		attach.Title = fmt.Sprintf("%s, %s (%s)", record.Name, record.Sys.Country, helpers.RoughDay(record.Date))
		attach.ThumbUrl = weatherIconUrl(record.Weather[0].Icon)
		attach.Text = record.Weather[0].Desc
		tempField := slack.Field{}
		tempField.Title = "Temp"
		tempField.Value = formatTemp(record.Main.Temp, w.Config.Units)
		tempField.Short = true
		attach.AddField(tempField)
		windField := slack.Field{}
		windField.Title = "Wind"
		windField.Value = formatWind(record.Wind.Speed, w.Config.Units)
		windField.Short = true
		attach.AddField(windField)
		humField := slack.Field{}
		humField.Title = "Humidity"
		humField.Value = strconv.FormatInt(record.Main.Humidity, 10) + "%"
		humField.Short = true
		response.AddAttachment(attach)
	}

	return response, nil
}

// handleForecast handles the request for a forecast
func (w *WeatherPlugin) handleForecast(request slack.SlackRequest) (slack.SlackResponse, error) {
	var city string
	if len(request.Text) > 9 {
		city = request.Text[9:]
	} else {
		city = w.Config.DefaultCity
	}
	city = url.QueryEscape(city)
	response := slack.SlackResponse{}
	url := fmt.Sprintf("%sforecast/daily?APPID=%s&q=%s&units=%s", w.Source, w.Config.ApiToken, city, w.Config.Units)
	resp, err := http.Get(url)
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()
	parsedResult := ForecastResponse{}

	if err := json.NewDecoder(resp.Body).Decode(&parsedResult); err != nil {
		return response, err
	}
	response.Text = "Your forecast request"
	for _, record := range parsedResult.List {
		attach := slack.Attachment{}
		attach.Title = fmt.Sprintf("%s, %s (%s)", parsedResult.City.Name, parsedResult.City.Country, helpers.RoughDay(record.Date))
		attach.ThumbUrl = weatherIconUrl(record.Weather[0].Icon)
		attach.Text = record.Weather[0].Desc
		mintempField := slack.Field{}
		mintempField.Title = "Min Temp"
		mintempField.Value = formatTemp(record.Temp.Min, w.Config.Units)
		mintempField.Short = true
		attach.AddField(mintempField)
		maxtempField := slack.Field{}
		maxtempField.Title = "Max Temp"
		maxtempField.Value = formatTemp(record.Temp.Max, w.Config.Units)
		maxtempField.Short = true
		attach.AddField(maxtempField)
		windField := slack.Field{}
		windField.Title = "Wind"
		windField.Value = formatWind(record.Windspeed, w.Config.Units)
		windField.Short = true
		attach.AddField(windField)
		humField := slack.Field{}
		humField.Title = "Humidity"
		humField.Value = strconv.FormatInt(record.Humidity, 10) + "%"
		humField.Short = true
		attach.AddField(humField)
		response.AddAttachment(attach)
	}

	return response, nil
}

func (p WeatherPlugin) Description() string {
	return p.description
}
func (p WeatherPlugin) Name() string {
	return p.name
}

// ParseWeatherConfig collects the config as defined in the config file for
// the weather plugin
func ParseWeatherConfig() WeatherConfig {
	configFile := config.GetConfigFile()

	config := struct {
		Weather map[string]string `yaml:"weather"`
	}{}

	err := yaml.Unmarshal(configFile, &config)
	if err != nil {
		panic(err)
	}
	weather := WeatherConfig{Units: "metric"}
	value, ok := config.Weather["default_city"]
	if ok {
		weather.DefaultCity = value
	}
	value, ok = config.Weather["api_token"]
	if ok {
		weather.ApiToken = value
	}
	value, ok = config.Weather["units"]
	if ok {
		weather.Units = value
	}
	return weather
}

// weatherIconUrl returns the image location for a weather icon
// based on the code provided
func weatherIconUrl(code string) string {
	return "http://openweathermap.org/img/w/" + code + ".png"
}

// formatTemp formats the temperature by rounding it and adding the unit type
func formatTemp(temp float64, units string) string {
	var value string
	switch units {
	case "metric":
		value = "C"
	case "imperial":
		value = "F"
	}
	return fmt.Sprintf("%s %s", strconv.FormatFloat(temp, 'f', 0, 64), value)
}

// formatWind formats the wind by rounding it and adding the unit type
func formatWind(speed float64, units string) string {
	var value string
	switch units {
	case "metric":
		value = "km/h"
	case "imperial":
		value = "mph"
	}
	return fmt.Sprintf("%s %s", strconv.FormatFloat(speed, 'f', 0, 64), value)
}

type (
	WeatherResponse struct {
		Message string `json:"message"`
		List    []List `json:"list"`
	}

	List struct {
		Name string   `json:"name"`
		Main MainList `json:"main"`
		Wind struct {
			Speed float64 `json:"speed"`
		} `json:"wind"`
		Sys struct {
			Country string `json:"country"`
		} `json:"sys"`
		Weather []Wthr `json:"weather"`
		Date    int64  `json:"dt"`
	}

	Wthr struct {
		Main string `json:"main"`
		Desc string `json:"description"`
		Icon string `json:"icon"`
	}

	MainList struct {
		Temp     float64 `json:"temp"`
		Humidity int64   `json:"humidity"`
	}

	ForecastResponse struct {
		City City           `json:"city"`
		List []ForecastList `json:"list"`
	}

	City struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	}

	ForecastList struct {
		Temp      TempList `json:"temp"`
		Weather   []Wthr   `json:"weather"`
		Date      int64    `json:"dt"`
		Windspeed float64  `json:"speed"`
		Humidity  int64    `json:"humidity"`
	}

	TempList struct {
		Min float64 `json:"min"`
		Max float64 `json:"max"`
	}

	WeatherPlugin struct {
		name        string
		description string
		Source      string
		Config      WeatherConfig
	}

	WeatherConfig struct {
		DefaultCity string
		ApiToken    string
		Units       string
	}
)
