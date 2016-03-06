package slack_test

import (
	"github.com/ArjenSchwarz/igor/slack"
	"testing"
)

var escapeTests = []struct {
	input    string
	expected string
}{
	{"&", "&amp;"},
	{"<", "&lt;"},
	{">", "&gt;"},
	{"test & <url>", "test &amp; &lt;url&gt;"},
	{"test && <url>><&", "test &amp;&amp; &lt;url&gt;&gt;&lt;&amp;"},
	{"&amp;", "&amp;"},
}

func TestEscapeString(t *testing.T) {
	for _, tt := range escapeTests {
		actual := slack.EscapeString(tt.input)
		if actual != tt.expected {
			t.Errorf("EscapeString(%v): expected %v, actual %v", tt.input, tt.expected, actual)
		}
	}
}
