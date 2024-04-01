package services

import (
	"fmt"
	"regexp"
	"strconv"
)

// Extract the numeric value from text
func ExtractNumeric(input string) (int, error) {
	// Define a regular expression to match numeric characters
	re := regexp.MustCompile(`(\d+)`)

	// Find the first match in the input string
	match := re.FindStringSubmatch(input)

	// Check if a match is found
	if len(match) < 2 {
		return 0, fmt.Errorf("no numeric value found in the input string")
	}

	// Convert the matched string to an integer
	value, err := strconv.Atoi(match[1])
	if err != nil {
		return 0, err
	}

	return value, nil
}
