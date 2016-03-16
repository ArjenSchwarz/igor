package plugins

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/ArjenSchwarz/igor/config"
	"github.com/ArjenSchwarz/igor/helpers"
	"github.com/ArjenSchwarz/igor/slack"
)

// WeatherPlugin provides weather information for the city you specify
type WeatherPlugin struct {
	name        string
	description string
	Source      string
	Config      weatherConfig
}

// Weather instantiates a WeatherPlugin
func Weather() IgorPlugin {
	pluginName := "weather"
	pluginConfig := parseWeatherConfig()
	description := fmt.Sprintf("Igor provides weather information for the city you specify. If no city is specified, the default city (%s) is used.", pluginConfig.DefaultCity)
	plugin := WeatherPlugin{
		name:        pluginName,
		Source:      "http://api.openweathermap.org/data/2.5/",
		description: description,
		Config:      pluginConfig,
	}
	return plugin
}

// Describe provides the triggers WeatherPlugin can handle
func (WeatherPlugin) Describe() map[string]string {
	descriptions := make(map[string]string)
	descriptions["weather [city]"] = "Show the current weather in the city provided as argument"
	descriptions["forecast [city]"] = "Shows a 7 day forecast for the city provided as argument"
	return descriptions
}

// Work parses the request and ensures a request comes through if any triggers
// are matched. Handled triggers:
//
// * weather
// * forecast
func (plugin WeatherPlugin) Work(request slack.Request) (slack.Response, error) {
	response := slack.Response{}
	if len(request.Text) >= 7 && request.Text[:7] == "weather" {
		response, err := plugin.handleWeather(request)
		return response, err
	} else if len(request.Text) >= 8 && request.Text[:8] == "forecast" {
		response, err := plugin.handleForecast(request)
		return response, err
	}

	return response, errors.New("No Match")
}

// handleWeather handles a request for the current Weather
func (plugin *WeatherPlugin) handleWeather(request slack.Request) (slack.Response, error) {
	var city string
	if len(request.Text) > 8 {
		city = request.Text[8:]
	} else {
		city = plugin.Config.DefaultCity
	}
	city = url.QueryEscape(city)
	response := slack.Response{}
	url := fmt.Sprintf("%sfind?APPID=%s&q=%s&units=%s",
		plugin.Source, plugin.Config.APIToken, city, plugin.Config.Units)
	resp, err := http.Get(url)
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()
	parsedResult := weatherResponse{}

	if err := json.NewDecoder(resp.Body).Decode(&parsedResult); err != nil {
		return response, err
	}
	response.Text = "Your weather request"
	for _, record := range parsedResult.List {
		attach := slack.Attachment{}
		attach.Title = fmt.Sprintf("%s, %s (%s)",
			record.Name,
			record.Sys.Country,
			helpers.RoughDay(record.Date))
		attach.ThumbURL = weatherIconURL(record.Weather[0].Icon)
		attach.Text = record.Weather[0].Desc
		tempField := slack.Field{}
		tempField.Title = "Temp"
		tempField.Value = formatTemp(record.Main.Temp, plugin.Config.Units)
		tempField.Short = true
		attach.AddField(tempField)
		windField := slack.Field{}
		windField.Title = "Wind"
		windField.Value = formatWind(record.Wind.Speed, plugin.Config.Units)
		windField.Short = true
		attach.AddField(windField)
		humField := slack.Field{}
		humField.Title = "Humidity"
		humField.Value = strconv.FormatInt(record.Main.Humidity, 10) + "%"
		humField.Short = true
		attach.AddField(humField)
		response.AddAttachment(attach)
	}

	return response, nil
}

// handleForecast handles the request for a forecast
func (plugin *WeatherPlugin) handleForecast(request slack.Request) (slack.Response, error) {
	var city string
	if len(request.Text) > 9 {
		city = request.Text[9:]
	} else {
		city = plugin.Config.DefaultCity
	}
	city = url.QueryEscape(city)
	response := slack.Response{}
	url := fmt.Sprintf("%sforecast/daily?APPID=%s&q=%s&units=%s",
		plugin.Source,
		plugin.Config.APIToken,
		city,
		plugin.Config.Units)
	resp, err := http.Get(url)
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()
	parsedResult := forecastResponse{}

	if err := json.NewDecoder(resp.Body).Decode(&parsedResult); err != nil {
		return response, err
	}
	response.Text = "Your forecast request"
	for _, record := range parsedResult.List {
		attach := slack.Attachment{}
		attach.Title = fmt.Sprintf("%s, %s (%s)",
			parsedResult.City.Name,
			parsedResult.City.Country,
			helpers.RoughDay(record.Date))
		attach.ThumbURL = weatherIconURL(record.Weather[0].Icon)
		attach.Text = record.Weather[0].Desc
		mintempField := slack.Field{}
		mintempField.Title = "Min Temp"
		mintempField.Value = formatTemp(record.Temp.Min, plugin.Config.Units)
		mintempField.Short = true
		attach.AddField(mintempField)
		maxtempField := slack.Field{}
		maxtempField.Title = "Max Temp"
		maxtempField.Value = formatTemp(record.Temp.Max, plugin.Config.Units)
		maxtempField.Short = true
		attach.AddField(maxtempField)
		windField := slack.Field{}
		windField.Title = "Wind"
		windField.Value = formatWind(record.Windspeed, plugin.Config.Units)
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

// Description returns a global description of the plugin
func (plugin WeatherPlugin) Description() string {
	return plugin.description
}

// Name returns the name of the plugin
func (plugin WeatherPlugin) Name() string {
	return plugin.name
}

// parseWeatherConfig collects the config as defined in the config file for
// the weather plugin
func parseWeatherConfig() weatherConfig {
	pluginConfig := struct {
		Weather map[string]string `yaml:"weather"`
	}{}
    err := config.ParsePluginConfig(&pluginConfig)
    if err != nil {
        panic(err)
    }

	weather := weatherConfig{Units: "metric"}
	value, ok := pluginConfig.Weather["default_city"]
	if ok {
		weather.DefaultCity = value
	}
	value, ok = pluginConfig.Weather["api_token"]
	if ok {
		weather.APIToken = value
	}
	value, ok = pluginConfig.Weather["units"]
	if ok {
		weather.Units = value
	}
	return weather
}

// weatherIconURL returns the image location for a weather icon
// based on the code provided
func weatherIconURL(code string) string {
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
	weatherResponse struct {
		Message string `json:"message"`
		List    []list `json:"list"`
	}

	list struct {
		Name string   `json:"name"`
		Main mainList `json:"main"`
		Wind struct {
			Speed float64 `json:"speed"`
		} `json:"wind"`
		Sys struct {
			Country string `json:"country"`
		} `json:"sys"`
		Weather []wthr `json:"weather"`
		Date    int64  `json:"dt"`
	}

	wthr struct {
		Main string `json:"main"`
		Desc string `json:"description"`
		Icon string `json:"icon"`
	}

	mainList struct {
		Temp     float64 `json:"temp"`
		Humidity int64   `json:"humidity"`
	}

	forecastResponse struct {
		City city           `json:"city"`
		List []forecastList `json:"list"`
	}

	city struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	}

	forecastList struct {
		Temp      tempList `json:"temp"`
		Weather   []wthr   `json:"weather"`
		Date      int64    `json:"dt"`
		Windspeed float64  `json:"speed"`
		Humidity  int64    `json:"humidity"`
	}

	tempList struct {
		Min float64 `json:"min"`
		Max float64 `json:"max"`
	}

	weatherConfig struct {
		DefaultCity string
		APIToken    string
		Units       string
	}
)
