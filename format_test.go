package sqli

import (
	"testing"
)

func Test_stripName(t *testing.T) {
	got := stripName("table`")
	expected := "table``"
	if got != expected {
		t.Errorf("Sanitized string is incorrect. Expected %s got %s", expected, got)
		t.Fail()
	}
}

func Test_sanitizeValue(t *testing.T) {
	var expected string
	var got string

	got = sanitizeValue("hello world")
	expected = "'hello world'"
	if got != expected {
		t.Errorf("Sanitized string is incorrect. Expected %s got %s", expected, got)
		t.Fail()
	}
	got = sanitizeValue(1)
	expected = "1"
	if got != expected {
		t.Errorf("Sanitized string is incorrect. Expected %s got %s", expected, got)
		t.Fail()
	}
	got = sanitizeValue(1.25)
	expected = "1.25"
	if got != expected {
		t.Errorf("Sanitized string is incorrect. Expected %s got %s", expected, got)
		t.Fail()
	}
	got = sanitizeValue(-3.14)
	expected = "-3.14"
	if got != expected {
		t.Errorf("Sanitized string is incorrect. Expected %s got %s", expected, got)
		t.Fail()
	}
	got = sanitizeValue(false)
	expected = "0"
	if got != expected {
		t.Errorf("Sanitized string is incorrect. Expected %s got %s", expected, got)
		t.Fail()
	}
	got = sanitizeValue([]byte("hello world"))
	expected = "'68656c6c6f20776f726c64'"
	if got != expected {
		t.Errorf("Sanitized string is incorrect. Expected %s got %s", expected, got)
		t.Fail()
	}
	got = sanitizeValue("wayne's world")
	expected = "'wayne''s world'"
	if got != expected {
		t.Errorf("Sanitized string is incorrect. Expected %s got %s", expected, got)
		t.Fail()
	}
}
