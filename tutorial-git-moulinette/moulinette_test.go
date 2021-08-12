package main

import (
	"testing"
)

func TestContains(t *testing.T) {
	var parameters = []struct {
	    content string
	    substring string
	    strict bool
	    expected bool
	}{
	    // Positive test
	    {"Jean-Edouard Dujardin hopla", "Jean-Edouard Dujardin", false, true},
	    // Check with - and and accent
	    {"Jean Edouard Dujardin hopla", "Jean-Edouard Dujardin", false, true},
	    {"Jean Edouard Dujardin hopla", "Jean-Ëdouard Dujardin", false, true},
	    // Strict
	    {"Jean Edouard Dujardin hopla", "Jean-Edouard Dujardin", true, false},
	    // Non-matching names
	    {"Jean Edouard Dujardin hopla", "Marc Edouard Dujardin", false, false},
	    {"Jean Edouard Dujardin hopla", "Marc-Edouard Dujardin", false, false},
	}
	for _, parameter := range parameters {
		result := Contains(parameter.content, parameter.substring, parameter.strict)
		if result != parameter.expected {
			t.Fatalf(`Contains("%s", "%s", %t) -> %t instead of %t`,
			parameter.content, parameter.substring, parameter.strict, result, parameter.expected)
		}
	}
}
