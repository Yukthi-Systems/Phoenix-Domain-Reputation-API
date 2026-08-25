/*
Copyright (C) 2026 Yukthi Systems Private Limited

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License version 3
as published by the Free Software Foundation.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
version 3 along with this program. If not, see
<https://www.gnu.org/licenses/>.
*/

// Package config loads application configuration from environment
// variables (optionally seeded from a .env file via godotenv).
package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Yukthi-Systems/Phoenix-Domain-Reputation-API/internal/ipfire"
	"github.com/joho/godotenv"
)

// CategoryConfig ties one IPFire category to its list URL and the score
// applied when a domain is found on that list.
type CategoryConfig struct {
	Name  string
	URL   string
	Score int
}

// Config holds all application settings.
type Config struct {
	ServerPort     string
	UpdateInterval time.Duration
	HTTPTimeout    time.Duration
	LogLevel       string
	Categories     []CategoryConfig
}

// defaultScores are used when a category's IPFIRE_<NAME>_SCORE variable is
// not set.
var defaultScores = map[string]int{
	"gambling":    1,
	"malware":     1,
	"phishing":    1,
	"pornography": 1,
	"violence":    1,
}

// Load reads configuration from the environment. If a .env file is present
// in the working directory it is loaded first (missing .env is not an
// error). Load fails fast on malformed values so misconfiguration is caught
// at startup rather than silently producing wrong scores.
func Load() (*Config, error) {
	_ = godotenv.Load()

	updateInterval, err := parseDuration("IPFIRE_UPDATE_INTERVAL", time.Hour)
	if err != nil {
		return nil, err
	}

	httpTimeout, err := parseDuration("IPFIRE_HTTP_TIMEOUT", 30*time.Second)
	if err != nil {
		return nil, err
	}

	categories := make([]CategoryConfig, 0, len(ipfire.Categories))
	for _, c := range ipfire.Categories {
		envVar := "IPFIRE_" + strings.ToUpper(c.Name) + "_SCORE"
		score, err := parseScore(envVar, defaultScores[c.Name])
		if err != nil {
			return nil, err
		}
		categories = append(categories, CategoryConfig{Name: c.Name, URL: c.URL, Score: score})
	}

	return &Config{
		ServerPort:     getEnv("SERVER_PORT", "8080"),
		UpdateInterval: updateInterval,
		HTTPTimeout:    httpTimeout,
		LogLevel:       getEnv("LOG_LEVEL", "INFO"),
		Categories:     categories,
	}, nil
}

// getEnv returns the value of the environment variable key, or fallback if
// it is unset or empty.
func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

// parseDuration reads and parses the environment variable key as a
// time.Duration, returning fallback if it is unset or empty. It returns an
// error if the variable is set to a value time.ParseDuration rejects.
func parseDuration(key string, fallback time.Duration) (time.Duration, error) {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		return fallback, nil
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return 0, fmt.Errorf("invalid duration for %s: %w", key, err)
	}
	return d, nil
}

// parseScore reads and parses the environment variable key as an integer
// category score, returning fallback if it is unset or empty. It returns
// an error if the variable is set to a non-integer value.
func parseScore(key string, fallback int) (int, error) {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		return fallback, nil
	}
	score, err := strconv.Atoi(v)
	if err != nil {
		return 0, fmt.Errorf("invalid score for %s: %w", key, err)
	}
	return score, nil
}
