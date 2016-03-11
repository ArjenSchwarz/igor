// Package slack provides all the Slack specific code for Igor
package slack

import (
	"strings"
)

// Response contains the fields for returning a response message to Slack
type Response struct {
	Text         string       `json:"text"`
	ResponseType string       `json:"response_type,omitempty"`
	Attachments  []Attachment `json:"attachments,omitempty"`
	UnfurlLinks  bool         `json:"unfurl_links,omitempty"`
	UnfurlMedia  bool         `json:"unfurl_media,omitempty"`
	Markdown     bool         `json:"mrkdwn,omitempty"`
	Username     string       `json:"username,omitempty"`
	IconEmoji    string       `json:"icon_emoji,omitempty"`
}

// Attachment contains the fields for a slack attachment
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
	ImageURL   string   `json:"image_url,omitempty"`
	ThumbURL   string   `json:"thumb_url,omitempty"`
	Fields     []Field  `json:"fields,omitempty"`
}

// Field contains the fields for a field within a slack attachment
type Field struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short,omitempty"`
}

// EnableMarkdownFor enables Markdown for a part of the attachment
func (a *Attachment) EnableMarkdownFor(value string) {
	//TODO validation
	var markin []string
	if a.Markdown != nil {
		markin = a.Markdown
	}
	a.Markdown = append(markin, value)
}

// AddAttachment adds an attachment to the response
func (response *Response) AddAttachment(a Attachment) {
	//TODO validation
	var attachments []Attachment
	if response.Attachments != nil {
		attachments = response.Attachments
	}
	response.Attachments = append(attachments, a)
}

// AddField adds a field to the attachment
func (a *Attachment) AddField(f Field) {
	var fields []Field
	if a.Fields != nil {
		fields = a.Fields
	}
	a.Fields = append(fields, f)
}

// SetPublic configures the response to show up publicly
func (response *Response) SetPublic() {
	response.ResponseType = "in_channel"
}

// IsPublic returns whether the response is configured to show up publicly
func (response *Response) IsPublic() bool {
	return response.ResponseType == "in_channel"
}

// NothingFoundResponse is a specific response for when no matching trigger
// is found
func NothingFoundResponse(request Request) Response {
	response := Response{}
	response.Text = "Our apologies. No Igor was able to handle your request."
	attach := Attachment{}
	attach.Color = "danger"
	attach.Text = "You tried to look for *" + request.Command + " " + request.Text + "\n"
	attach.Text += "Please try *" + request.Command + " help* to see which Igors are available"
	attach.EnableMarkdownFor("text")
	response.AddAttachment(attach)
	return response
}

// ValidationErrorResponse is a specific response for when validation failed
func ValidationErrorResponse() Response {
	response := Response{}
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

// Escape escapes all values in the Response that need to be escaped
func (response *Response) Escape() {
	response.Text = EscapeString(response.Text)
	for _, attach := range response.Attachments {
		attach.Title = EscapeString(attach.Title)
		attach.Text = EscapeString(attach.Text)
		attach.PreText = EscapeString(attach.PreText)
		for _, field := range attach.Fields {
			field.Title = EscapeString(field.Title)
			field.Value = EscapeString(field.Value)
		}
	}
}
