package configParser

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/ForTheTrashBin/LooLid/helper/constants"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

//-----------------------------------------------------------------------------
// Creates a new RulesConfig with default values
//-----------------------------------------------------------------------------

func NewRulesConfig() *RulesConfig {

	rulesConfig := &RulesConfig{}

	rulesConfig.Blacklist.inherit = true
	rulesConfig.Whitelist.inherit = true

	return rulesConfig
}

//-----------------------------------------------------------------------------
// Creates a new RulesConfig as a deep copy of the specified RulesConfig
//-----------------------------------------------------------------------------

func NewRulesConfigFromRulesConfig(rulesConfig *RulesConfig) *RulesConfig {

	if rulesConfig == nil {

		return nil
	}

	return &RulesConfig{

		FrontMatter: CloneFrontMatter(rulesConfig.FrontMatter),
		Blacklist:   CloneRuleSet(&rulesConfig.Blacklist),
		Whitelist:   CloneRuleSet(&rulesConfig.Whitelist),
	}
}

//-----------------------------------------------------------------------------
// Creates a new RulesConfig from the specified file
//-----------------------------------------------------------------------------

func NewRulesConfigFromFile(localizer *i18n.Localizer, fileName string, patternDepth int) (*RulesConfig, error) {

	data, err := os.ReadFile(fileName)

	if err != nil {

		if os.IsNotExist(err) {

			return NewRulesConfig(), nil
		}

		return nil, err
	}

	//-------------------------------------------------------------------------

	rulesConfig := NewRulesConfig()

	frontMatter, remainingContent, err := parseFrontmatterBlock(string(data))

	if err != nil {

		relativePath, relError := filepath.Rel(constants.AppConfig_DefaultInputFolder, fileName)

		if relError != nil {

			relativePath = filepath.Base(fileName)
		}

		if localizer != nil {

			return nil, errors.New(localizer.MustLocalize(&i18n.LocalizeConfig{

				MessageID: constants.BlackWhiteError_ReadConfig,
				TemplateData: map[string]string{

					"value1": relativePath,
					"value2": err.Error(),
				}}))
		}

		return nil, fmt.Errorf("read config %s: %w", relativePath, err)
	}

	if frontMatter != nil {

		rulesConfig.FrontMatter = frontMatter
	}

	if _, err := toml.Decode(remainingContent, rulesConfig); err != nil {

		relativePath, relError := filepath.Rel(constants.AppConfig_DefaultInputFolder, fileName)

		if relError != nil {

			relativePath = filepath.Base(fileName)
		}

		if localizer != nil {

			return nil, errors.New(localizer.MustLocalize(&i18n.LocalizeConfig{

				MessageID: constants.BlackWhiteError_ReadConfig,
				TemplateData: map[string]string{

					"value1": relativePath,
					"value2": err.Error(),
				}}))
		}

		return nil, fmt.Errorf("read config %s: %w", relativePath, err)
	}

	if rulesConfig.FrontMatter == nil {

		rulesConfig.FrontMatter = map[string]any{}
	}

	return rulesConfig, rulesConfig.prepare(localizer, patternDepth)
}

func parseFrontmatterBlock(content string) (map[string]any, string, error) {

	trimmed := strings.TrimSpace(content)

	if trimmed == "" {
		return nil, content, nil
	}

	lines := strings.Split(trimmed, "\n")

	if len(lines) < 2 {
		return nil, content, nil
	}

	openingDelimiter := ""
	closingDelimiter := ""

	if strings.TrimSpace(lines[0]) == "+++" {

		openingDelimiter = "+++"
		closingDelimiter = "+++"
	} else if strings.TrimSpace(lines[0]) == "---" {

		openingDelimiter = "---"
		closingDelimiter = "---"
	} else {

		return nil, content, nil
	}

	closingIndex := -1

	for idx := 1; idx < len(lines); idx++ {

		trimmedLine := strings.TrimSpace(lines[idx])

		if trimmedLine == closingDelimiter {

			closingIndex = idx
			break
		}

		if trimmedLine == "+++" || trimmedLine == "---" {

			return nil, content, errors.New("mixed frontmatter delimiters are not allowed")
		}
	}

	if closingIndex == -1 {

		return nil, content, errors.New("unterminated frontmatter block")
	}

	frontMatterLines := lines[1:closingIndex]
	frontMatterContent := strings.Join(frontMatterLines, "\n")

	frontMatter := map[string]any{}

	if openingDelimiter == "+++" {

		if _, err := toml.Decode(frontMatterContent, &frontMatter); err != nil {

			return nil, content, err
		}
	} else {

		parsed, err := parseSimpleYAMLFrontMatter(frontMatterContent)

		if err != nil {

			return nil, content, err
		}

		frontMatter = parsed
	}

	remainingLines := append([]string{}, lines[closingIndex+1:]...)
	remainingContent := strings.Join(remainingLines, "\n")

	return frontMatter, remainingContent, nil
}

func parseSimpleYAMLFrontMatter(content string) (map[string]any, error) {

	frontMatter := map[string]any{}

	for _, rawLine := range strings.Split(content, "\n") {

		line := strings.TrimSpace(rawLine)

		if line == "" {

			continue
		}

		parts := strings.SplitN(line, ":", 2)

		if len(parts) != 2 {

			return nil, fmt.Errorf("invalid YAML frontmatter line: %q", line)
		}

		key := strings.TrimSpace(parts[0])

		if key == "" {

			return nil, fmt.Errorf("invalid YAML frontmatter line: %q", line)
		}

		value := strings.TrimSpace(parts[1])

		if value == "" {

			frontMatter[key] = ""

			continue
		}

		if (len(value) >= 2 && value[0] == '"' && value[len(value)-1] == '"') ||
			(len(value) >= 2 && value[0] == '\'' && value[len(value)-1] == '\'') {

			frontMatter[key] = value[1 : len(value)-1]

			continue
		}

		if boolValue, err := strconv.ParseBool(value); err == nil {

			frontMatter[key] = boolValue

			continue
		}

		if intValue, err := strconv.Atoi(value); err == nil {

			frontMatter[key] = intValue

			continue
		}

		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {

			frontMatter[key] = floatValue

			continue
		}

		frontMatter[key] = value
	}

	return frontMatter, nil
}

//-----------------------------------------------------------------------------
// CloneRuleSet creates a deep copy of the specified RuleSet.
//-----------------------------------------------------------------------------

func CloneRuleSet(ruleSet *RuleSet) RuleSet {

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

	//-------------------------------------------------------------------------

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
//-----------------------------------------------------------------------------

func CloneFrontMatter(frontMatter map[string]any) map[string]any {

	if frontMatter == nil { // Security prompt

		return map[string]any{} // An empty map
	}

	//-------------------------------------------------------------------------

	copiedFrontMatter := make(map[string]any, len(frontMatter))

	for key, value := range frontMatter {

		if subMap, isSubMap := value.(map[string]any); isSubMap {

			copiedFrontMatter[key] = CloneFrontMatter(subMap)
		} else {

			copiedFrontMatter[key] = value // Primitive types such as strings and integers are safe
		}
	}

	return copiedFrontMatter
}

func MergeFrontMatter(base, addendum map[string]any) map[string]any {

	result := CloneFrontMatter(base)

	for addendumKey, addendumVal := range addendum {

		if baseVal, exists := result[addendumKey]; exists {

			baseMap, baseOk := baseVal.(map[string]any)
			addendumMap, addendumOk := addendumVal.(map[string]any)

			if baseOk && addendumOk {

				result[addendumKey] = MergeFrontMatter(baseMap, addendumMap)
			} else {
				if subMap2, ok := addendumVal.(map[string]any); ok {

					result[addendumKey] = CloneFrontMatter(subMap2)
				} else {

					result[addendumKey] = addendumVal
				}
			}
		} else {

			if subMap2, ok := addendumVal.(map[string]any); ok {

				result[addendumKey] = CloneFrontMatter(subMap2)
			} else {

				result[addendumKey] = addendumVal
			}
		}
	}

	return result
}

//-----------------------------------------------------------------------------
// Checks all rules of Black- and Whitelist
//-----------------------------------------------------------------------------

func (rulesConfig *RulesConfig) IsListed(name string, isDir bool) bool {

	bestDepth := -1

	blackMatch := false
	whiteMatch := false

	for idx := range rulesConfig.Blacklist.Rules {

		rule := &rulesConfig.Blacklist.Rules[idx]

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

	for idx := range rulesConfig.Whitelist.Rules {

		rule := &rulesConfig.Whitelist.Rules[idx]

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

func (rulesConfig *RulesConfig) IsBlacklisted(name string, isDir bool) bool {

	for idx := range rulesConfig.Blacklist.Rules {

		if rulesConfig.Blacklist.Rules[idx].matchRule(name, isDir) {

			return true
		}
	}

	return false
}

func (rulesConfig *RulesConfig) DoBlacklistInherit() bool {

	return rulesConfig.Blacklist.inherit
}

//-----------------------------------------------------------------------------
// IsWhitelisted checks the Whitelist only
//-----------------------------------------------------------------------------

func (rulesConfig *RulesConfig) IsWhitelisted(name string, isDir bool) bool {

	for idx := range rulesConfig.Whitelist.Rules {

		if rulesConfig.Whitelist.Rules[idx].matchRule(name, isDir) {

			return true
		}
	}

	return false
}

func (rulesConfig *RulesConfig) DoWhitelistInherit() bool {

	return rulesConfig.Whitelist.inherit
}

//-----------------------------------------------------------------------------

func (ruleSet *RuleSet) MergeRuleSet(src *RuleSet) {

	ruleSet.TomlInherit = src.TomlInherit

	ruleSet.Rules = append(ruleSet.Rules, src.Rules...)

	ruleSet.inherit = src.inherit
}

//-----------------------------------------------------------------------------
