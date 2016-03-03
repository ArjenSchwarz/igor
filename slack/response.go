package slack

type SlackResponse struct {
	Text         string        `json:"text"`
	ResponseType string        `json:"response_type"`
	Attachments  []*Attachment `json:"attachments"`
}

type Attachment struct {
	Title   string `json:"title"`
	Text    string `json:"text"`
	PreText string `json:"pretext"`
}
