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
	Color    string   `json:"color"`
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
