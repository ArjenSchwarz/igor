package plugins

import (
	"testing"
)

func TestName(t *testing.T) {
	expectedResult := "testname"
	plugin := HelpPlugin{name: expectedResult}
	if plugin.Name() != expectedResult {
		t.Error("Name method not working correctly")
	}
}

func TestDescription(t *testing.T) {
	expectedResult := "testdescription"
	plugin := HelpPlugin{description: expectedResult}
	if plugin.Description() != expectedResult {
		t.Error("Description method not working correctly")
	}
}
