package blackwhite

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/ForTheTrashBin/LooLid/helper/constants"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

//-----------------------------------------------------------------------------
// Creates a new blackwhite.BWConfig with default values
//-----------------------------------------------------------------------------

func NewBWConfig() *BWConfig {

	bwConfig := &BWConfig{}

	bwConfig.Blacklist.inherit = true
	bwConfig.Whitelist.inherit = true

	return bwConfig
}

//-----------------------------------------------------------------------------
// NewRuleSetFromRuleSet creates a deep copy of the specified RuleSet.
//-----------------------------------------------------------------------------

func NewRuleSetFromRuleSet(ruleSet *RuleSet) RuleSet {

	copyBoolPtr := func(src *bool) *bool {

		if src == nil {

			return nil
		}

		dst := *src

		return &dst
	}

	copyRules := func(rules []Rule) []Rule {

		if rules == nil {

			return nil
		}

		dst := make([]Rule, len(rules))

		for i := range rules {

			dst[i] = Rule{

				TomlPattern:    rules[i].TomlPattern,
				TomlScope:      rules[i].TomlScope,
				patternDepth:   rules[i].patternDepth,
				patternMatcher: rules[i].patternMatcher,
				patternScope:   rules[i].patternScope,
			}
		}

		return dst
	}

	if ruleSet == nil {

		return RuleSet{}
	}

	return RuleSet{
		TomlInherit: copyBoolPtr(ruleSet.TomlInherit),
		Rules:       copyRules(ruleSet.Rules),
		inherit:     ruleSet.inherit,
	}
}

//-----------------------------------------------------------------------------
// Creates a new BWConfig as a deep copy of the specified BWConfig
//-----------------------------------------------------------------------------

func NewBWConfigFromBWConfig(bwConfig *BWConfig) *BWConfig {

	if bwConfig == nil {

		return nil
	}

	return &BWConfig{

		Blacklist: NewRuleSetFromRuleSet(&bwConfig.Blacklist),
		Whitelist: NewRuleSetFromRuleSet(&bwConfig.Whitelist),
	}
}

//-----------------------------------------------------------------------------
// Creates a new blackwhite.BWConfig from the specified file
//-----------------------------------------------------------------------------

func NewBWConfigFromFile(localizer *i18n.Localizer, fileName string, patternDepth int) (*BWConfig, error) {

	data, err := os.ReadFile(fileName)

	if err != nil {

		if os.IsNotExist(err) {

			return NewBWConfig(), nil
		}

		return nil, err
	}

	//-------------------------------------------------------------------------

	bwConfig := NewBWConfig()

	if _, err := toml.Decode(string(data), bwConfig); err != nil {

		relativePath, relError := filepath.Rel(constants.AppConfig_DefaultInputFolder, fileName)

		if relError != nil {

			relativePath = filepath.Base(fileName)
		}

		return nil, errors.New(localizer.MustLocalize(&i18n.LocalizeConfig{

			MessageID: constants.BlackWhiteError_ReadConfig,
			TemplateData: map[string]string{

				"value1": relativePath,
				"value2": err.Error(),
			}}))
	}

	return bwConfig, bwConfig.prepare(localizer, patternDepth)
}

//-----------------------------------------------------------------------------
// Checks all rules of Black- and Whitelist
//-----------------------------------------------------------------------------

func (bwConfig *BWConfig) IsListed(name string, isDir bool) bool {

	bestDepth := -1

	blackMatch := false
	whiteMatch := false

	for idx := range bwConfig.Blacklist.Rules {

		rule := &bwConfig.Blacklist.Rules[idx]

		if !rule.matchRule(name, isDir) {

			continue
		}

		if rule.patternDepth > bestDepth {

			bestDepth = rule.patternDepth
			blackMatch = true
			whiteMatch = false
		} else if rule.patternDepth == bestDepth {

			blackMatch = true
		}
	}

	for idx := range bwConfig.Whitelist.Rules {

		rule := &bwConfig.Whitelist.Rules[idx]

		if !rule.matchRule(name, isDir) {

			continue
		}

		if rule.patternDepth > bestDepth {

			bestDepth = rule.patternDepth
			blackMatch = false
			whiteMatch = true
		} else if rule.patternDepth == bestDepth {

			whiteMatch = true
		}
	}

	if whiteMatch {

		return false
	}

	return blackMatch
}

//-----------------------------------------------------------------------------
// IsBlacklisted checks the Blacklist only
//-----------------------------------------------------------------------------

func (bwConfig *BWConfig) IsBlacklisted(name string, isDir bool) bool {

	for idx := range bwConfig.Blacklist.Rules {

		if bwConfig.Blacklist.Rules[idx].matchRule(name, isDir) {

			return true
		}
	}

	return false
}

func (bwConfig *BWConfig) DoBlacklistInherit() bool {

	return bwConfig.Blacklist.inherit
}

//-----------------------------------------------------------------------------
// IsWhitelisted checks the Whitelist only
//-----------------------------------------------------------------------------

func (bwConfig *BWConfig) IsWhitelisted(name string, isDir bool) bool {

	for idx := range bwConfig.Whitelist.Rules {

		if bwConfig.Whitelist.Rules[idx].matchRule(name, isDir) {

			return true
		}
	}

	return false
}

func (bwConfig *BWConfig) DoWhitelistInherit() bool {

	return bwConfig.Whitelist.inherit
}

//-----------------------------------------------------------------------------

func (ruleSet *RuleSet) MergeRuleSet(src *RuleSet) {

	ruleSet.TomlInherit = src.TomlInherit

	ruleSet.Rules = append(ruleSet.Rules, src.Rules...)

	ruleSet.inherit = src.inherit
}

//-----------------------------------------------------------------------------
