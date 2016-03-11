// Package slack provides all the Slack specific code for Igor
package slack

import (
	"net/url"

	"github.com/ArjenSchwarz/igor/config"
)

// Request contains the information sent through the request from Slack
type Request struct {
	Token       string
	TeamID      string
	TeamDomain  string
	ChannelID   string
	ChannelName string
	UserID      string
	UserName    string
	Command     string
	Text        string
	ResponseURL string
}

// LoadRequestFromQuery translates the query string sent by Slack into a Request struct
func LoadRequestFromQuery(query string) Request {
	parsedQuery, _ := url.ParseQuery(query)
	request := Request{}
	request.Token = parsedQuery.Get("token")
	request.TeamID = parsedQuery.Get("team_id")
	request.TeamDomain = parsedQuery.Get("team_domain")
	request.ChannelID = parsedQuery.Get("channel_id")
	request.ChannelName = parsedQuery.Get("channel_name")
	request.UserID = parsedQuery.Get("user_id")
	request.UserName = parsedQuery.Get("user_name")
	request.Command = parsedQuery.Get("command")
	request.Text = parsedQuery.Get("text")
	request.ResponseURL = parsedQuery.Get("response_url")
	return request
}

// Validate ensures the request comes from the configured Slack team
func (request *Request) Validate(config config.Config) bool {
	return request.Token == config.Token
}
