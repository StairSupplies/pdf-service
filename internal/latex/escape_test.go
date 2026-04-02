package latex_test

import (
	"testing"

	"github.com/StairSupplies/pdf-service/internal/latex"
)

func TestEscape(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Individual special characters
		{name: "backslash", input: `\`, want: `\textbackslash{}`},
		{name: "dollar", input: `$`, want: `\$`},
		{name: "percent", input: `%`, want: `\%`},
		{name: "underscore", input: `_`, want: `\_`},
		{name: "hash", input: `#`, want: `\#`},
		{name: "ampersand", input: `&`, want: `\&`},
		{name: "open brace", input: `{`, want: `\{`},
		{name: "close brace", input: `}`, want: `\}`},
		{name: "caret", input: `^`, want: `\textasciicircum{}`},
		{name: "tilde", input: `~`, want: `\textasciitilde{}`},

		// Plain text is unchanged
		{name: "plain text", input: "Hello World", want: "Hello World"},
		{name: "empty string", input: "", want: ""},

		// Combinations
		{
			name:  "price tag",
			input: `100% off $10`,
			want:  `100\% off \$10`,
		},
		{
			name:  "variable reference",
			input: `foo_bar`,
			want:  `foo\_bar`,
		},
		{
			name:  "math expression",
			input: `$x^2 + y^2$`,
			want:  `\$x\textasciicircum{}2 + y\textasciicircum{}2\$`,
		},
		{
			name:  "all special chars",
			input: `\$%_#&{}^~`,
			want:  `\textbackslash{}\$\%\_\#\&\{\}\textasciicircum{}\textasciitilde{}`,
		},
		{
			name:  "backslash not double-escaped",
			input: `C:\Users\name`,
			want:  `C:\textbackslash{}Users\textbackslash{}name`,
		},
		{
			name:  "braces in combination",
			input: `{foo}`,
			want:  `\{foo\}`,
		},
		{
			name:  "tilde in url-like string",
			input: `~/documents`,
			want:  `\textasciitilde{}/documents`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := latex.Escape(tc.input)
			if got != tc.want {
				t.Errorf("Escape(%q)\n  got  %q\n  want %q", tc.input, got, tc.want)
			}
		})
	}
}
