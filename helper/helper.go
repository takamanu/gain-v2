package helper

import "regexp"

func RemoveNumberAndPeriods(input string) string {
	// Define a regex to remove leading numbers and periods (e.g., "5. " or "10. ")
	re := regexp.MustCompile(`^\d+\.\s*`)
	return re.ReplaceAllString(input, "")
}
