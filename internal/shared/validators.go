package shared

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/itchyny/gojq"
)

func ParseBool(key string) func(string) (any, error) {
	return func(s string) (any, error) {
		switch s {
		case "true":
			return true, nil
		case "false":
			return false, nil
		default:
			return nil, fmt.Errorf("%s must be 'true' or 'false'", key)
		}
	}
}

func ParseIntInRange(key string, min, max int) func(string) (any, error) {
	return func(s string) (any, error) {
		n, err := strconv.Atoi(s)
		if err != nil || n < min || n > max {
			return nil, fmt.Errorf("%s must be an integer between %d and %d", key, min, max)
		}
		return n, nil
	}
}

func ParseEnum(key string, allowed ...string) func(string) (any, error) {
	return func(s string) (any, error) {
		for _, v := range allowed {
			if s == v {
				return s, nil
			}
		}
		return nil, fmt.Errorf("%s must be one of: %v", key, allowed)
	}
}

func StringPassThrough(_ string) func(string) (any, error) {
	return func(s string) (any, error) {
		return s, nil
	}
}

func ValidateRegex(key string) func(string) (any, error) {
	return func(s string) (any, error) {
		if _, err := regexp.Compile(s); err != nil {
			return nil, fmt.Errorf("%s is not a valid regex: %w", key, err)
		}
		return s, nil
	}
}

func ValidateJQ(key string) func(string) (any, error) {
	return func(s string) (any, error) {
		if _, err := gojq.Parse(s); err != nil {
			return nil, fmt.Errorf("%s is not a valid jq expression: %w", key, err)
		}
		return s, nil
	}
}

func SuggestBool() []string {
	return []string{"true", "false"}
}
