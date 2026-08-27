package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/destyk/commitkit/internal/lint"
	"github.com/destyk/commitkit/internal/policy"
	"gopkg.in/yaml.v3"
)

// FileName is the default configuration file name.
const FileName = ".commitkit.yml"

// Config is the on-disk representation of .commitkit.yml.
//
// Omitted fields fall back to Conventional Commits defaults when
// building the rule set via Rules().
type Config struct {
	// Types restricts allowed commit types. Empty means the default
	// Conventional Commits type set.
	Types []string `yaml:"types"`

	Description DescriptionConfig `yaml:"description"`
	Scope       ScopeConfig       `yaml:"scope"`
	Header      HeaderConfig      `yaml:"header"`
	Rules       RulesConfig       `yaml:"rules"`
}

// DescriptionConfig controls description validation.
type DescriptionConfig struct {
	// Min is the minimum description length in runes. 0 = default (1).
	Min *int `yaml:"min"`
	// Max is the maximum description length in runes. 0 disables the limit
	// when explicitly set; nil means default (72).
	Max *int `yaml:"max"`
	// Lowercase requires the first letter to be lowercase.
	// nil means default (true).
	Lowercase *bool `yaml:"lowercase"`
}

// ScopeConfig controls scope validation.
type ScopeConfig struct {
	// Required demands a non-empty scope.
	Required bool `yaml:"required"`
	// Enum, when non-empty, restricts scope to one of the listed values.
	Enum []string `yaml:"enum"`
}

// HeaderConfig controls whole-header validation.
type HeaderConfig struct {
	// MaxLength limits the header line length. nil / 0 = disabled.
	MaxLength *int `yaml:"max_length"`
}

// RulesConfig toggles individual rules.
type RulesConfig struct {
	// NoTrailingPeriod forbids a trailing '.' in the description.
	// nil means default (false – not part of stock Conventional Commits).
	NoTrailingPeriod *bool `yaml:"no_trailing_period"`
	// BreakingChangeFooter requires BREAKING CHANGE footer when '!' is used.
	// nil means default (true).
	BreakingChangeFooter *bool `yaml:"breaking_change_footer"`
}

// LoadResult is returned by Load / FindAndLoad.
type LoadResult struct {
	Config Config
	Path   string // absolute path of the file that was loaded; empty if defaults
	Found  bool
}

// Default returns a Config equivalent to policy.ConventionalCommits().
func Default() Config {
	min := 1
	max := 72
	lowercase := true
	breaking := true

	return Config{
		Types: append([]string(nil), defaultTypes...),
		Description: DescriptionConfig{
			Min:       &min,
			Max:       &max,
			Lowercase: &lowercase,
		},
		Rules: RulesConfig{
			BreakingChangeFooter: &breaking,
		},
	}
}

var defaultTypes = []string{
	"feat", "fix", "docs", "style", "refactor", "perf",
	"test", "build", "ci", "chore", "revert",
}

// Load reads and parses a config file at the given path.
func Load(path string) (LoadResult, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return LoadResult{}, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return LoadResult{}, fmt.Errorf("parse config %s: %w", path, err)
	}

	abs, err := filepath.Abs(path)
	if err != nil {
		abs = path
	}

	return LoadResult{
		Config: cfg,
		Path:   abs,
		Found:  true,
	}, nil
}

// FindAndLoad walks up from startDir looking for .commitkit.yml.
// If none is found, returns Default() with Found=false.
func FindAndLoad(startDir string) (LoadResult, error) {
	if startDir == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return LoadResult{}, err
		}
		startDir = cwd
	}

	dir, err := filepath.Abs(startDir)
	if err != nil {
		return LoadResult{}, err
	}

	for {
		candidate := filepath.Join(dir, FileName)
		info, err := os.Stat(candidate)
		if err == nil && !info.IsDir() {
			return Load(candidate)
		}

		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return LoadResult{}, err
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached filesystem root.
			return LoadResult{
				Config: Default(),
				Found:  false,
			}, nil
		}
		dir = parent
	}
}

// ToRules builds a lint.Rules set from the config, applying Conventional
// Commits defaults for any field that was left unspecified.
func (c Config) ToRules() lint.Rules {
	var rules lint.Rules

	// Types
	types := c.Types
	if len(types) == 0 {
		types = defaultTypes
	}
	rules = append(rules, policy.TypeEnum(types...))

	// Description length
	min := 1
	if c.Description.Min != nil {
		min = *c.Description.Min
	}
	max := 72
	if c.Description.Max != nil {
		max = *c.Description.Max
	}
	rules = append(rules, policy.DescriptionLength(min, max))

	// Description lowercase
	lowercase := true
	if c.Description.Lowercase != nil {
		lowercase = *c.Description.Lowercase
	}
	if lowercase {
		rules = append(rules, policy.DescriptionLowercase())
	}

	// Scope required
	if c.Scope.Required {
		rules = append(rules, policy.RequireScope())
	}

	// Scope enum
	if len(c.Scope.Enum) > 0 {
		rules = append(rules, policy.ScopeEnum(c.Scope.Enum...))
	}

	// Header max length
	if c.Header.MaxLength != nil && *c.Header.MaxLength > 0 {
		rules = append(rules, policy.HeaderMaxLength(*c.Header.MaxLength))
	}

	// No trailing period
	noPeriod := false
	if c.Rules.NoTrailingPeriod != nil {
		noPeriod = *c.Rules.NoTrailingPeriod
	}
	if noPeriod {
		rules = append(rules, policy.NoTrailingPeriod())
	}

	// Breaking change footer
	breaking := true
	if c.Rules.BreakingChangeFooter != nil {
		breaking = *c.Rules.BreakingChangeFooter
	}
	if breaking {
		rules = append(rules, policy.BreakingChangeFooter())
	}

	return rules
}
