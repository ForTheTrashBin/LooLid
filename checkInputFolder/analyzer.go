package checkInputFolder

import (
	"strings"

	"golang.org/x/text/unicode/norm"
)

func analyzeCollisions(entries []Entry) []Issue {

	var issues []Issue

	issues = append(
		issues,
		analyzeCaseCollisions(entries)...,
	)

	issues = append(
		issues,
		analyzeUnicodeCollisions(entries)...,
	)

	return issues
}

func analyzeCaseCollisions(entries []Entry) []Issue {

	var issues []Issue

	seen := make(
		map[string]string,
	)

	for _, entry := range entries {

		key := strings.ToLower(
			entry.Path,
		)

		if old, exists := seen[key]; exists {

			issues = append(
				issues,
				Issue{
					Path: entry.Path,
					Info: "conflicts with " + old,
				},
			)

		} else {

			seen[key] = entry.Path
		}
	}

	return issues
}

func analyzeUnicodeCollisions(entries []Entry) []Issue {

	var issues []Issue
	seen := make(
		map[string]string,
	)

	for _, entry := range entries {

		key := strings.ToLower(
			norm.NFC.String(
				entry.Path,
			),
		)

		if old, exists := seen[key]; exists {

			issues = append(
				issues,
				Issue{
					Path: entry.Path,
					Info: "unicode conflict with " + old,
				},
			)

		} else {

			seen[key] = entry.Path
		}

	}
	return issues
}
