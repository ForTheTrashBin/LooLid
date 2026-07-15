package blackwhite

import (
	"errors"
	"strings"

	"github.com/ForTheTrashBin/LooLid/helper/constants"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

//-----------------------------------------------------------------------------
// 'Matcher' is the internal abstraction for all pattern types.
//-----------------------------------------------------------------------------

type patternMatcher interface {
	matchPattern(name string) bool
}

//-----------------------------------------------------------------------------
// Scope specifies whether a rule applies to files, directories or both.
//-----------------------------------------------------------------------------

type patternScope uint8

const (
	patternScopeBoth patternScope = iota
	patternScopeFile
	patternScopeDir
)

//-----------------------------------------------------------------------------

const (
	patternPrefixGlob   = "glob:" // This is the default
	patternPrefixRegex  = "regex:"
	patternPreficExact  = "exact:"
	patternPrefixPrefix = "prefix:"
	patternPrefixSuffix = "suffix:"
)

//-----------------------------------------------------------------------------
// prepare all rules
//-----------------------------------------------------------------------------

func (cfg *Config) prepare(localizer *i18n.Localizer) error {

	if err := cfg.Blacklist.prepare(localizer); err != nil {

		return errors.New(localizer.MustLocalize(&i18n.LocalizeConfig{

			MessageID: constants.BlackWhiteError_Blacklist,
			TemplateData: map[string]string{

				"value1": err.Error(),
			}}))
	}

	if err := cfg.Whitelist.prepare(localizer); err != nil {

		return errors.New(localizer.MustLocalize(&i18n.LocalizeConfig{

			MessageID: constants.BlackWhiteError_Whitelist,
			TemplateData: map[string]string{

				"value1": err.Error(),
			}}))
	}

	return nil
}

//-----------------------------------------------------------------------------

func (rs *RuleSet) prepare(localizer *i18n.Localizer) error {

	if rs.TomlInherit == nil {

		rs.inherit = true
	} else {

		rs.inherit = *rs.TomlInherit
	}

	//-------------------------------------------------------------------------

	for idx := range rs.Rules {

		rule := &rs.Rules[idx]

		//---------------------------------------------------------------------

		patternMatcher, err := createPatternMatcher(localizer, rule.TomlPattern)

		if err != nil {

			return err
		}

		rule.patternMatcher = patternMatcher

		//---------------------------------------------------------------------

		patternScope, err := parseTomlScope(localizer, rule.TomlScope)

		if err != nil {

			return err
		}

		rule.patternScope = patternScope
	}

	return nil
}

//-----------------------------------------------------------------------------

func parseTomlScope(localizer *i18n.Localizer, tomlScope string) (patternScope, error) {

	switch strings.ToLower(tomlScope) {

	case "", "both":

		return patternScopeBoth, nil

	case "file":

		return patternScopeFile, nil

	case "dir":

		return patternScopeDir, nil

	default:

		return patternScopeBoth, errors.New(localizer.MustLocalize(&i18n.LocalizeConfig{

			MessageID: constants.BlackWhiteError_InvalidScope,
			TemplateData: map[string]string{

				"value1": tomlScope,
			}}))
	}
}

//-----------------------------------------------------------------------------

func (ar *Rule) matchRule(name string, isDir bool) bool {

	switch ar.patternScope {

	case patternScopeFile:

		if isDir {

			return false
		}

	case patternScopeDir:

		if !isDir {

			return false
		}
	}

	return ar.patternMatcher.matchPattern(name)
}

//-----------------------------------------------------------------------------
