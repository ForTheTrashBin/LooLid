package blackwhite

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

//-----------------------------------------------------------------------------
// loadLocalFile reads teh specified 'LooLid.comfig'
//-----------------------------------------------------------------------------

func LoadConfigFile(fileName string) (*Config, bool, error) {

	var doDebug bool = false

	if doDebug {

		fmt.Println("*********************************************************************************")
		fmt.Print("*** LoadConfigFile fileName <", fileName, ">\n")
		fmt.Println("*********************************************************************************")
	}

	data, err := os.ReadFile(fileName)

	if err != nil {

		if os.IsNotExist(err) {

			return &Config{}, false, nil
		}

		return nil, false, err
	}

	if doDebug {

		fmt.Println("*** ConfigFile read")
	}

	var cfg Config

	if _, err := toml.Decode(string(data), &cfg); err != nil {

		if doDebug {

			out := fmt.Errorf("%s: %w", fileName, err)

			fmt.Println(out)
		}

		return nil, false, fmt.Errorf("%s: %w", fileName, err)
	}

	if doDebug {

		fmt.Println("*** ConfigFile decoded")
	}

	return &cfg, true, cfg.prepareConfig()
}

//-----------------------------------------------------------------------------
// Checks all rules of Black- and Whitelist
//-----------------------------------------------------------------------------

func (cfg *Config) IsListed(name string, isDir bool) bool {

	isListed := cfg.IsBlacklisted(name, isDir)

	if isListed {

		if cfg.IsWhitelisted(name, isDir) {

			isListed = false
		}
	}

	return isListed
}

//-----------------------------------------------------------------------------
// IsBlacklisted checks the Blacklist only
//-----------------------------------------------------------------------------

func (cfg *Config) IsBlacklisted(name string, isDir bool) bool {

	for idx := range cfg.Blacklist.Rules {

		if cfg.Blacklist.Rules[idx].matchRule(name, isDir) {

			return true
		}
	}

	return false
}

//-----------------------------------------------------------------------------
// IsWhitelisted checks the Whitelist only
//-----------------------------------------------------------------------------

func (cfg *Config) IsWhitelisted(name string, isDir bool) bool {

	for idx := range cfg.Whitelist.Rules {

		if cfg.Whitelist.Rules[idx].matchRule(name, isDir) {

			return true
		}
	}

	return false
}

//-----------------------------------------------------------------------------

func (rs *RuleSet) MergeRuleSet(src *RuleSet) {

	rs.Rules = append(rs.Rules, src.Rules...)
}

//-----------------------------------------------------------------------------
