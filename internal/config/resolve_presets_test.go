package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolveRegex(t *testing.T) {
	cfg := &Config{
		RegexPresets: map[string]string{
			"custom": `(?P<foo>.+)`,
		},
	}

	t.Run("raw regex takes priority", func(t *testing.T) {
		regex, err := ResolveRegex("custom", `(?P<bar>.*)`, cfg)
		assert.NoError(t, err)
		assert.Equal(t, `(?P<bar>.*)`, regex)
	})

	t.Run("preset works when raw is empty", func(t *testing.T) {
		regex, err := ResolveRegex("custom", "", cfg)
		assert.NoError(t, err)
		assert.Equal(t, `(?P<foo>.+)`, regex)
	})

	t.Run("unknown preset errors", func(t *testing.T) {
		_, err := ResolveRegex("doesnotexist", "", cfg)
		assert.Error(t, err)
	})

	t.Run("default fallback", func(t *testing.T) {
		regex, err := ResolveRegex("", "", cfg)
		assert.NoError(t, err)
		assert.Equal(t, "^(?P<message>.*)$", regex)
	})
}

func TestResolveJqFilter(t *testing.T) {
	cfg := &Config{
		JQPresets: map[string]string{
			"simple": ".message",
		},
	}

	t.Run("raw jq takes priority", func(t *testing.T) {
		jq, err := ResolveJqFilter("simple", ".raw", cfg)
		assert.NoError(t, err)
		assert.Equal(t, ".raw", jq)
	})

	t.Run("preset works when raw is empty", func(t *testing.T) {
		jq, err := ResolveJqFilter("simple", "", cfg)
		assert.NoError(t, err)
		assert.Equal(t, ".message", jq)
	})

	t.Run("unknown preset errors", func(t *testing.T) {
		_, err := ResolveJqFilter("doesnotexist", "", cfg)
		assert.Error(t, err)
	})

	t.Run("default fallback is empty", func(t *testing.T) {
		jq, err := ResolveJqFilter("", "", cfg)
		assert.NoError(t, err)
		assert.Equal(t, "", jq)
	})
}

func TestGetRegexPresets(t *testing.T) {
	cfg := &Config{
		RegexPresets: map[string]string{
			"nginx": `(?P<ip>\d+)`,
		},
	}

	presets := GetRegexPresets(cfg)
	assert.Contains(t, presets, "apache")
	assert.Contains(t, presets, "nginx")
	assert.Equal(t, `(?P<ip>\d+)`, presets["nginx"])
}

func TestGetJqPresets(t *testing.T) {
	cfg := &Config{
		JQPresets: map[string]string{
			"short": ".msg",
		},
	}

	presets := GetJqPresets(cfg)
	assert.Contains(t, presets, "short")
	assert.Equal(t, ".msg", presets["short"])
}

func TestResolveParsingOptions(t *testing.T) {
	cfg := &Config{
		RegexPresets: map[string]string{"foo": `(?P<message>abc)`},
		JQPresets:    map[string]string{"bar": `.log`},
	}

	resolved, err := ResolveParsingOptions("foo", "", "bar", "", cfg)
	assert.NoError(t, err)
	assert.Equal(t, `(?P<message>abc)`, resolved.Regex)
	assert.Equal(t, `.log`, resolved.JQ)
}
