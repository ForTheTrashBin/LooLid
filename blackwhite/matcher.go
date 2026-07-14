package blackwhite

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

//-----------------------------------------------------------------------------
// 'createPatternMatcher' generates a fully-formed matcher from a pattern.
//
// Examples:
//
//	*.tmp
//	glob:*.tmp
//	regex:^cache_[0-9]+$
//	exact:README.md
//	prefix:cache_
//	suffix:.bak
//
//-----------------------------------------------------------------------------

const (
	patternPrefixRegex  = "regex:"
	patternPrefixGlob   = "glob:"
	patternPreficExact  = "exact:"
	patternPrefixPrefix = "prefix:"
	patternPrefixSuffix = "suffix:"
)

//-----------------------------------------------------------------------------

func createPatternMatcher(pattern string) (PatternMatcher, error) {

	if pattern == "" {

		return nil, fmt.Errorf("empty pattern")
	}

	switch {

	case strings.HasPrefix(pattern, patternPrefixRegex):

		return newPatternMatcherRegex(strings.TrimPrefix(pattern, patternPrefixRegex))

	case strings.HasPrefix(pattern, patternPrefixGlob):

		return newPatternMatcherGlob(strings.TrimPrefix(pattern, patternPrefixGlob))

	case strings.HasPrefix(pattern, patternPreficExact):

		return newPatternMatcherExact(strings.TrimPrefix(pattern, patternPreficExact))

	case strings.HasPrefix(pattern, patternPrefixPrefix):

		return newPatternMatcherPrefix(strings.TrimPrefix(pattern, patternPrefixPrefix))

	case strings.HasPrefix(pattern, patternPrefixSuffix):

		return newPatternMatcherSuffix(strings.TrimPrefix(pattern, patternPrefixSuffix))

	default:

		return newPatternMatcherGlob(pattern) // Standard: Glob
	}
}

//-----------------------------------------------------------------------------
// Regex PatternMatcher
//-----------------------------------------------------------------------------

type patternMatcherRegex struct {
	re *regexp.Regexp
}

func newPatternMatcherRegex(expr string) (PatternMatcher, error) {

	re, err := regexp.Compile(expr)

	if err != nil {

		return nil, fmt.Errorf("invalid regex %q: %w", expr, err)
	}

	return &patternMatcherRegex{re: re}, nil
}

func (matcher *patternMatcherRegex) MatchPattern(name string) bool {

	return matcher.re.MatchString(name)
}

//-----------------------------------------------------------------------------
// Glob PatternMatcher
//-----------------------------------------------------------------------------

type patternMatcherGlob struct {
	value string
}

func newPatternMatcherGlob(value string) (PatternMatcher, error) {

	if value == "" {

		return nil, fmt.Errorf("empty glob pattern")
	}

	return &patternMatcherGlob{value: value}, nil
}

func (matcher *patternMatcherGlob) MatchPattern(name string) bool {

	ok, err := filepath.Match(matcher.value, name)

	if err != nil { // This should actually be detected when loading the rules.

		return false
	}

	return ok
}

//-----------------------------------------------------------------------------
// Exact PatternMatcher
//-----------------------------------------------------------------------------

type patternMatcherExact struct {
	value string
}

func newPatternMatcherExact(value string) (PatternMatcher, error) {

	if value == "" {

		return nil, fmt.Errorf("empty exact value")
	}

	return &patternMatcherExact{value: value}, nil
}

func (matcher *patternMatcherExact) MatchPattern(name string) bool {

	return name == matcher.value
}

//-----------------------------------------------------------------------------
// Prefix PatternMatcher
//-----------------------------------------------------------------------------

type patternMatcherPrefix struct {
	value string
}

func newPatternMatcherPrefix(value string) (PatternMatcher, error) {

	if value == "" {

		return nil, fmt.Errorf("empty prefix")
	}

	return &patternMatcherPrefix{value: value}, nil
}

func (matcher *patternMatcherPrefix) MatchPattern(name string) bool {

	return strings.HasPrefix(name, matcher.value)
}

//-----------------------------------------------------------------------------
// Suffix PatternMatcher
//-----------------------------------------------------------------------------

type patternMatcherSuffix struct {
	value string
}

func newPatternMatcherSuffix(value string) (PatternMatcher, error) {

	if value == "" {

		return nil, fmt.Errorf("empty suffix")
	}

	return &patternMatcherSuffix{value: value}, nil
}

func (matcher *patternMatcherSuffix) MatchPattern(name string) bool {

	return strings.HasSuffix(name, matcher.value)
}

//-----------------------------------------------------------------------------
