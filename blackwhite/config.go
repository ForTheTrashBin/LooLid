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
// loadLocalFile reads teh specified 'LooLid.comfig'
//-----------------------------------------------------------------------------

func LoadConfigFile(localizer *i18n.Localizer, fileName string) (*Config, bool, error) {

	data, err := os.ReadFile(fileName)

	if err != nil {

		if os.IsNotExist(err) {

			return &Config{}, false, nil
		}

		return nil, false, err
	}

	//-------------------------------------------------------------------------

	var cfg Config

	if _, err := toml.Decode(string(data), &cfg); err != nil {

		relativePath, relError := filepath.Rel(constants.AppConfig_DefaultInputFolder, fileName)

		if relError != nil {

			relativePath = filepath.Base(fileName)
		}

		return nil, false, errors.New(localizer.MustLocalize(&i18n.LocalizeConfig{

			MessageID: constants.BlackWhiteError_ReadConfig,
			TemplateData: map[string]string{

				"value1": relativePath,
				"value2": err.Error(),
			}}))
	}

	return &cfg, true, cfg.prepare(localizer)
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

func (cfg *Config) DoBlacklistInherit() bool {

	return cfg.Blacklist.inherit
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

func (cfg *Config) DoWhitelistInherit() bool {

	return cfg.Whitelist.inherit
}

//-----------------------------------------------------------------------------

func (rs *RuleSet) MergeRuleSet(src *RuleSet) {

	rs.Rules = append(rs.Rules, src.Rules...)
}

//-----------------------------------------------------------------------------
