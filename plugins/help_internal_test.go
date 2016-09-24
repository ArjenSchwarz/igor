package plugins

import "testing"

func TestName(t *testing.T) {
	expectedResult := "testname"
	plugin := HelpPlugin{name: expectedResult}
	if plugin.Name() != expectedResult {
		t.Error("Name method not working correctly")
	}
}
