package blackwhite

//-----------------------------------------------------------------------------
// Rule describes a single blacklist/whitelist rule.
//-----------------------------------------------------------------------------

type Rule struct {
	TomlPattern string `toml:"pattern"`
	TomlScope   string `toml:"scope"`

	// generated internal values

	patternMatcher patternMatcher
	patternScope   patternScope
	patternDepth   int
}

//-----------------------------------------------------------------------------
// RuleSet is a Black- or Whitelist.
//-----------------------------------------------------------------------------

type RuleSet struct {
	TomlInherit *bool  `toml:"inherit"`
	Rules       []Rule `toml:"rule"`

	// generated internal values

	inherit bool
}

//-----------------------------------------------------------------------------
// Config represents the Black- and Whitelist configuration
//-----------------------------------------------------------------------------

type BWConfig struct {
	Blacklist RuleSet `toml:"blacklist"`
	Whitelist RuleSet `toml:"whitelist"`
}

//-----------------------------------------------------------------------------
