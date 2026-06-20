package maxcomments

import (
	"fmt"
	"regexp"
)

// compileIgnore compiles each ignore pattern into a regular expression. It
// fails loudly on the first invalid pattern rather than silently skipping it,
// so a typo in configuration cannot quietly disable a check.
func compileIgnore(patterns []string) ([]*regexp.Regexp, error) {
	compiled := make([]*regexp.Regexp, 0, len(patterns))
	for _, p := range patterns {
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, fmt.Errorf("invalid ignore pattern %q: %w", p, err)
		}
		compiled = append(compiled, re)
	}

	return compiled, nil
}

// matchesAny reports whether path matches any of the compiled patterns.
func matchesAny(patterns []*regexp.Regexp, path string) bool {
	for _, re := range patterns {
		if re.MatchString(path) {
			return true
		}
	}

	return false
}
