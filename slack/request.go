package slack

import (
	"net/url"

	"github.com/ArjenSchwarz/igor/config"
)

type SlackRequest struct {
	Token       string
	TeamId      string
	TeamDomain  string
	ChannelId   string
	ChannelName string
	UserId      string
	UserName    string
	Command     string
	Text        string
	ResponseUrl string
}

func LoadRequestFromQuery(query string) SlackRequest {
	parsedQuery, _ := url.ParseQuery(query)
	request := SlackRequest{}
	request.Token = parsedQuery.Get("token")
	request.TeamId = parsedQuery.Get("team_id")
	request.TeamDomain = parsedQuery.Get("team_domain")
	request.ChannelId = parsedQuery.Get("channel_id")
	request.ChannelName = parsedQuery.Get("channel_name")
	request.UserId = parsedQuery.Get("user_id")
	request.UserName = parsedQuery.Get("user_name")
	request.Command = parsedQuery.Get("command")
	request.Text = parsedQuery.Get("text")
	request.ResponseUrl = parsedQuery.Get("response_url")
	return request
}

func (s *SlackRequest) Validate(config config.Config) bool {
	return s.Token == config.Token
}
