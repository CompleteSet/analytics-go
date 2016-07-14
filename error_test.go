package analytics

import "testing"

func TestConfigError(t *testing.T) {
	e := ConfigError{
		Reason: "testing",
		Field:  "Answer",
		Value:  42,
	}

	if s := e.Error(); s != "ConfigError: testing (Answer: 42)" {
		t.Error("invalid error message returned by config error:", s)
	}
}

func TestFieldError(t *testing.T) {
	e := FieldError{
		Type:  "testing.T",
		Name:  "Answer",
		Value: 42,
	}

	if s := e.Error(); s != "FieldError: invalid field value (testing.T.Answer: 42)" {
		t.Error("invalid error message returned by field error:", s)
	}
}
