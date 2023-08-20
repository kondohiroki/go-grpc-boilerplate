package timeconvert

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/golang-module/carbon/v2"
	"github.com/kondohiroki/go-grpc-boilerplate/internal/logger"
	"go.uber.org/zap"
)

func ToString(t *time.Time) *string {
	if t == nil {
		return nil
	}

	result := carbon.Parse(t.String()).SetTimezone(carbon.UTC).ToString()
	return &result
}

func FromStringUnixToTime(s string) *time.Time {
	if s == "" {
		return nil
	}

	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		logger.Log.Error("Failed to parse string to int64", zap.Error(err))
		return nil
	}

	t := time.Unix(i, 0)

	return &t
}

func FromPointerString(s *string) *time.Time {
	if s == nil {
		return nil
	}

	stime := *s
	t := carbon.Parse(stime).ToStdTime()

	return &t
}

// Helper function to parse duration string into time.Duration
// Example: 1m, 2h, 3d, or 1m2h3d
func ParseDurationString(duration string) (time.Duration, error) {
	parts := strings.Split(duration, "")

	var totalDuration time.Duration
	var numStr string
	var timeUnit time.Duration
	for _, part := range parts {
		// Check if the part is a numeric value
		if unicode.IsDigit([]rune(part)[0]) {
			numStr += part
			continue
		}

		if numStr == "" {
			return 0, fmt.Errorf("invalid duration: %s", duration)
		}

		// Extract the numeric part of the duration
		num, err := strconv.Atoi(numStr)
		if err != nil {
			return 0, fmt.Errorf("invalid duration: %s", duration)
		}

		// Extract the unit part of the duration
		unit := part

		// Determine the time unit based on the provided string
		switch strings.ToLower(unit) {
		case "m":
			timeUnit = time.Minute
		case "h":
			timeUnit = time.Hour
		case "d":
			timeUnit = time.Hour * 24
		default:
			return 0, errors.New("invalid duration unit")
		}

		totalDuration += time.Duration(num) * timeUnit
		numStr = ""
	}

	// Handle the case when the duration string ends with a numeric value
	if numStr != "" {
		num, err := strconv.Atoi(numStr)
		if err != nil {
			return 0, fmt.Errorf("invalid duration: %s", duration)
		}

		totalDuration += time.Duration(num) * timeUnit
	}

	return totalDuration, nil
}
