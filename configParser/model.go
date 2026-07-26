package configParser

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
// FrontMatter stores optional metadata for a configuration file.
//-----------------------------------------------------------------------------

type FrontMatter map[string]any

//-----------------------------------------------------------------------------
// Config represents the Black- and Whitelist configuration together with optional frontmatter.
//-----------------------------------------------------------------------------

type RulesConfig struct {
	FrontMatter FrontMatter `toml:"frontmatter"`
	Blacklist   RuleSet     `toml:"blacklist"`
	Whitelist   RuleSet     `toml:"whitelist"`
}

//-----------------------------------------------------------------------------
