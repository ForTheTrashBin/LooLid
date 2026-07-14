package blackwhite

import (
	"fmt"
	"strings"
)

//-----------------------------------------------------------------------------
// prepare all rules
//-----------------------------------------------------------------------------

func (cfg *Config) prepareConfig() error {

	if err := prepareRuleSet(&cfg.Blacklist); err != nil {

		return fmt.Errorf("blacklist: %w", err)
	}

	if err := prepareRuleSet(&cfg.Whitelist); err != nil {

		return fmt.Errorf("whitelist: %w", err)
	}

	return nil
}

//-----------------------------------------------------------------------------

func prepareRuleSet(set *RuleSet) error {

	if set.TomlInherit == nil {

		set.Inherit = true
	} else {

		set.Inherit = *set.TomlInherit
	}

	//-------------------------------------------------------------------------

	for i := range set.Rules {

		rule := &set.Rules[i]

		//---------------------------------------------------------------------

		patternMatcher, err := createPatternMatcher(rule.TomlPattern)

		if err != nil {

			return fmt.Errorf("rule %d (%s): %w", i+1, rule.TomlPattern, err)
		}

		rule.patternMatcher = patternMatcher

		//---------------------------------------------------------------------

		patternScope, err := parseTomlScope(rule.TomlScope)

		if err != nil {

			return fmt.Errorf("rule %d: %w", i+1, err)
		}

		rule.patternScope = patternScope

	}

	return nil
}

//-----------------------------------------------------------------------------

func parseTomlScope(tomlScope string) (PatternScope, error) {

	switch strings.ToLower(tomlScope) {

	case "", "both":

		return PatternScopeBoth, nil

	case "file":

		return PatternScopeFile, nil

	case "dir":

		return PatternScopeDir, nil

	default:

		return PatternScopeBoth, fmt.Errorf("unknown scope %q", tomlScope)
	}
}

//-----------------------------------------------------------------------------

func (ar *Rule) matchRule(name string, isDir bool) bool {

	switch ar.patternScope {

	case PatternScopeFile:

		if isDir {

			return false
		}

	case PatternScopeDir:

		if !isDir {

			return false
		}
	}

	return ar.patternMatcher.MatchPattern(name)
}

//-----------------------------------------------------------------------------
