package latex

import "strings"

// specialChars maps each LaTeX special character to its escaped equivalent.
// Order matters: backslash must be first so we don't double-escape later replacements.
var replacer = strings.NewReplacer(
	`\`, `\textbackslash{}`,
	`$`, `\$`,
	`%`, `\%`,
	`_`, `\_`,
	`#`, `\#`,
	`&`, `\&`,
	`{`, `\{`,
	`}`, `\}`,
	`^`, `\textasciicircum{}`,
	`~`, `\textasciitilde{}`,
)

// Escape replaces LaTeX special characters in s with their safe equivalents.
func Escape(s string) string {
	return replacer.Replace(s)
}
