package configParser

import (
	"errors"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ForTheTrashBin/LooLid/helper/constants"
	"github.com/nicksnyder/go-i18n/v2/i18n"
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

func createPatternMatcher(localizer *i18n.Localizer, pattern string) (patternMatcher, error) {

	if pattern == "" {

		return nil, errors.New(localizer.MustLocalize(&i18n.LocalizeConfig{

			MessageID: constants.BlackWhiteError_EmptyPattern,
		}))
	}

	switch {

	case strings.HasPrefix(pattern, patternPrefixGlob):

		return newPatternMatcherGlob(localizer, strings.TrimPrefix(pattern, patternPrefixGlob))

	case strings.HasPrefix(pattern, patternPrefixRegex):

		return newPatternMatcherRegex(localizer, strings.TrimPrefix(pattern, patternPrefixRegex))

	case strings.HasPrefix(pattern, patternPreficExact):

		return newPatternMatcherExact(localizer, strings.TrimPrefix(pattern, patternPreficExact))

	case strings.HasPrefix(pattern, patternPrefixPrefix):

		return newPatternMatcherPrefix(localizer, strings.TrimPrefix(pattern, patternPrefixPrefix))

	case strings.HasPrefix(pattern, patternPrefixSuffix):

		return newPatternMatcherSuffix(localizer, strings.TrimPrefix(pattern, patternPrefixSuffix))

	default:

		return newPatternMatcherGlob(localizer, pattern) // Standard: Glob
	}
}

//-----------------------------------------------------------------------------
// Glob PatternMatcher
//-----------------------------------------------------------------------------

type patternMatcherGlob struct {
	value string
}

func newPatternMatcherGlob(localizer *i18n.Localizer, value string) (patternMatcher, error) {

	if value == "" {

		return nil, errors.New(localizer.MustLocalize(&i18n.LocalizeConfig{

			MessageID: constants.BlackWhiteError_EmptyGlobPattern,
		}))
	}

	return &patternMatcherGlob{value: value}, nil
}

func (matcher *patternMatcherGlob) matchPattern(name string) bool {

	ok, err := filepath.Match(matcher.value, name)

	if err != nil { // This should actually be detected when loading the rules.

		return false
	}

	return ok
}

//-----------------------------------------------------------------------------
// Regex PatternMatcher
//-----------------------------------------------------------------------------

type patternMatcherRegex struct {
	re *regexp.Regexp
}

func newPatternMatcherRegex(localizer *i18n.Localizer, expr string) (patternMatcher, error) {

	re, err := regexp.Compile(expr)

	if err != nil {

		return nil, errors.New(localizer.MustLocalize(&i18n.LocalizeConfig{

			MessageID: constants.BlackWhiteError_InvalidRegex,
			TemplateData: map[string]string{

				"value1": expr,
				"value2": err.Error(),
			}}))
	}

	return &patternMatcherRegex{re: re}, nil
}

func (matcher *patternMatcherRegex) matchPattern(name string) bool {

	return matcher.re.MatchString(name)
}

//-----------------------------------------------------------------------------
// Exact PatternMatcher
//-----------------------------------------------------------------------------

type patternMatcherExact struct {
	value string
}

func newPatternMatcherExact(localizer *i18n.Localizer, value string) (patternMatcher, error) {

	if value == "" {

		return nil, errors.New(localizer.MustLocalize(&i18n.LocalizeConfig{

			MessageID: constants.BlackWhiteError_EmptyExactValue,
		}))
	}

	return &patternMatcherExact{value: value}, nil
}

func (matcher *patternMatcherExact) matchPattern(name string) bool {

	return name == matcher.value
}

//-----------------------------------------------------------------------------
// Prefix PatternMatcher
//-----------------------------------------------------------------------------

type patternMatcherPrefix struct {
	value string
}

func newPatternMatcherPrefix(localizer *i18n.Localizer, value string) (patternMatcher, error) {

	if value == "" {

		return nil, errors.New(localizer.MustLocalize(&i18n.LocalizeConfig{

			MessageID: constants.BlackWhiteError_EmptyPrefix,
		}))
	}

	return &patternMatcherPrefix{value: value}, nil
}

func (matcher *patternMatcherPrefix) matchPattern(name string) bool {

	return strings.HasPrefix(name, matcher.value)
}

//-----------------------------------------------------------------------------
// Suffix PatternMatcher
//-----------------------------------------------------------------------------

type patternMatcherSuffix struct {
	value string
}

func newPatternMatcherSuffix(localizer *i18n.Localizer, value string) (patternMatcher, error) {

	if value == "" {

		return nil, errors.New(localizer.MustLocalize(&i18n.LocalizeConfig{

			MessageID: constants.BlackWhiteError_EmptySuffix,
		}))
	}

	return &patternMatcherSuffix{value: value}, nil
}

func (matcher *patternMatcherSuffix) matchPattern(name string) bool {

	return strings.HasSuffix(name, matcher.value)
}

//-----------------------------------------------------------------------------
