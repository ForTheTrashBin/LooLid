package blackwhite

//-----------------------------------------------------------------------------
// 'Matcher' is the internal abstraction for all pattern types.
//-----------------------------------------------------------------------------

type PatternMatcher interface {
	MatchPattern(name string) bool
}

//-----------------------------------------------------------------------------
// Scope specifies whether a rule applies to files, directories or both.
//-----------------------------------------------------------------------------

type PatternScope uint8

const (
	PatternScopeBoth PatternScope = iota
	PatternScopeFile
	PatternScopeDir
)

//-----------------------------------------------------------------------------
// Rule describes a single blacklist/whitelist rule.
//-----------------------------------------------------------------------------

type Rule struct {
	TomlPattern string `toml:"pattern"`
	TomlScope   string `toml:"scope"`

	// generated internal values

	patternMatcher PatternMatcher
	patternScope   PatternScope
}

//-----------------------------------------------------------------------------
// RuleSet is a Black- or Whitelist.
//-----------------------------------------------------------------------------

type RuleSet struct {
	TomlInherit *bool  `toml:"inherit"`
	Rules       []Rule `toml:"rule"`

	// generated internal values

	Inherit bool
}

//-----------------------------------------------------------------------------
// Config represents the Back- and Whitelist configuration
//-----------------------------------------------------------------------------

type Config struct {
	Blacklist RuleSet `toml:"blacklist"`
	Whitelist RuleSet `toml:"whitelist"`

	baseDir string
}

//-----------------------------------------------------------------------------
