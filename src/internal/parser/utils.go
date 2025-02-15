package parser

import (
	"regexp"
	"strings"
)

func sanitazeExpression(traceId, expression string) (string, error) {
	if expression[len(expression)-1:] != ";" {
		return "", errorWrongExpressionFormat(traceId, expression)
	}
	return expression[:len(expression)-1], nil
}

func isMatch(pattern regex_pattern, expression string) bool {
	r := regexp.MustCompile(string(pattern))
	return r.MatchString(expression)
}

func findMatch(pattern regex_pattern, expression string) (bool, string) {
	r := regexp.MustCompile(string(pattern))
	str := r.FindString(expression)
	if str == "" {
		return false, str
	}
	return true, str
}

func extractKeyFromParanthesis(expression string) (int, []string) {
	r := regexp.MustCompile(string(valueInParanthesis))
	matches := r.FindAllString(expression, -1)

	if matches == nil {
		return 0, nil
	}

	for i := range matches {
		matches[i] = strings.Trim(matches[i], "'")
	}

	return len(matches), matches
}

func extractKeysFromBrackets(expression string) (int, [][]string) {
	r := regexp.MustCompile(string(valueInBrackets))
	matches := r.FindAllString(expression, -1)

	if matches == nil {
		return 0, nil
	}
	result := [][]string{}

	for i := range matches {
		content := strings.Trim(matches[i], "()")
		items := strings.Split(content, ",")

		partResult := []string{}
		for _, item := range items {
			partResult = append(partResult, strings.Trim(strings.TrimSpace(item), "'"))
		}
		result = append(result, partResult)
	}

	return len(matches), result
}
