package plugins

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/ArjenSchwarz/igor/config"
	"github.com/ArjenSchwarz/igor/helpers"
	"github.com/ArjenSchwarz/igor/slack"
)

// WeatherPlugin provides weather information for the city you specify
type WeatherPlugin struct {
	name        string
	description string
	Source      string
	config      weatherConfig
	request     slack.Request
}

// Config returns the plugin configuration
func (plugin WeatherPlugin) Config() IgorConfig {
	return plugin.config
}

func (config weatherConfig) Languages() map[string]config.LanguagePluginDetails {
	return config.languages
}

func (config weatherConfig) ChosenLanguage() string {
	return config.chosenLanguage
}

// Weather instantiates a WeatherPlugin
func Weather(request slack.Request) (IgorPlugin, error) {
	pluginName := "weather"
	pluginConfig, err := parseWeatherConfig()
	if err != nil {
		return WeatherPlugin{}, err
	}
	plugin := WeatherPlugin{
		name:    pluginName,
		Source:  "http://api.openweathermap.org/data/2.5/",
		config:  pluginConfig,
		request: request,
	}
	return plugin, nil
}

// Describe provides the triggers WeatherPlugin can handle
func (plugin WeatherPlugin) Describe(language string) map[string]string {

	descriptions := make(map[string]string)
	for _, values := range getAllCommands(plugin, language) {
		descriptions[values.Command] = values.Description
	}
	return descriptions
}

// Work parses the request and ensures a request comes through if any triggers
// are matched. Handled triggers:
//
// * weather
// * forecast
func (plugin WeatherPlugin) Work() (slack.Response, error) {
	response := slack.Response{}
	message, language := getCommandName(plugin)
	plugin.config.chosenLanguage = language
	switch message {
	case "weather":
		return plugin.handleWeather()
	case "forecast":
		return plugin.handleForecast()
	}

	return response, CreateNoMatchError("Nothing found")
}

// handleWeather handles a request for the current Weather
func (plugin *WeatherPlugin) handleWeather() (slack.Response, error) {
	var city string
	parts := strings.Split(plugin.Message(), " ")
	if len(parts) > 1 {
		city = strings.Replace(plugin.Message(), parts[0], "", 1)
	} else {
		city = plugin.config.determineDefaultWeatherCity(plugin.request)
	}
	city = url.QueryEscape(city)
	if isSpecialWeather(city) {
		return getSpecialWeather(city)
	}
	response := slack.Response{}
	url := fmt.Sprintf("%sfind?APPID=%s&q=%s&units=%s",
		plugin.Source, plugin.config.APIToken, city, plugin.config.Units)
	resp, err := http.Get(url)
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()
	parsedResult := weatherResponse{}

	if err := json.NewDecoder(resp.Body).Decode(&parsedResult); err != nil {
		return response, err
	}
	commandDetails := getCommandDetails(plugin, "weather")
	response.Text = commandDetails.Texts["response_text"]
	for _, record := range parsedResult.List {
		attach := slack.Attachment{}
		attach.Title = fmt.Sprintf("%s, %s (%s)",
			record.Name,
			record.Sys.Country,
			helpers.RoughDay(record.Date))
		attach.ThumbURL = weatherIconURL(record.Weather[0].Icon)
		attach.Text = record.Weather[0].Desc
		tempField := slack.Field{}
		tempField.Title = commandDetails.Texts["temperature"]
		tempField.Value = formatTemp(record.Main.Temp, plugin.config.Units)
		tempField.Short = true
		attach.AddField(tempField)
		windField := slack.Field{}
		windField.Title = commandDetails.Texts["wind"]
		windField.Value = formatWind(record.Wind.Speed, plugin.config.Units)
		windField.Short = true
		attach.AddField(windField)
		humField := slack.Field{}
		humField.Title = commandDetails.Texts["humidity"]
		humField.Value = strconv.FormatInt(record.Main.Humidity, 10) + "%"
		humField.Short = true
		attach.AddField(humField)
		response.AddAttachment(attach)
	}

	return response, nil
}

// handleForecast handles the request for a forecast
func (plugin *WeatherPlugin) handleForecast() (slack.Response, error) {
	var city string
	parts := strings.Split(plugin.Message(), " ")
	if len(parts) > 1 {
		city = strings.Replace(plugin.Message(), parts[0], "", 1)
	} else {
		city = plugin.config.determineDefaultWeatherCity(plugin.request)
	}
	city = url.QueryEscape(city)
	response := slack.Response{}
	url := fmt.Sprintf("%sforecast/daily?APPID=%s&q=%s&units=%s",
		plugin.Source,
		plugin.config.APIToken,
		city,
		plugin.config.Units)
	resp, err := http.Get(url)
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()
	parsedResult := forecastResponse{}

	if err := json.NewDecoder(resp.Body).Decode(&parsedResult); err != nil {
		return response, err
	}
	commandDetails := getCommandDetails(plugin, "forecast")
	response.Text = commandDetails.Texts["response_text"]
	for _, record := range parsedResult.List {
		attach := slack.Attachment{}
		attach.Title = fmt.Sprintf("%s, %s (%s)",
			parsedResult.City.Name,
			parsedResult.City.Country,
			helpers.RoughDay(record.Date))
		attach.ThumbURL = weatherIconURL(record.Weather[0].Icon)
		attach.Text = record.Weather[0].Desc
		mintempField := slack.Field{}
		mintempField.Title = commandDetails.Texts["min_temperature"]
		mintempField.Value = formatTemp(record.Temp.Min, plugin.config.Units)
		mintempField.Short = true
		attach.AddField(mintempField)
		maxtempField := slack.Field{}
		maxtempField.Title = commandDetails.Texts["max_temperature"]
		maxtempField.Value = formatTemp(record.Temp.Max, plugin.config.Units)
		maxtempField.Short = true
		attach.AddField(maxtempField)
		windField := slack.Field{}
		windField.Title = commandDetails.Texts["wind"]
		windField.Value = formatWind(record.Windspeed, plugin.config.Units)
		windField.Short = true
		attach.AddField(windField)
		humField := slack.Field{}
		humField.Title = commandDetails.Texts["humidity"]
		humField.Value = strconv.FormatInt(record.Humidity, 10) + "%"
		humField.Short = true
		attach.AddField(humField)
		response.AddAttachment(attach)
	}

	return response, nil
}

// determineDefaultWeatherCity checks if there are defaults for specific channels
// and returns those.
func (config weatherConfig) determineDefaultWeatherCity(request slack.Request) string {
	if config.ChannelCity != nil {
		if val, ok := config.ChannelCity[request.ChannelID]; ok {
			return val
		}
		if val, ok := config.ChannelCity[request.ChannelName]; ok {
			return val
		}
	}
	return config.DefaultCity
}

// Description returns a global description of the plugin
func (plugin WeatherPlugin) Description(language string) string {
	descrString := strings.Replace(getDescriptionText(plugin, language), "[replace]", "%s", -1)
	return fmt.Sprintf(descrString, plugin.config.determineDefaultWeatherCity(plugin.request))
}

// Name returns the name of the plugin
func (plugin WeatherPlugin) Name() string {
	return plugin.name
}

// Message returns the request sent
func (plugin WeatherPlugin) Message() string {
	return strings.ToLower(plugin.request.Text)
}

// parseWeatherConfig collects the config as defined in the config file for
// the weather plugin
func parseWeatherConfig() (weatherConfig, error) {
	pluginConfig := struct {
		Weather weatherConfig
	}{}

	err := config.ParseConfig(&pluginConfig)
	if err != nil {
		return pluginConfig.Weather, err
	}
	pluginConfig.Weather.languages = getPluginLanguages("weather")

	if pluginConfig.Weather.Units == "" {
		pluginConfig.Weather.Units = "metric"
	}

	return pluginConfig.Weather, nil
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

func isSpecialWeather(location string) bool {
	locations := getSpecialWeatherMap()
	if _, ok := locations[location]; ok {
		return true
	}
	return false
}

func getSpecialWeather(location string) (slack.Response, error) {
	response := slack.Response{Text: "Your weather request"}

	locations := getSpecialWeatherMap()
	if val, ok := locations[location]; ok {
		response.AddAttachment(val)
		return response, nil
	}
	return response, CreateNoMatchError("No special location")
}

func getSpecialWeatherMap() map[string]slack.Attachment {
	plugins := make(map[string]slack.Attachment)
	hothResponse := slack.Attachment{
		Title:    "Hoth, A galaxy far far away (A long time ago)",
		Text:     "Tauntaun freezing cold",
		ThumbURL: weatherIconURL("13d"),
	}
	plugins["hoth"] = hothResponse
	tatResponse := slack.Attachment{
		Title:    "Tatooine, A galaxy far far away (A long time ago)",
		Text:     "So hot milk turns blue",
		ThumbURL: weatherIconURL("01d"),
	}
	plugins["tatooine"] = tatResponse

	yakkuResponse := slack.Attachment{
		Title:    "Yakku, A galaxy far far away (A long time ago)",
		Text:     "Hot enough to make Star Destroyers crash",
		ThumbURL: weatherIconURL("01d"),
	}
	plugins["yakku"] = yakkuResponse

	dagobahResponse := slack.Attachment{
		Title:    "Dagobah, A galaxy far far away (A long time ago)",
		Text:     "Hot and humid with no visibility. Not recommended to fly",
		ThumbURL: weatherIconURL("09d"),
	}
	plugins["dagobah"] = dagobahResponse
	return plugins
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
		DefaultCity    string
		APIToken       string
		Units          string
		ChannelCity    map[string]string
		languages      map[string]config.LanguagePluginDetails
		chosenLanguage string
	}
)
