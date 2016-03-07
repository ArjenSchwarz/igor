// Package slack provides all the Slack specific code for Igor
package slack

import (
	"strings"
)

type SlackResponse struct {
	Text         string       `json:"text"`
	ResponseType string       `json:"response_type,omitempty"`
	Attachments  []Attachment `json:"attachments,omitempty"`
	UnfurlLinks  bool         `json:"unfurl_links,omitempty"`
	UnfurlMedia  bool         `json:"unfurl_media,omitempty"`
	Markdown     bool         `json:"mrkdwn,omitempty"`
	Username     string       `json:"username,omitempty"`
	IconEmoji    string       `json:"icon_emoji,omitempty"`
}

type Attachment struct {
	Title      string   `json:"title,omitempty"`
	TitleLink  string   `json:"title_link,omitempty"`
	Text       string   `json:"text,omitempty"`
	PreText    string   `json:"pretext,omitempty"`
	Color      string   `json:"color,omitempty"`
	Markdown   []string `json:"mrkdwn_in,omitempty"`
	Fallback   string   `json:"fallback,omitempty"`
	AuthorName string   `json:"author_name,omitempty"`
	AuthorLink string   `json:"author_link,omitempty"`
	AuthorIcon string   `json:"author_icon,omitempty"`
	ImageUrl   string   `json:"image_url,omitempty"`
	ThumbUrl   string   `json:"thumb_url,omitempty"`
	Fields     []Field  `json:"fields,omitempty"`
}

type Field struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short,omitempty"`
}

//TODO validation
func (a *Attachment) EnableMarkdownFor(value string) {
	var markin []string
	if a.Markdown != nil {
		markin = a.Markdown
	}
	a.Markdown = append(markin, value)
}

//TODO validation
func (s *SlackResponse) AddAttachment(a Attachment) {
	var attachments []Attachment
	if s.Attachments != nil {
		attachments = s.Attachments
	}
	s.Attachments = append(attachments, a)
}

func (a *Attachment) AddField(f Field) {
	var fields []Field
	if a.Fields != nil {
		fields = a.Fields
	}
	a.Fields = append(fields, f)
}

func (s *SlackResponse) SetPublic() {
	s.ResponseType = "in_channel"
}

func (s *SlackResponse) IsPublic() bool {
	return s.ResponseType == "in_channel"
}

func NothingFoundResponse(request SlackRequest) SlackResponse {
	response := SlackResponse{}
	response.Text = "Our apologies. No Igor was able to handle your request."
	attach := Attachment{}
	attach.Color = "danger"
	attach.Text = "You tried to look for *" + request.Command + " " + request.Text + "\n"
	attach.Text += "Please try *" + request.Command + " help* to see which Igors are available"
	attach.EnableMarkdownFor("text")
	response.AddAttachment(attach)
	return response
}

func ValidationErrorResponse() SlackResponse {
	response := SlackResponse{}
	response.Text = "Invalid token."
	return response
}

// EscapeString escapes any values as demanded by Slack
// This means it HTML escapes '&', '<', and '>'
// It doesn't double escape. If a string is already escaped it won't do it
// again. Meaning, if you supply "&amp;" it will return "&amp;"
func EscapeString(toEscape string) string {
	toEscape = strings.Replace(toEscape, "&amp;", "&", -1)
	toEscape = strings.Replace(toEscape, "&lt;", "<", -1)
	toEscape = strings.Replace(toEscape, "&gt;", ">", -1)
	toEscape = strings.Replace(toEscape, "&", "&amp;", -1)
	toEscape = strings.Replace(toEscape, "<", "&lt;", -1)
	toEscape = strings.Replace(toEscape, ">", "&gt;", -1)
	return toEscape
}

// Escape escapes all values in the SlackResponse that need to be escaped
func (s *SlackResponse) Escape() {
	s.Text = EscapeString(s.Text)
	for _, attach := range s.Attachments {
		attach.Title = EscapeString(attach.Title)
		attach.Text = EscapeString(attach.Text)
		attach.PreText = EscapeString(attach.PreText)
		for _, field := range attach.Fields {
			field.Title = EscapeString(field.Title)
			field.Value = EscapeString(field.Value)
		}
	}
}
