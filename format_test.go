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

	// string
	got = sanitizeValue("hello world")
	expected = "'hello world'"
	if got != expected {
		t.Errorf("Sanitized string is incorrect. Expected %s got %s", expected, got)
		t.Fail()
	}

	// *string
	s := "hello world"
	got = sanitizeValue(&s)
	expected = "'hello world'"
	if got != expected {
		t.Errorf("Sanitized string is incorrect. Expected %s got %s", expected, got)
		t.Fail()
	}

	// nil *string
	var n *string
	got = sanitizeValue(n)
	expected = ""
	if got != expected {
		t.Errorf("Sanitized string is incorrect. Expected %s got %s", expected, got)
		t.Fail()
	}

	// int
	got = sanitizeValue(1)
	expected = "1"
	if got != expected {
		t.Errorf("Sanitized string is incorrect. Expected %s got %s", expected, got)
		t.Fail()
	}

	// float
	got = sanitizeValue(1.25)
	expected = "1.25"
	if got != expected {
		t.Errorf("Sanitized string is incorrect. Expected %s got %s", expected, got)
		t.Fail()
	}

	// -float
	got = sanitizeValue(-3.14)
	expected = "-3.14"
	if got != expected {
		t.Errorf("Sanitized string is incorrect. Expected %s got %s", expected, got)
		t.Fail()
	}

	// bool
	got = sanitizeValue(false)
	expected = "0"
	if got != expected {
		t.Errorf("Sanitized string is incorrect. Expected %s got %s", expected, got)
		t.Fail()
	}

	// bytes
	got = sanitizeValue([]byte("hello world"))
	expected = "'68656c6c6f20776f726c64'"
	if got != expected {
		t.Errorf("Sanitized string is incorrect. Expected %s got %s", expected, got)
		t.Fail()
	}

	// escape string
	got = sanitizeValue("wayne's world")
	expected = "'wayne''s world'"
	if got != expected {
		t.Errorf("Sanitized string is incorrect. Expected %s got %s", expected, got)
		t.Fail()
	}
}
