package slack

type SlackResponse struct {
	Text         string       `json:"text"`
	ResponseType string       `json:"response_type,omitempty"`
	Attachments  []Attachment `json:"attachments,omitempty"`
}

type Attachment struct {
	Title    string   `json:"title,omitempty"`
	Text     string   `json:"text,omitempty"`
	PreText  string   `json:"pretext,omitempty"`
	Markdown []string `json:"mrkdwn_in,omitempty"`
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

func (s *SlackResponse) SetPublic() {
	s.ResponseType = "in_channel"
}
