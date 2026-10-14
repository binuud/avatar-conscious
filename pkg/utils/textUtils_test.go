package utils

import (
	"testing"
)

var sample_text = "this is a a, ball demo text with some   needle works elephant  man elephant sample superman"
var needles_present = [...]string{"ball", "elephant", "man"}
var needles_some_fail = [...]string{"ball", "elephant", "super"}
var needles_all_fail = [...]string{"one", "two", "three"}

func Test_ContainsAllWords(t *testing.T) {
	t.Logf("Testing utils.ContainsAllWords(hay, needles)")
	if ContainsAllWords(sample_text, needles_present[:]) {
		t.Log("words found")
	} else {
		t.Errorf("words not found : %v \n", needles_present)
	}
	if ContainsAllWords(sample_text, needles_some_fail[:]) {
		t.Errorf("words should have failed search: %v \n", needles_some_fail)
	} else {
		t.Logf("words have failed search as required: %v \n", needles_some_fail)
	}
	if ContainsAllWords(sample_text, needles_some_fail[:]) {
		t.Logf("words should have failed search: %v \n", needles_all_fail)
	} else {
		t.Logf("words have failed search as required: %v \n", needles_all_fail)
	}
}
