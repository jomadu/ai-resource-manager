package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// parseDuration parses a duration string with units (m, h, d)
// Examples: "30m", "2h", "7d", "1h30m"
func parseDuration(s string) (time.Duration, error) {
	if s == "" {
		return 0, fmt.Errorf("empty duration string")
	}

	// Remove any whitespace
	s = strings.TrimSpace(s)

	// Use regex to parse the duration string
	// Supports formats like: "30m", "2h", "7d", "1h30m", "2d12h"
	re := regexp.MustCompile(`^(\d+)([mhd])$|^(\d+)([mhd])(\d+)([mhd])$`)
	matches := re.FindStringSubmatch(s)

	if len(matches) == 0 {
		return 0, fmt.Errorf("invalid duration format: %s (expected format like '30m', '2h', '7d')", s)
	}

	var totalMinutes int64

	switch {
	case matches[1] != "" && matches[2] != "":
		// Single unit format: "30m", "2h", "7d"
		value, err := strconv.ParseInt(matches[1], 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid number in duration: %s", matches[1])
		}
		unit := matches[2]

		switch unit {
		case "m":
			totalMinutes = value
		case "h":
			totalMinutes = value * 60
		case "d":
			totalMinutes = value * 60 * 24
		default:
			return 0, fmt.Errorf("unsupported unit: %s (supported: m, h, d)", unit)
		}
	case matches[3] != "" && matches[4] != "" && matches[5] != "" && matches[6] != "":
		// Two unit format: "1h30m", "2d12h"
		value1, err := strconv.ParseInt(matches[3], 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid number in duration: %s", matches[3])
		}
		unit1 := matches[4]
		value2, err := strconv.ParseInt(matches[5], 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid number in duration: %s", matches[5])
		}
		unit2 := matches[6]

		// Convert both units to minutes and add
		var minutes1, minutes2 int64

		switch unit1 {
		case "m":
			minutes1 = value1
		case "h":
			minutes1 = value1 * 60
		case "d":
			minutes1 = value1 * 60 * 24
		default:
			return 0, fmt.Errorf("unsupported unit: %s (supported: m, h, d)", unit1)
		}

		switch unit2 {
		case "m":
			minutes2 = value2
		case "h":
			minutes2 = value2 * 60
		case "d":
			minutes2 = value2 * 60 * 24
		default:
			return 0, fmt.Errorf("unsupported unit: %s (supported: m, h, d)", unit2)
		}

		totalMinutes = minutes1 + minutes2
	default:
		return 0, fmt.Errorf("invalid duration format: %s", s)
	}

	return time.Duration(totalMinutes) * time.Minute, nil
}
